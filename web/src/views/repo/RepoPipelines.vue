<template>
  <div class="lha-ci-activity">
    <header class="lha-ci-activity__header">
      <div>
        <p class="lha-ci-kicker">LHA Play CI</p>
        <h1 class="lha-ci-activity__title">Build history</h1>
        <p class="lha-ci-activity__intro">
          Recent validation, build, release and deployment activity for {{ repo.full_name }}.
        </p>
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
const runningCount = computed(
  () =>
    pipelines.value?.filter(
      (pipeline) =>
        pipeline.status === 'started' || pipeline.status === 'running' || pipeline.status === 'pending',
    ).length ?? 0,
);
const successCount = computed(
  () => pipelines.value?.filter((pipeline) => pipeline.status === 'success').length ?? 0,
);
const failedCount = computed(
  () =>
    pipelines.value?.filter((pipeline) => pipeline.status === 'failure' || pipeline.status === 'error').length ?? 0,
);
const supersededCount = computed(
  () =>
    pipelines.value?.filter(
      (pipeline) => pipeline.status === 'killed' && pipeline.cancel_info?.superseded_by,
    ).length ?? 0,
);

async function loadMore() {
  page.value += 1;
  await pipelineStore.loadRepoPipelines(repo.value.id, page.value);
}

const { t } = useI18n();
useWPTitle(computed(() => [t('repo.activity'), repo.value.full_name]));
</script>

<style scoped>
.lha-ci-activity {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.lha-ci-activity__header {
  display: flex;
  justify-content: space-between;
  gap: 2rem;
  padding: 1.5rem;
  border: 1px solid var(--lha-ci-border);
  border-radius: 1rem;
  background:
    radial-gradient(circle at top right, var(--lha-ci-accent-soft), transparent 42%),
    var(--lha-ci-surface);
}

.lha-ci-activity__title {
  margin: 0;
  color: var(--wp-text-200);
  font-size: 1.65rem;
  font-weight: 750;
}

.lha-ci-activity__intro {
  max-width: 680px;
  margin: 0.45rem 0 0;
  color: var(--wp-text-alt-100);
}

.lha-ci-activity__latest {
  display: grid;
  min-width: 150px;
  align-self: flex-start;
  gap: 0.15rem;
  padding: 0.75rem 0.9rem;
  border: 1px solid var(--lha-ci-border);
  border-radius: 0.75rem;
  background: var(--lha-ci-surface-muted);
}

.lha-ci-activity__latest span,
.lha-ci-activity__summary span {
  color: var(--wp-text-alt-100);
  font-size: 0.7rem;
  font-weight: 700;
  letter-spacing: 0.05em;
  text-transform: uppercase;
}

.lha-ci-activity__latest strong,
.lha-ci-activity__summary strong {
  color: var(--wp-text-200);
  font-size: 1.05rem;
}

.lha-ci-activity__latest small {
  color: var(--wp-text-alt-100);
  font-family: monospace;
}

.lha-ci-activity__summary {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 0.75rem;
}

.lha-ci-activity__summary > div {
  display: flex;
  flex-direction: column;
  gap: 0.2rem;
  padding: 0.85rem 1rem;
  border: 1px solid var(--lha-ci-border);
  border-radius: 0.75rem;
  background: var(--lha-ci-surface);
}

@media (max-width: 760px) {
  .lha-ci-activity__header {
    flex-direction: column;
  }

  .lha-ci-activity__latest {
    width: 100%;
  }

  .lha-ci-activity__summary {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}
</style>
