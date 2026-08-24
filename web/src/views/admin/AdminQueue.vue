<template>
  <Settings :title="$t('admin.settings.queue.queue')" :description="$t('admin.settings.queue.desc')">
    <template #headerActions>
      <div v-if="queueInfo" class="flex items-center gap-2">
        <Button v-if="queueInfo.paused" :text="$t('admin.settings.queue.resume')" start-icon="play" @click="resumeQueue" />
        <Button v-else :text="$t('admin.settings.queue.pause')" start-icon="pause" @click="pauseQueue" />
        <Icon
          :name="queueInfo.paused ? 'pause' : 'play'"
          class="h-6 w-6"
          :class="{ 'text-wp-error-100': queueInfo.paused, 'text-wp-state-ok-100': !queueInfo.paused }"
        />
      </div>
    </template>

    <div class="flex flex-col gap-4">
      <div class="lha-ci-shell-intro">
        <div class="lha-ci-shell-intro__copy">
          <p class="lha-ci-kicker">Scheduler</p>
          <h2 class="lha-ci-shell-intro__title">Execution queue</h2>
          <p class="lha-ci-shell-intro__text">See what is running now, what is waiting for capacity, and which tasks are blocked by dependencies.</p>
        </div>
        <div class="lha-ci-shell-intro__context">
          <span>Queue state</span>
          <strong>{{ queueInfo?.paused ? 'Paused' : 'Accepting work' }}</strong>
        </div>
      </div>

      <div class="grid grid-cols-1 gap-3 sm:grid-cols-3">
        <div class="border-wp-background-400 bg-wp-background-100 dark:border-wp-background-100 dark:bg-wp-background-200 rounded-lg border p-4">
          <span class="text-wp-text-alt-100 text-xs font-bold uppercase tracking-wide">Running</span>
          <strong class="text-wp-text-200 mt-1 block text-2xl">{{ runningCount }}</strong>
        </div>
        <div class="border-wp-background-400 bg-wp-background-100 dark:border-wp-background-100 dark:bg-wp-background-200 rounded-lg border p-4">
          <span class="text-wp-text-alt-100 text-xs font-bold uppercase tracking-wide">Pending</span>
          <strong class="text-wp-text-200 mt-1 block text-2xl">{{ pendingCount }}</strong>
        </div>
        <div class="border-wp-background-400 bg-wp-background-100 dark:border-wp-background-100 dark:bg-wp-background-200 rounded-lg border p-4">
          <span class="text-wp-text-alt-100 text-xs font-bold uppercase tracking-wide">Blocked</span>
          <strong class="text-wp-text-200 mt-1 block text-2xl">{{ blockedCount }}</strong>
        </div>
      </div>

      <AdminQueueStats :stats="queueInfo?.stats" />

      <div v-if="tasks.length > 0" class="flex flex-col gap-2">
        <div class="lha-ci-section-heading mb-0!">
          <div>
            <h2>{{ $t('admin.settings.queue.tasks') }}</h2>
            <span>{{ tasks.length }} active queue item{{ tasks.length === 1 ? '' : 's' }}</span>
          </div>
        </div>

        <ListItem
          v-for="task in tasks"
          :key="task.id"
          class="bg-wp-background-100! dark:bg-wp-background-200! border-wp-background-400 dark:border-wp-background-100 flex-col items-center gap-3 rounded-lg border"
        >
          <div class="flex w-full items-center justify-between gap-3">
            <div class="flex min-w-0 items-center gap-3">
              <Icon
                :name="task.status === 'pending' ? 'status-pending' : task.status === 'running' ? 'status-running' : 'status-declined'"
                :class="{
                  'text-wp-error-100': task.status === 'waiting_on_deps',
                  'text-wp-state-info-100': task.status === 'running',
                  'text-wp-state-neutral-100': task.status === 'pending',
                }"
              />
              <div class="min-w-0">
                <span class="text-wp-text-200 block truncate font-semibold">{{ task.name }}</span>
                <span class="text-wp-text-alt-100 text-xs">{{ task.status.replaceAll('_', ' ') }}</span>
              </div>
            </div>

            <div class="ml-auto flex items-center gap-2">
              <Badge v-if="task.agent_name" :label="$t('admin.settings.queue.agent')" :value="task.agent_name" />
              <Badge
                v-if="task.dependencies.length > 0"
                :label="$t('admin.settings.queue.waiting_for')"
                :value="task.dependencies.join(', ')"
              />
              <IconButton
                v-if="task.pipeline_number"
                icon="chevron-right"
                :title="$t('repo.pipeline.view')"
                class="h-8 w-8"
                :to="{ name: 'repo-pipeline', params: { repoId: task.repo_id, pipelineId: task.pipeline_number, stepId: task.pid } }"
              />
            </div>
          </div>

          <div v-if="Object.values(task.labels).some(Boolean)" class="border-wp-background-400 dark:border-wp-background-100 flex w-full flex-wrap gap-2 border-t pt-3">
            <template v-for="(value, label) in task.labels">
              <Badge v-if="value" :key="label" :label="label.toString()" :value="value" />
            </template>
          </div>
        </ListItem>
      </div>
    </div>
  </Settings>
