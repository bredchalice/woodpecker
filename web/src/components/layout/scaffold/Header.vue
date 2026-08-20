<template>
  <header class="lha-ci-page-header text-wp-text-100" :class="{ 'md:px-4': fullWidth }">
    <Container :full-width="fullWidth" class="relative py-0!">
      <div class="lha-ci-page-header__inner">
        <div class="flex w-full flex-col gap-2 py-3 md:flex-row md:items-center md:justify-between md:gap-10">
          <div
            class="flex min-h-10 content-start items-center"
            :class="{
              'md:flex-1': searchBoxPresent,
            }"
          >
            <IconButton
              v-if="goBack"
              icon="back"
              :title="$t('back')"
              class="md:display-unset mr-2 hidden h-8 w-8 shrink-0 md:justify-between"
              @click="goBack"
            />
            <h1 class="lha-ci-page-header__title flex min-w-0 items-center gap-x-2">
              <slot name="title" />
            </h1>
          </div>
          <TextField
            v-if="searchBoxPresent"
            class="order-3 w-full grow md:order-none md:w-auto"
            :aria-label="$t('search')"
            :placeholder="$t('search')"
            :model-value="search"
            @update:model-value="(value: string) => $emit('update:search', value)"
          />
          <div
            v-if="$slots.headerActions"
            class="flex min-w-0 items-center gap-x-2 md:justify-end"
            :class="{
              'md:flex-1': searchBoxPresent,
            }"
          >
            <slot name="headerActions" />
          </div>
        </div>

        <div v-if="enableTabs" class="flex flex-col gap-4 md:flex-row md:items-center md:justify-between md:py-0">
          <Tabs class="order-2 md:order-none" />
          <div v-if="$slots.tabActions" class="flex flex-wrap content-start md:justify-end">
            <slot name="tabActions" />
          </div>
        </div>
      </div>
    </Container>
  </header>
</template>

<script setup lang="ts">
import { computed } from 'vue';

import IconButton from '~/components/atomic/IconButton.vue';
import TextField from '~/components/form/TextField.vue';
import Container from '~/components/layout/Container.vue';

import Tabs from './Tabs.vue';

const props = defineProps<{
  goBack?: () => void;
  enableTabs?: boolean;
  search?: string;
  fullWidth?: boolean;
}>();

defineEmits<{
  (event: 'update:search', query: string): void;
}>();

const searchBoxPresent = computed(() => props.search !== undefined);
</script>

<style scoped>
.lha-ci-page-header {
  background:
    linear-gradient(180deg, var(--lha-ci-accent-soft), transparent 85%),
    var(--wp-background-100);
}

.lha-ci-page-header__inner {
  border-bottom: 1px solid var(--lha-ci-border);
}

.lha-ci-page-header__title {
  color: var(--wp-text-200);
  font-size: 1.1rem;
  font-weight: 760;
  letter-spacing: -0.02em;
}
</style>
