<template>
  <div class="lha-ci-activity">
    <header class="lha-ci-activity__header">
      <div>
        <p class="lha-ci-kicker">LHA Play CI</p>
        <h1 class="lha-ci-activity__title">Build history</h1>
        <p class="lha-ci-activity__intro">Recent validation, build, release and deployment activity for {{ repo.full_name }}.</p>
      </div>

      <div v-if="latestPipeline" class="lha-ci-activity__latest">
        <span>Latest run</span>
        <strong>#{{ latestPipeline.number }}</strong>
        <small>{{ latestPipeline.commit.slice(0, 10) }}</small>
      </div>
    </header>

    <section v-if="pipelines?.length" class="lha-ci-activity__summary" aria-label="Build summary">
      <div>
        <span>Running</span>
        <strong>{{ runningCount }}</strong>
      </div>
      <div>
        <span>Passed</span>
        <strong>{{ successCount }}</strong>
      </div>
      <div>
        <span>Failed</span>
        <strong>{{ failedCount }}</strong>
      </div>
      <div>
        <span>Superseded</span>
        <strong>{{ supersededCount }}</strong>
      </div>
    </section>

    <PipelineList
      :pipelines="pipelines"
      :loading="pipelineStore.loading"
      :has-more="pipelineStore.hasMore"
      @load-more="loadMore"
    />
  </div>
</template>

<script lang="ts" setup>
import { computed, ref } from 'vue';
import { useI18n } from 'vue-i18n';

import PipelineList from '~/components/repo/pipeline/PipelineList.vue';
import { requiredInject } from '~/compositions/useInjectProvide';
import { useWPTitle } from '~/compositions/useWPTitle';
import { usePipelineStore } from '~/store/pipelines';

// TODO(4626): Refactor to use usePagination with server-side filtering,
// so pipeline loading can move from RepoWrapper to individual list views.
const repo = requiredInject('repo');
const pipelines = requiredInject('pipelines');
const pipelineStore = usePipelineStore();

const page = ref(1);

const latestPipeline = computed(() => pipelines.value?.[0]);
const runningCount = computed(() => pipelines.value?.filter((p) => p.status === 'started' || p.status === 'running' || p.status === 'pending').length ?? 0);
const successCount = computed(() => pipelines.value?.filter((p) => p.status === 'success').length ?? 0);
const failedCount = computed(() => pipelines.value?.filter((p) => p.status === 'failure' || p.status === 'error').length ?? 0);
const supersededCount = computed(() => pipelines.value?.filter((p) => p.status === 'killed' && p.cancel_info?.superseded_by).length ?? 0);

async function loadMore() {
  page.value += 1;
  await pipelineStore.loadRepoPipelines(repo.value.id, page.value);
}

const { t } = useI18n();
useWPTitle(computed(() => [t('repo.activity'), repo.value.full_name]));
</script>
