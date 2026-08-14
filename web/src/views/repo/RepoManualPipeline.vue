<template>
  <Panel v-if="!loading">
    <form @submit.prevent="triggerManualPipeline">
      <span class="text-wp-text-100 text-xl">{{ $t('repo.manual_pipeline.title') }}</span>
      <InputField v-slot="{ id }" :label="$t('repo.manual_pipeline.select_branch')">
        <SelectField :id="id" v-model="payload.branch" :options="branches" required />
      </InputField>

      <div v-if="manualInputEntries.length" class="mb-4">
        <template v-for="[name, input] in manualInputEntries" :key="name">
          <Checkbox
            v-if="input.type === 'boolean'"
            v-model="typedValues[name] as boolean"
            :label="formatInputName(name)"
            :description="input.description"
          />

          <InputField v-else v-slot="{ id }" :label="formatInputName(name)">
            <span v-if="input.description" class="text-wp-text-alt-100 mb-2 text-sm">{{ input.description }}</span>
            <SelectField
              v-if="input.type === 'choice'"
              :id="id"
              v-model="typedValues[name] as string"
              :options="input.options.map((option) => ({ text: option, value: option }))"
              :placeholder="input.required ? undefined : '-'"
            />
            <TextField v-else :id="id" v-model="typedValues[name] as string" />
          </InputField>
        </template>
      </div>

      <InputField v-slot="{ id }" :label="$t('repo.manual_pipeline.variables.title')">
        <span class="text-wp-text-alt-100 mb-2 text-sm">{{ $t('repo.manual_pipeline.variables.desc') }}</span>
        <KeyValueEditor
          :id="id"
          v-model="payload.variables"
          :key-placeholder="$t('repo.manual_pipeline.variables.name')"
          :value-placeholder="$t('repo.manual_pipeline.variables.value')"
          :delete-title="$t('repo.manual_pipeline.variables.delete')"
          @update:is-valid="isVariablesValid = $event"
        />
      </InputField>
      <Button type="submit" :text="$t('repo.manual_pipeline.trigger')" :disabled="!isFormValid" />
    </form>
  </Panel>
  <div v-else class="text-wp-text-100 flex justify-center">
    <Icon name="spinner" />
  </div>
</template>

