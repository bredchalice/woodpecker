<template>
  <div class="flex flex-col gap-4">
    <div class="lha-ci-shell-intro">
      <div class="lha-ci-shell-intro__copy">
        <p class="lha-ci-kicker">Change validation</p>
        <h2 class="lha-ci-shell-intro__title">Pull requests</h2>
        <p class="lha-ci-shell-intro__text">Inspect CI activity for open changes without mixing pull-request validation into normal branch history.</p>
      </div>
      <div class="lha-ci-shell-intro__context">
        <span>Open changes</span>
        <strong>{{ pullRequests.length }}</strong>
      </div>
    </div>

    <div class="grid gap-2">
      <ListItem
        v-for="pullRequest in pullRequests"
        :key="pullRequest.index"
        class="text-wp-text-100 border-wp-background-400 bg-wp-background-100 dark:border-wp-background-100 dark:bg-wp-background-200 rounded-lg border"
        :to="{ name: 'repo-pull-request', params: { pullRequest: pullRequest.index } }"
      >
        <div class="flex min-w-0 items-center gap-3">
          <Icon name="pull-request" />
          <span class="text-wp-text-alt-100 shrink-0 font-mono text-xs">#{{ pullRequest.index }}</span>
          <span class="text-wp-text-100 overflow-hidden text-ellipsis whitespace-nowrap font-semibold">{{ pullRequest.title }}</span>
        </div>
      </ListItem>
    </div>

    <div v-if="loading" class="text-wp-text-100 flex justify-center"><Icon name="spinner" /></div>
    <Panel v-else-if="pullRequests.length === 0" class="flex justify-center">{{ $t('empty_list', { entity: $t('repo.pull_requests') }) }}</Panel>
  </div>
</template>

<script lang="ts" setup>
import { computed, watch } from 'vue';
import { useI18n } from 'vue-i18n';

import Icon from '~/components/atomic/Icon.vue';
import ListItem from '~/components/atomic/ListItem.vue';
import Panel from '~/components/layout/Panel.vue';
import useApiClient from '~/compositions/useApiClient';
import { requiredInject } from '~/compositions/useInjectProvide';
import { usePagination } from '~/compositions/usePaginate';
import { useWPTitle } from '~/compositions/useWPTitle';
import type { PullRequest } from '~/lib/api/types';

const apiClient = useApiClient();
const repo = requiredInject('repo');
if (!repo.value.pr_enabled || !repo.value.allow_pr) {
  throw new Error('Unexpected: pull requests are disabled for repo');
}

async function loadPullRequests(page: number): Promise<PullRequest[]> {
  return apiClient.getRepoPullRequests(repo.value.id, { page });
}

const { resetPage, data: pullRequests, loading } = usePagination(loadPullRequests);
watch(repo, resetPage);
const { t } = useI18n();
useWPTitle(computed(() => [t('repo.pull_requests'), repo.value.full_name]));
</script>
