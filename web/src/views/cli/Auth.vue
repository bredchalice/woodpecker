<template>
  <main class="lha-ci-state-page">
    <section class="lha-ci-state-card lha-ci-cli-auth">
      <div class="lha-ci-cli-auth__mark">
        <WoodpeckerLogo class="h-12 w-12" />
      </div>
      <p class="lha-ci-kicker">LHA Play CI</p>

      <template v-if="state === 'confirm'">
        <h1 class="lha-ci-state-title">{{ $t('login_to_cli') }}</h1>
        <p class="lha-ci-state-copy">{{ $t('login_to_cli_description') }}</p>
        <div class="lha-ci-cli-auth__actions">
          <Button :text="$t('login_to_cli')" color="green" :is-loading="isSubmitting" @click="sendToken(false)" />
          <Button :text="$t('abort')" color="red" :disabled="isSubmitting" @click="abortLogin" />
        </div>
      </template>

      <template v-else-if="state === 'success'">
        <div class="lha-ci-state-code">OK</div>
        <h1 class="lha-ci-state-title">{{ $t('cli_login_success') }}</h1>
        <p class="lha-ci-state-copy">{{ $t('return_to_cli') }}</p>
      </template>

      <template v-else-if="state === 'failed'">
        <div class="lha-ci-state-code">!</div>
        <h1 class="lha-ci-state-title">{{ $t('cli_login_failed') }}</h1>
        <p class="lha-ci-state-copy">{{ $t('return_to_cli') }}</p>
      </template>

      <template v-else-if="state === 'denied'">
        <div class="lha-ci-state-code">×</div>
        <h1 class="lha-ci-state-title">{{ $t('cli_login_denied') }}</h1>
        <p class="lha-ci-state-copy">{{ $t('return_to_cli') }}</p>
      </template>
    </section>
  </main>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRoute } from 'vue-router';

import WoodpeckerLogo from '~/assets/logo.svg?component';
import Button from '~/components/atomic/Button.vue';
import useApiClient from '~/compositions/useApiClient';

const apiClient = useApiClient();
const route = useRoute();
const state = ref<'confirm' | 'success' | 'failed' | 'denied'>('confirm');
const isSubmitting = ref(false);

async function sendToken(abort = false) {
  const port = route.query.port as string;
  if (!port) {
    state.value = 'failed';
    return;
  }

  isSubmitting.value = true;
  try {
    const address = `http://localhost:${port}`;
    const token = abort ? '' : await apiClient.getToken();
    const response = await fetch(`${address}/token`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ token }),
    });

    if (abort) {
      state.value = 'denied';
      window.close();
      return;
    }

    const data = (await response.json()) as { ok: string };
    state.value = data.ok === 'true' ? 'success' : 'failed';
  } catch {
    state.value = abort ? 'denied' : 'failed';
  } finally {
    isSubmitting.value = false;
  }
}

async function abortLogin() {
  await sendToken(true);
}
</script>