<script lang="ts" setup>
import { useNotification } from '@kyvg/vue3-notification';
import { computed, onMounted, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';

import Button from '~/components/atomic/Button.vue';
import Icon from '~/components/atomic/Icon.vue';
import Checkbox from '~/components/form/Checkbox.vue';
import InputField from '~/components/form/InputField.vue';
import KeyValueEditor from '~/components/form/KeyValueEditor.vue';
import SelectField from '~/components/form/SelectField.vue';
import TextField from '~/components/form/TextField.vue';
import Panel from '~/components/layout/Panel.vue';
import useApiClient from '~/compositions/useApiClient';
import useConfig from '~/compositions/useConfig';
import { requiredInject } from '~/compositions/useInjectProvide';
import { usePaginate } from '~/compositions/usePaginate';
import { useWPTitle } from '~/compositions/useWPTitle';

defineProps<{
  open: boolean;
}>();

const emit = defineEmits<{
  (event: 'close'): void;
}>();

interface ManualInput {
  type: 'string' | 'choice' | 'boolean';
  description?: string;
  required?: boolean;
  default?: string | boolean;
  options: string[];
}

interface ManualInputsResponse {
  inputs: Record<string, ManualInput>;
}

const apiClient = useApiClient();
const config = useConfig();
const notifications = useNotification();
const i18n = useI18n();

const repo = requiredInject('repo');
const repoPermissions = requiredInject('repo-permissions');

const router = useRouter();
const branches = ref<{ text: string; value: string }[]>([]);
const payload = ref<{ branch: string; variables: Record<string, string> }>({
  branch: 'main',
  variables: {},
});
const manualInputs = ref<Record<string, ManualInput>>({});
const typedValues = ref<Record<string, string | boolean>>({});
const loadingManualInputs = ref(false);
const initialized = ref(false);
let manualInputsRequest = 0;

const manualInputEntries = computed(() => Object.entries(manualInputs.value));
const isVariablesValid = ref(true);

const areRequiredInputsValid = computed(() =>
  manualInputEntries.value.every(([name, input]) => {
    if (!input.required || input.type === 'boolean') return true;
    return String(typedValues.value[name] ?? '').trim() !== '';
  }),
);

const isFormValid = computed(() => {
  return payload.value.branch !== '' && isVariablesValid.value && areRequiredInputsValid.value && !loadingManualInputs.value;
});

const pipelineOptions = computed(() => {
  const typedVariables: Record<string, string> = {};
  for (const [name, input] of manualInputEntries.value) {
    const value = typedValues.value[name];
    if (input.type === 'boolean') {
      typedVariables[name] = String(value ?? false);
    } else if (value !== undefined && value !== '') {
      typedVariables[name] = String(value);
    }
  }

  return {
    ...payload.value,
    variables: {
      ...payload.value.variables,
      ...typedVariables,
    },
  };
});

const loading = ref(true);
onMounted(async () => {
  if (!repoPermissions.value.push) {
    notifications.notify({ type: 'error', title: i18n.t('repo.settings.not_allowed') });
    await router.replace({ name: 'home' });
    return;
  }

  const data = await usePaginate((page) => apiClient.getRepoBranches(repo.value.id, { page }));
  branches.value = data.map((e) => ({
    text: e,
    value: e,
  }));
  payload.value.branch = repo.value.branch || data[0] || 'main';
  initialized.value = true;
  await loadManualInputs(payload.value.branch);
  loading.value = false;
});

watch(
  () => payload.value.branch,
  async (branch) => {
    if (initialized.value && branch) {
      await loadManualInputs(branch);
    }
  },
);

function formatInputName(name: string): string {
  return name
    .replaceAll('_', ' ')
    .replaceAll('-', ' ')
    .replace(/\b\w/g, (character) => character.toUpperCase());
}

async function loadManualInputs(branch: string) {
  const request = ++manualInputsRequest;
  loadingManualInputs.value = true;

  try {
    const rootPath = config.rootPath.replace(/\/$/, '');
    const response = await fetch(
      `${rootPath}/api/repos/${repo.value.id}/manual-pipeline-inputs?branch=${encodeURIComponent(branch)}`,
      { credentials: 'same-origin' },
    );
    if (!response.ok) {
      throw new Error(await response.text());
    }

    const data = (await response.json()) as ManualInputsResponse;
    if (request !== manualInputsRequest) return;

    manualInputs.value = data.inputs ?? {};
    const values: Record<string, string | boolean> = {};
    for (const [name, input] of Object.entries(manualInputs.value)) {
      if (input.type === 'boolean') {
        values[name] = typeof input.default === 'boolean' ? input.default : false;
      } else {
        values[name] = typeof input.default === 'string' ? input.default : '';
      }
    }
    typedValues.value = values;
  } catch (error) {
    if (request !== manualInputsRequest) return;
    manualInputs.value = {};
    typedValues.value = {};
    notifications.notify({ type: 'error', title: String(error) });
  } finally {
    if (request === manualInputsRequest) {
      loadingManualInputs.value = false;
    }
  }
}

async function triggerManualPipeline() {
  loading.value = true;
  const pipeline = await apiClient.createPipeline(repo.value.id, pipelineOptions.value);

  emit('close');

  if (typeof pipeline == 'string') {
    // if this is a string (http 204) there is no workflow to run with the 'manual' event
    await router.push({
      name: 'repo',
    });

    notifications.notify({ type: 'warn', title: i18n.t('repo.manual_pipeline.no_manual_workflows') });
  } else {
    await router.push({
      name: 'repo-pipeline',
      params: {
        pipelineId: pipeline.number,
      },
    });
  }

  loading.value = false;
}

useWPTitle(computed(() => [i18n.t('repo.manual_pipeline.trigger'), repo.value.full_name]));
</script>
