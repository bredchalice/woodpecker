<template>
  <div>
    <header class="lha-ci-workspace-head">
      <div class="lha-ci-workspace-head__copy">
        <p class="lha-ci-kicker">Pipeline diagnostics</p>
        <h2 class="lha-ci-workspace-head__title">{{ hasErrors ? 'Errors and warnings' : 'Warnings' }}</h2>
        <p class="lha-ci-workspace-head__text">Configuration and execution diagnostics reported for build #{{ pipeline.number }}.</p>
      </div>
      <div class="lha-ci-workspace-head__context">
        <span>Diagnostics</span>
        <strong>{{ pipeline.errors?.length ?? 0 }}</strong>
      </div>
    </header>

    <Panel class="lha-ci-workspace-card">
      <div class="lha-ci-workspace-list">
        <template v-for="(error, index) in pipeline.errors" :key="index">
          <article class="rounded-lg border p-4">
            <div class="grid gap-2 md:grid-cols-[minmax(10rem,auto)_3fr]">
              <span class="flex items-center gap-x-2">
                <Icon
                  name="alert"
                  class="my-1 shrink-0"
                  :class="{
                    'text-wp-state-warn-100': error.is_warning,
                    'text-wp-error-100': !error.is_warning,
                  }"
                />
                <code>{{ error.type }}</code>
              </span>
              <span
                v-if="isLinterError(error) || isDeprecationError(error) || isBadHabitError(error)"
                class="flex min-w-0 items-center gap-x-2"
              >
                <span class="min-w-0 truncate">
                  <span v-if="error.data?.file" class="font-bold">{{ error.data.file }}: </span>
                  <span>{{ error.data?.field }}</span>
                </span>
                <DocsLink
                  v-if="isDeprecationError(error) || isBadHabitError(error)"
                  :topic="error.data?.field || ''"
                  :url="error.data?.docs || ''"
                />
              </span>
            </div>
            <div class="mt-3 md:pl-10">
              <RenderMarkdown :content="error.message" />
            </div>
          </article>
        </template>
      </div>
    </Panel>
  </div>
</template>

<script lang="ts" setup>
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';

import DocsLink from '~/components/atomic/DocsLink.vue';
import Icon from '~/components/atomic/Icon.vue';
import RenderMarkdown from '~/components/atomic/RenderMarkdown.vue';
import Panel from '~/components/layout/Panel.vue';
import { requiredInject } from '~/compositions/useInjectProvide';
import { useWPTitle } from '~/compositions/useWPTitle';
import type { PipelineError } from '~/lib/api/types';

const repo = requiredInject('repo');
const pipeline = requiredInject('pipeline');
const hasErrors = computed(() => pipeline.value.errors?.some((error) => !error.is_warning) ?? false);

function isLinterError(error: PipelineError): error is PipelineError<{ file?: string; field: string }> {
  return error.type === 'linter';
}

function isDeprecationError(
  error: PipelineError,
): error is PipelineError<{ file: string; field: string; docs: string }> {
  return error.type === 'deprecation';
}

function isBadHabitError(error: PipelineError): error is PipelineError<{ file?: string; field: string; docs: string }> {
  return error.type === 'bad_habit';
}

const { t } = useI18n();
useWPTitle(
  computed(() => [
    hasErrors.value ? t('repo.pipeline.errors') : t('repo.pipeline.warnings'),
    t('repo.pipeline.pipeline', { pipelineId: pipeline.value.number }),
    repo.value.full_name,
  ]),
);
</script>
