<template>
  <router-link
    v-if="repo"
    :to="{ name: 'repo', params: { repoId: repo.id } }"
    class="lha-ci-repo-card"
  >
    <div class="lha-ci-repo-card__top">
      <div class="min-w-0">
        <span class="lha-ci-repo-card__owner">{{ repo.owner }}</span>
        <strong class="lha-ci-repo-card__name">{{ repo.name }}</strong>
      </div>
      <div class="text-wp-text-alt-100">
        <div
          v-if="repo.visibility === RepoVisibility.Private"
          :title="`${$t('repo.visibility.visibility')}: ${$t(`repo.visibility.private.private`)}`"
        >
          <Icon name="visibility-private" />
        </div>
        <div
          v-else-if="repo.visibility === RepoVisibility.Internal"
          :title="`${$t('repo.visibility.visibility')}: ${$t(`repo.visibility.internal.internal`)}`"
        >
          <Icon name="visibility-internal" />
        </div>
      </div>
    </div>

    <div class="lha-ci-repo-card__activity">
      <template v-if="lastPipeline">
        <div class="lha-ci-repo-card__message">
          <PipelineStatusIcon :status="lastPipeline.status" />
          <RenderMarkdown
            class="overflow-hidden text-ellipsis whitespace-nowrap"
            :title="message"
            :content="shortMessage"
            inline
          />
        </div>
        <div class="lha-ci-repo-card__since">
          <Icon name="since" />
          <span>{{ since }}</span>
        </div>
      </template>

      <div v-else class="lha-ci-repo-card__empty">
        <span>{{ $t('repo.pipeline.no_pipelines') }}</span>
      </div>
    </div>
  </router-link>
</template>

<script lang="ts" setup>
import { computed } from 'vue';

import Icon from '~/components/atomic/Icon.vue';
import RenderMarkdown from '~/components/atomic/RenderMarkdown.vue';
import PipelineStatusIcon from '~/components/repo/pipeline/PipelineStatusIcon.vue';
import usePipeline from '~/compositions/usePipeline';
import type { Repo } from '~/lib/api/types';
import { RepoVisibility } from '~/lib/api/types';

const props = defineProps<{
  repo: Repo;
}>();

const lastPipeline = computed(() => props.repo.last_pipeline);
const { since, shortMessage, message } = usePipeline(lastPipeline);
</script>
