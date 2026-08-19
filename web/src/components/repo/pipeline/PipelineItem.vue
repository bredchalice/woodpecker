<template>
  <ListItem v-if="pipeline" class="lha-ci-pipeline-item w-full p-0!">
    <div class="lha-ci-pipeline-item__status" :data-status="semanticStatus">
      <PipelineRunningIcon v-if="pipeline.status === 'started' || pipeline.status === 'running'" />
      <PipelineStatusIcon v-else :status="pipeline.status" />
    </div>

    <div class="lha-ci-pipeline-item__body">
      <div class="lha-ci-pipeline-item__main">
        <div class="lha-ci-pipeline-item__headline">
          <span class="lha-ci-pipeline-item__number">#{{ pipeline.number }}</span>
          <span class="lha-ci-pipeline-item__semantic">{{ semanticLabel }}</span>
        </div>
        <RenderMarkdown class="lha-ci-pipeline-item__message" :title="message" :content="shortMessage" inline />
        <p v-if="statusDetail" class="lha-ci-pipeline-item__detail">{{ statusDetail }}</p>
      </div>

      <div class="lha-ci-pipeline-item__meta">
        <div class="flex min-w-0 items-center space-x-2">
          <span :title="pipelineEventTitle">
            <Icon v-if="pipeline.event === 'pull_request'" name="pull-request" />
            <Icon v-else-if="pipeline.event === 'pull_request_closed'" name="pull-request-closed" />
            <Icon v-else-if="pipeline.event === 'pull_request_metadata'" name="pull-request-metadata" />
            <Icon v-else-if="pipeline.event === 'deployment'" name="deployment" />
            <Icon v-else-if="pipeline.event === 'tag' || pipeline.event === 'release'" name="tag" />
            <Icon v-else-if="pipeline.event === 'cron'" name="stopwatch" />
            <Icon v-else-if="pipeline.event === 'manual'" name="manual-pipeline" />
            <Icon v-else name="branch" />
          </span>
          <span class="truncate">{{ prettyRef }}</span>
        </div>

        <div class="flex min-w-0 items-center space-x-2">
          <Icon name="commit" />
          <span class="truncate">{{ pipeline.commit.slice(0, 10) }}</span>
        </div>

        <div
          class="flex min-w-0 items-center space-x-2"
          :title="durationElapsed > 0 ? $t('repo.pipeline.duration', { duration: durationAsNumber(durationElapsed) }) : ''"
        >
          <Icon name="duration" />
          <span class="truncate">{{ duration }}</span>
        </div>

        <div class="flex min-w-0 items-center space-x-2" :title="$t('repo.pipeline.created', { created })">
          <Icon name="since" />
          <span class="truncate">{{ since }}</span>
        </div>
      </div>
    </div>
  </ListItem>
</template>

<script lang="ts" setup>
import { computed, toRef } from 'vue';
import { useI18n } from 'vue-i18n';

import Icon from '~/components/atomic/Icon.vue';
import ListItem from '~/components/atomic/ListItem.vue';
import RenderMarkdown from '~/components/atomic/RenderMarkdown.vue';
import PipelineRunningIcon from '~/components/repo/pipeline/PipelineRunningIcon.vue';
import PipelineStatusIcon from '~/components/repo/pipeline/PipelineStatusIcon.vue';
import { useDate } from '~/compositions/useDate';
import usePipeline from '~/compositions/usePipeline';
import type { Pipeline } from '~/lib/api/types';

const props = defineProps<{
  pipeline: Pipeline;
}>();

const { t } = useI18n();
const { durationAsNumber } = useDate();

const pipeline = toRef(props, 'pipeline');
const { since, duration, durationElapsed, message, shortMessage, prettyRef, created } = usePipeline(pipeline);

const semanticStatus = computed(() => {
  if (pipeline.value.status === 'killed' && pipeline.value.cancel_info?.superseded_by) return 'superseded';
  if (pipeline.value.status === 'killed') return 'cancelled';
  if (pipeline.value.status === 'success') return 'passed';
  if (pipeline.value.status === 'failure' || pipeline.value.status === 'error') return 'failed';
  if (pipeline.value.status === 'started' || pipeline.value.status === 'running') return 'running';
  if (pipeline.value.status === 'pending') return 'queued';
  return 'neutral';
});

const semanticLabel = computed(() => {
  switch (semanticStatus.value) {
    case 'superseded':
      return 'Superseded';
    case 'cancelled':
      return 'Cancelled';
    case 'passed':
      return 'Passed';
    case 'failed':
      return 'Failed';
    case 'running':
      return 'Running';
    case 'queued':
      return 'Queued';
    default:
      return pipeline.value.status;
  }
});

const statusDetail = computed(() => {
  const info = pipeline.value.cancel_info;
  if (pipeline.value.status !== 'killed' || !info) return '';
  if (info.superseded_by) return `Replaced by build #${info.superseded_by}. No action required.`;
  if (info.canceled_by_user) return `Cancelled by ${info.canceled_by_user}.`;
  if (info.canceled_by_step) return `Cancelled by step ${info.canceled_by_step}.`;
  return '';
});

const pipelineEventTitle = computed(() => {
  switch (pipeline.value.event) {
    case 'pull_request':
      return t('repo.pipeline.event.pr');
    case 'pull_request_closed':
      return t('repo.pipeline.event.pr_closed');
    case 'pull_request_metadata':
      return t('repo.pipeline.event.pr_metadata');
    case 'deployment':
      return t('repo.pipeline.event.deploy');
    case 'tag':
      return t('repo.pipeline.event.tag');
    case 'release':
      return t('repo.pipeline.event.release');
    case 'cron':
      return t('repo.pipeline.event.cron');
    case 'manual':
      return t('repo.pipeline.event.manual');
    default:
      return t('repo.pipeline.event.push');
  }
});
</script>
