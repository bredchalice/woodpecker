<template>
  <Panel class="lha-ci-settings-panel">
    <header class="lha-ci-settings-head">
      <div class="lha-ci-settings-head__main">
        <div class="flex min-w-0 items-center justify-between gap-3">
          <h1 class="lha-ci-settings-head__title">
            {{ title }}
            <DocsLink v-if="docsUrl" :topic="title" :url="docsUrl" />
          </h1>
          <slot v-if="$slots.titleActions" name="titleActions" />
        </div>

        <div class="flex flex-wrap items-center justify-between gap-x-4 gap-y-2">
          <p v-if="description" class="lha-ci-settings-head__description">{{ description }}</p>
          <div v-if="$slots.headerActions" class="shrink-0">
            <slot name="headerActions" />
          </div>
        </div>
        <slot name="headerEnd" />
      </div>
    </header>

    <div class="lha-ci-settings-body">
      <slot />
    </div>
  </Panel>
</template>

<script setup lang="ts">
import DocsLink from '~/components/atomic/DocsLink.vue';
import Panel from '~/components/layout/Panel.vue';

defineProps<{
  title: string;
  description?: string;
  docsUrl?: string;
}>();
</script>
