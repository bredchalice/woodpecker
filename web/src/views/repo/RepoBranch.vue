<template>
  <div class="flex flex-col gap-4">
    <div class="lha-ci-shell-intro">
      <div class="lha-ci-shell-intro__copy">
        <p class="lha-ci-kicker">Branch activity</p>
        <h2 class="lha-ci-shell-intro__title">{{ branch }}</h2>
        <p class="lha-ci-shell-intro__text">Pipeline history for this source branch, excluding pull-request validation runs.</p>
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

const props = defineProps<{ branch: string }>();
const branch = toRef(props, 'branch');
const repo = requiredInject('repo');
const allPipelines = requiredInject('pipelines');
const pipelines = computed(() =>
  allPipelines.value.filter(
    (b) => b.branch === branch.value && b.event !== 'pull_request' && b.event !== 'pull_request_closed' && b.event !== 'pull_request_metadata',
  ),
);

const { t } = useI18n();
useWPTitle(computed(() => [t('repo.activity'), branch.value, repo.value.full_name]));
</script>
