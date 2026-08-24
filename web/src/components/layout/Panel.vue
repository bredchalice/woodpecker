<template>
  <div class="lha-ci-panel w-full overflow-hidden">
    <component
      :is="collapsable ? 'button' : 'div'"
      v-if="title"
      type="button"
      class="lha-ci-panel__title flex w-full gap-2 px-4 py-2 font-bold"
      :class="{
        'cursor-pointer': collapsable,
      }"
      @click="_collapsed = !_collapsed"
    >
      <Icon
        v-if="collapsable"
        name="chevron-right"
        class="h-6 min-w-6 transition-transform duration-150"
        :class="{ 'rotate-90 transform': !collapsed }"
      />
      {{ title }}
    </component>
    <div
      :class="{
        'max-h-auto': !collapsed,
        'max-h-0': collapsed,
      }"
      class="transition-height overflow-hidden duration-150"
    >
      <div class="text-wp-text-100 w-full p-4">
        <slot />
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed, ref } from 'vue';

import Icon from '~/components/atomic/Icon.vue';

const props = defineProps<{
  title?: string;
  collapsable?: boolean;
  collapsedByDefault?: boolean;
}>();

const _collapsed = ref(props.collapsedByDefault || false);
const collapsed = computed(() => props.collapsable && _collapsed.value);
</script>

<style scoped>
.lha-ci-panel {
  border: 1px solid var(--lha-ci-border);
  border-radius: 0.9rem;
  background: var(--lha-ci-surface);
}

.lha-ci-panel__title {
  border-bottom: 1px solid var(--lha-ci-border);
  background: var(--lha-ci-surface-muted);
  color: var(--wp-text-200);
}
</style>
