<template>
  <Container full-width class="md:min-h-xs flex grow-0 flex-col md:grow md:px-4">
    <div class="flex min-h-0 w-full grow flex-wrap-reverse md:flex-nowrap md:gap-4">
      <PipelineStepList
        v-model:selected-step-id="selectedStepId"
        :class="{ 'hidden md:flex': pipeline!.status === 'blocked' }"
        :pipeline="pipeline!"
      />

      <div class="relative flex min-w-0 grow basis-full flex-col items-stretch gap-3 md:basis-auto">
        <div v-if="hasErrors" class="lha-ci-failure-summary">
          <div class="lha-ci-failure-summary__icon">
            <Icon name="status-error" class="text-wp-error-100" size="1.4rem" />
          </div>
          <div class="lha-ci-failure-summary__body">
            <p class="lha-ci-kicker">Build failed</p>
            <strong>{{ $t('repo.pipeline.we_got_some_errors') }}</strong>
            <span>The failed step remains selected below so the useful log stays visible.</span>
          </div>
          <Button color="red" :text="$t('repo.pipeline.show_errors')" :to="{ name: 'repo-pipeline-errors' }" />
        </div>

        <div v-if="pipeline!.status === 'blocked'" class="mb-4 w-full md:mb-auto">
          <Panel>
            <div class="flex flex-col items-center gap-4">
              <Icon name="status-blocked" size="1.5rem" class="h-16 w-16" />
              <span class="text-xl">{{ $t('repo.pipeline.protected.awaits') }}</span>
              <div v-if="repoPermissions!.push" class="flex flex-wrap items-center justify-center gap-2">
                <Button
                  color="green"
                  :text="$t('repo.pipeline.protected.approve')"
                  :is-loading="isApprovingPipeline"
                  @click="approvePipeline"
                />
                <Button
                  color="red"
                  :text="$t('repo.pipeline.protected.decline')"
                  :is-loading="isDecliningPipeline"
                  @click="declinePipeline"
                />
              </div>
            </div>
          </Panel>
        </div>

        <div v-else-if="pipeline!.status === 'declined'" class="mb-4 w-full md:mb-auto">
          <Panel>
            <div class="flex flex-col items-center gap-4">
              <Icon name="status-declined" size="1.5rem" class="text-wp-error-100 h-16 w-16" />
              <p class="text-xl">{{ $t('repo.pipeline.protected.declined') }}</p>
            </div>
          </Panel>
        </div>

        <PipelineLog v-else-if="selectedStepId !== null" v-model:step-id="selectedStepId" :pipeline="pipeline!" />
      </div>
    </div>
  </Container>
</template>

<script lang="ts" setup>
import { computed, toRef } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRoute, useRouter } from 'vue-router';

import Button from '~/components/atomic/Button.vue';
import Icon from '~/components/atomic/Icon.vue';
import Container from '~/components/layout/Container.vue';
import Panel from '~/components/layout/Panel.vue';
import PipelineLog from '~/components/repo/pipeline/PipelineLog.vue';
import PipelineStepList from '~/components/repo/pipeline/PipelineStepList.vue';
import useApiClient from '~/compositions/useApiClient';
import { useAsyncAction } from '~/compositions/useAsyncAction';
import { requiredInject } from '~/compositions/useInjectProvide';
import useNotifications from '~/compositions/useNotifications';
import { useWPTitle } from '~/compositions/useWPTitle';
import type { PipelineStep } from '~/lib/api/types';

const props = defineProps<{
  stepId?: string | null;
}>();

const apiClient = useApiClient();
const router = useRouter();
const route = useRoute();
const notifications = useNotifications();
const i18n = useI18n();

const pipeline = requiredInject('pipeline');
const repo = requiredInject('repo');
const repoPermissions = requiredInject('repo-permissions');

const stepId = toRef(props, 'stepId');
const hasErrors = computed(() => pipeline.value?.errors?.some((error) => !error.is_warning) ?? false);

const defaultStepId = computed(() => pipeline.value?.workflows?.[0].children?.[0].pid ?? null);

const selectedStepId = computed({
  get() {
    if (stepId.value !== '' && stepId.value !== null && stepId.value !== undefined) {
      const id = Number.parseInt(stepId.value, 10);

      let step = pipeline.value.workflows?.find((workflow) => workflow.pid === id)?.children[0];
      if (step) {
        return step.pid;
      }

      step = pipeline.value?.workflows?.reduce(
        (prev, workflow) => prev || workflow.children?.find((child) => child.pid === id),
        undefined as PipelineStep | undefined,
      );
      if (step) {
        return step.pid;
      }

      return defaultStepId.value;
    }

    if (window.innerWidth > 768) {
      return defaultStepId.value;
    }

    return null;
  },
  set(_selectedStepId: number | null) {
    if (_selectedStepId === null) {
      router.replace({ params: { ...route.params, stepId: '' } });
      return;
    }

    router.replace({ params: { ...route.params, stepId: `${_selectedStepId}` } });
  },
});

const { doSubmit: approvePipeline, isLoading: isApprovingPipeline } = useAsyncAction(async () => {
  await apiClient.approvePipeline(repo.value.id, `${pipeline.value.number}`);
  notifications.notify({ title: i18n.t('repo.pipeline.protected.approve_success'), type: 'success' });
});

const { doSubmit: declinePipeline, isLoading: isDecliningPipeline } = useAsyncAction(async () => {
  await apiClient.declinePipeline(repo.value.id, `${pipeline.value.number}`);
  notifications.notify({ title: i18n.t('repo.pipeline.protected.decline_success'), type: 'success' });
});

useWPTitle(
  computed(() => [
    i18n.t('repo.pipeline.tasks'),
    i18n.t('repo.pipeline.pipeline', { pipelineId: pipeline.value.number }),
    repo.value.full_name,
  ]),
);
</script>

<style scoped>
.lha-ci-failure-summary {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1rem 1.1rem;
  border: 1px solid color-mix(in srgb, var(--wp-error-100) 45%, var(--lha-ci-border));
  border-radius: 0.9rem;
  background: color-mix(in srgb, var(--wp-error-100) 8%, var(--lha-ci-surface));
}

.lha-ci-failure-summary__icon {
  display: flex;
  width: 2.75rem;
  height: 2.75rem;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  border-radius: 0.75rem;
  background: color-mix(in srgb, var(--wp-error-100) 12%, transparent);
}

.lha-ci-failure-summary__body {
  display: flex;
  min-width: 0;
  flex: 1;
  flex-direction: column;
  gap: 0.15rem;
}

.lha-ci-failure-summary__body strong {
  color: var(--wp-text-200);
}

.lha-ci-failure-summary__body span {
  color: var(--wp-text-alt-100);
  font-size: 0.82rem;
}

@media (max-width: 760px) {
  .lha-ci-failure-summary {
    align-items: stretch;
    flex-direction: column;
  }
}
</style>
