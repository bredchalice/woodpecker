<template>
  <Container full-width class="md:min-h-xs flex grow-0 flex-col md:grow md:px-4">
    <div class="flex min-h-0 w-full grow flex-wrap-reverse md:flex-nowrap md:gap-4">
      <PipelineStepList
        v-model:selected-step-id="selectedStepId"
        :class="{ 'hidden md:flex': pipeline!.status === 'blocked' }"
        :pipeline="pipeline!"
      />

      <div class="relative flex min-w-0 grow basis-full flex-col items-stretch gap-3 md:basis-auto">
        <div v-if="failureSummary" class="lha-ci-failure-summary">
          <div class="lha-ci-failure-summary__icon">
            <Icon name="status-error" class="text-wp-error-100" size="1.4rem" />
          </div>
          <div class="lha-ci-failure-summary__body">
            <p class="lha-ci-kicker">{{ failureSummary.kicker }}</p>
            <strong>{{ failureSummary.title }}</strong>
            <span>{{ failureSummary.detail }}</span>
          </div>
          <Button
            v-if="hasErrors"
            color="red"
            :text="$t('repo.pipeline.show_errors')"
            :to="{ name: 'repo-pipeline-errors' }"
          />
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

const orderedSteps = computed(() =>
  (pipeline.value?.workflows ?? []).flatMap((workflow) => workflow.children ?? []),
);

const automaticStepId = computed(() => {
  const steps = orderedSteps.value;
  if (steps.length === 0) return null;

  // A live pipeline should always show the first currently executing step. This
  // also makes parallel workflows deterministic: the first running step in UI
  // order wins.
  const running = steps.find((step) => step.state === 'running' || step.state === 'started');
  if (running) return running.pid;

  // Once a build has failed, put the useful failed log in front of the user.
  const failed = steps.find((step) => step.state === 'failure' || step.state === 'error');
  if (failed) return failed.pid;

  // Between steps (or for a completed successful build), keep the most recently
  // completed step visible instead of jumping back to the first step.
  const completedStates = new Set(['success', 'failure', 'error', 'killed', 'canceled', 'skipped']);
  const completed = [...steps].reverse().find((step) => completedStates.has(step.state));
  if (completed) return completed.pid;

  // Nothing has started yet. Prefer the first pending step and finally fall back
  // to the first step if the backend reports a state we do not know yet.
  return steps.find((step) => step.state === 'pending')?.pid ?? steps[0]?.pid ?? null;
});

function findStep(id: number | null): PipelineStep | undefined {
  if (id === null) return undefined;
  return orderedSteps.value.find((step) => step.pid === id);
}

const selectedStepId = computed({
  get() {
    // An explicit step in the route is a user/manual selection and must remain
    // pinned. With no route step, selection stays automatic and follows pipeline
    // execution as states change.
    if (stepId.value !== '' && stepId.value !== null && stepId.value !== undefined) {
      const id = Number.parseInt(stepId.value, 10);
      const explicitStep = findStep(id);
      return explicitStep?.pid ?? automaticStepId.value;
    }

    if (window.innerWidth > 768) {
      return automaticStepId.value;
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

const selectedStep = computed(() => findStep(selectedStepId.value));
const firstPipelineError = computed(() => pipeline.value?.errors?.find((error) => !error.is_warning));

const failureSummary = computed(() => {
  const step = selectedStep.value;
  if (step && (step.state === 'failure' || step.state === 'error')) {
    const exitCode = step.exit_code !== undefined && step.exit_code !== null ? `Exit code ${step.exit_code}.` : '';
    const detail = step.error || firstPipelineError.value?.message || exitCode || 'The step failed. Open its log for the command output.';
    return {
      kicker: 'Step failed',
      title: step.name,
      detail: exitCode && detail !== exitCode ? `${exitCode} ${detail}` : detail,
    };
  }

  if (firstPipelineError.value) {
    return {
      kicker: 'Pipeline error',
      title: firstPipelineError.value.type || i18n.t('repo.pipeline.we_got_some_errors'),
      detail: firstPipelineError.value.message || i18n.t('repo.pipeline.we_got_some_errors'),
    };
  }

  return null;
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
  line-height: 1.45;
}

@media (max-width: 760px) {
  .lha-ci-failure-summary {
    align-items: stretch;
    flex-direction: column;
  }
}
</style>
