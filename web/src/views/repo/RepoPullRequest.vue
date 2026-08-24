<template>
  <div class="flex flex-col gap-4">
    <div class="lha-ci-shell-intro">
      <div class="lha-ci-shell-intro__copy">
        <p class="lha-ci-kicker">Pull request validation</p>
        <h2 class="lha-ci-shell-intro__title">PR #{{ pullRequest }}</h2>
        <p class="lha-ci-shell-intro__text">Validation runs associated with this pull request, including metadata and close events when available.</p>
      </div>
      <div class="lha-ci-shell-intro__context">
        <span>Runs loaded</span>
        <strong>{{ pipelines.length }}</strong>
      </div>
    </div>
    <PipelineList :pipelines="pipelines" :repo="repo" />
  </div>
</template>

<script lang="ts" setup>
import { computed, toRef } from 'vue';
import { useI18n } from 'vue-i18n';

import PipelineList from '~/components/repo/pipeline/PipelineList.vue';
import { requiredInject } from '~/compositions/useInjectProvide';
import { useWPTitle } from '~/compositions/useWPTitle';

const props = defineProps<{ pullRequest: string }>();
const pullRequest = toRef(props, 'pullRequest');
const repo = requiredInject('repo');
if (!repo.value.pr_enabled || !repo.value.allow_pr) {
  throw new Error('Unexpected: pull requests are disabled for repo');
}

const allPipelines = requiredInject('pipelines');
const pipelines = computed(() =>
  allPipelines.value.filter(
    (b) =>
      (b.event === 'pull_request' || b.event === 'pull_request_closed' || b.event === 'pull_request_metadata') &&
      b.ref
        .replaceAll('refs/pull/', '')
        .replaceAll('refs/merge-requests/', '')
        .replaceAll('refs/pull-requests/', '')
        .replaceAll('/from', '')
        .replaceAll('/merge', '')
        .replaceAll('/head', '') === pullRequest.value,
  ),
);

const { t } = useI18n();
useWPTitle(computed(() => [t('repo.activity'), pullRequest.value, repo.value.full_name]));
</script>
