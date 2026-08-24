<template>
  <Panel v-if="!loading" class="lha-ci-manual-panel">
    <form class="lha-ci-manual" @submit.prevent="triggerManualPipeline">
      <header class="lha-ci-manual__header">
        <div>
          <p class="lha-ci-kicker">LHA Play CI</p>
          <h1 class="lha-ci-manual__title">{{ $t('repo.manual_pipeline.title') }}</h1>
          <p class="lha-ci-manual__intro">
            Choose the source branch and only the inputs this pipeline actually supports. Advanced variables are optional.
          </p>
        </div>
        <div class="lha-ci-manual__repo">
          <span>Repository</span>
          <strong>{{ repo.full_name }}</strong>
        </div>
      </header>

      <section class="lha-ci-manual__section">
        <div class="lha-ci-manual__section-heading">
          <div>
            <span class="lha-ci-manual__step">01</span>
            <h2>Source</h2>
          </div>
          <span class="lha-ci-manual__hint">Select the commit source for this run.</span>
        </div>

        <InputField v-slot="{ id }" :label="$t('repo.manual_pipeline.select_branch')">
          <SelectField :id="id" v-model="payload.branch" :options="branches" required />
        </InputField>
      </section>

      <section v-if="manualInputEntries.length" class="lha-ci-manual__section">
        <div class="lha-ci-manual__section-heading">
          <div>
            <span class="lha-ci-manual__step">02</span>
            <h2>Build options</h2>
          </div>
          <span class="lha-ci-manual__hint">{{ manualInputEntries.length }} pipeline input{{ manualInputEntries.length === 1 ? '' : 's' }}</span>
        </div>

        <div class="lha-ci-manual__inputs">
          <template v-for="[name, input] in manualInputEntries" :key="name">
            <div class="lha-ci-input-card" :class="{ 'lha-ci-input-card--boolean': input.type === 'boolean' }">
              <Checkbox
                v-if="input.type === 'boolean'"
                :model-value="getBooleanValue(name)"
                :label="formatInputName(name)"
                :description="input.description"
                @update:model-value="setTypedValue(name, $event)"
              />

              <InputField v-else v-slot="{ id }" :label="formatInputName(name)">
                <span v-if="input.description" class="text-wp-text-alt-100 mb-2 text-sm">{{ input.description }}</span>
                <SelectField
                  v-if="input.type === 'choice'"
                  :id="id"
                  :model-value="getStringValue(name)"
                  :options="input.options.map((option) => ({ text: option, value: option }))"
                  :placeholder="getStringValue(name) === '' ? '-' : undefined"
                  @update:model-value="setTypedValue(name, $event)"
                />
                <TextField
                  v-else
                  :id="id"
                  :model-value="getStringValue(name)"
                  @update:model-value="setTypedValue(name, $event)"
                />
              </InputField>
            </div>
          </template>
        </div>
      </section>

      <section class="lha-ci-manual__section lha-ci-manual__summary">
        <div class="lha-ci-manual__section-heading">
          <div>
            <span class="lha-ci-manual__step">03</span>
            <h2>Run summary</h2>
          </div>
        </div>

        <div class="lha-ci-run-summary">
          <div>
            <span>Branch</span>
            <strong>{{ payload.branch }}</strong>
          </div>
          <div>
            <span>Configured inputs</span>
            <strong>{{ configuredInputCount }}</strong>
          </div>
          <div>
            <span>Action</span>
            <strong>Start pipeline</strong>
          </div>
        </div>
      </section>

      <details class="lha-ci-advanced">
        <summary>Advanced variables</summary>
        <div class="lha-ci-advanced__body">
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
        </div>
      </details>

      <footer class="lha-ci-manual__actions">
        <div class="lha-ci-manual__validation">
          <span v-if="loadingManualInputs">Loading pipeline inputs…</span>
          <span v-else-if="!areRequiredInputsValid">Complete the required build options before starting.</span>
          <span v-else>Ready to run {{ payload.branch }}.</span>
        </div>
        <Button type="submit" :text="$t('repo.manual_pipeline.trigger')" :disabled="!isFormValid" />
      </footer>
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

const manualInputEntries = computed(() =>
  Object.entries(manualInputs.value).sort(([left], [right]) => left.localeCompare(right)),
);
const isVariablesValid = ref(true);

const areRequiredInputsValid = computed(() =>
  manualInputEntries.value.every(([name, input]) => {
    if (!input.required || input.type === 'boolean') return true;
    return getStringValue(name).trim() !== '';
  }),
);

const configuredInputCount = computed(() => {
  let count = Object.keys(payload.value.variables).length;
  for (const [name, input] of manualInputEntries.value) {
    if (input.type === 'boolean') {
      if (getBooleanValue(name)) count += 1;
      continue;
    }
    if (getStringValue(name).trim() !== '') count += 1;
  }
  return count;
});

const isFormValid = computed(() => {
  return payload.value.branch !== '' && isVariablesValid.value && areRequiredInputsValid.value && !loadingManualInputs.value;
});

const pipelineOptions = computed(() => {
  const typedVariables: Record<string, string> = {};
  for (const [name, input] of manualInputEntries.value) {
    if (input.type === 'boolean') {
      typedVariables[name] = String(getBooleanValue(name));
      continue;
    }

    const value = getStringValue(name);
    if (value !== '') {
      typedVariables[name] = value;
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
  payload.value.branch = repo.value.default_branch || data[0] || 'main';
  await loadManualInputs(payload.value.branch);
  initialized.value = true;
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

function getStringValue(name: string): string {
  const value = typedValues.value[name];
  return typeof value === 'string' ? value : '';
}

function getBooleanValue(name: string): boolean {
  return typedValues.value[name] === true;
}

function setTypedValue(name: string, value: string | boolean) {
  typedValues.value[name] = value;
}

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
