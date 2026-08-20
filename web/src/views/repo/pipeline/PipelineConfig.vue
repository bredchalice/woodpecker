<template>
  <div>
    <header class="lha-ci-workspace-head">
      <div class="lha-ci-workspace-head__copy">
        <p class="lha-ci-kicker">Resolved workflow</p>
        <h2 class="lha-ci-workspace-head__title">Pipeline configuration</h2>
        <p class="lha-ci-workspace-head__text">The exact workflow configuration Woodpecker resolved for build #{{ pipeline.number }}.</p>
      </div>
      <div class="lha-ci-workspace-head__context">
        <span>Configs</span>
        <strong>{{ pipelineConfigsDecoded.length }}</strong>
      </div>
    </header>

    <div class="lha-ci-workspace-list">
      <Panel
        v-for="pipelineConfig in pipelineConfigsDecoded"
        :key="pipelineConfig.hash"
        class="lha-ci-workspace-card"
        :collapsable="pipelineConfigsDecoded.length > 1"
        collapsed-by-default
        :title="pipelineConfigsDecoded.length > 1 ? pipelineConfig.name : ''"
      >
        <SyntaxHighlight
          class="code-box overflow-auto font-mono whitespace-pre"
          language="yaml"
          :code="pipelineConfig.data"
        />
      </Panel>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { decode } from 'js-base64';
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';

import SyntaxHighlight from '~/components/atomic/SyntaxHighlight';
import Panel from '~/components/layout/Panel.vue';
import { requiredInject } from '~/compositions/useInjectProvide';
import { useWPTitle } from '~/compositions/useWPTitle';

const repo = requiredInject('repo');
const pipeline = requiredInject('pipeline');
const pipelineConfigs = requiredInject('pipeline-configs');

const pipelineConfigsDecoded = computed(
  () =>
    pipelineConfigs.value?.map((i) => ({
      ...i,
      data: decode(i.data),
    })) ?? [],
);

const { t } = useI18n();
useWPTitle(
  computed(() => [
    t('repo.pipeline.config'),
    t('repo.pipeline.pipeline', { pipelineId: pipeline.value.number }),
    repo.value.full_name,
  ]),
);
</script>