</template>

<script lang="ts" setup>
import { computed, ref } from 'vue';
import { useI18n } from 'vue-i18n';

import AdminQueueStats from '~/components/admin/settings/queue/AdminQueueStats.vue';
import Badge from '~/components/atomic/Badge.vue';
import Button from '~/components/atomic/Button.vue';
import Icon from '~/components/atomic/Icon.vue';
import IconButton from '~/components/atomic/IconButton.vue';
import ListItem from '~/components/atomic/ListItem.vue';
import Settings from '~/components/layout/Settings.vue';
import useApiClient from '~/compositions/useApiClient';
import { useInterval } from '~/compositions/useInterval';
import useNotifications from '~/compositions/useNotifications';
import { useWPTitle } from '~/compositions/useWPTitle';
import type { QueueInfo } from '~/lib/api/types';

const apiClient = useApiClient();
const notifications = useNotifications();
const { t } = useI18n();
const queueInfo = ref<QueueInfo>();
const runningCount = computed(() => queueInfo.value?.running?.length ?? 0);
const pendingCount = computed(() => queueInfo.value?.pending?.length ?? 0);
const blockedCount = computed(() => queueInfo.value?.waiting_on_deps?.length ?? 0);

const tasks = computed(() => {
  const _tasks = [];
  if (queueInfo.value?.running) _tasks.push(...queueInfo.value.running.map((task) => ({ ...task, status: 'running' })));
  if (queueInfo.value?.pending) _tasks.push(...queueInfo.value.pending.map((task) => ({ ...task, status: 'pending' })));
  if (queueInfo.value?.waiting_on_deps) _tasks.push(...queueInfo.value.waiting_on_deps.map((task) => ({ ...task, status: 'waiting_on_deps' })));
  return _tasks
    .map((task) => ({ ...task, labels: Object.fromEntries(Object.entries(task.labels).filter(([key]) => key !== 'org-id')) }))
    .toSorted((a, b) => a.id - b.id);
});

async function loadQueueInfo() {
  queueInfo.value = await apiClient.getQueueInfo();
}

async function pauseQueue() {
  await apiClient.pauseQueue();
  await loadQueueInfo();
  notifications.notify({ title: t('admin.settings.queue.paused'), type: 'success' });
}

async function resumeQueue() {
  await apiClient.resumeQueue();
  await loadQueueInfo();
  notifications.notify({ title: t('admin.settings.queue.resumed'), type: 'success' });
}

useInterval(loadQueueInfo, 5 * 1000);
useWPTitle(computed(() => [t('admin.settings.queue.queue'), t('admin.settings.settings')]));
</script>
