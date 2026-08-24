<template>
  <main class="lha-ci-login">
    <Error v-if="errorMessage" class="lha-ci-login__error">
      <span class="whitespace-pre">{{ errorMessage }}</span>
      <span v-if="errorDescription" class="mt-1 whitespace-pre">{{ errorDescription }}</span>
      <a
        v-if="errorUri"
        :href="errorUri"
        target="_blank"
        class="text-wp-link-100 hover:text-wp-link-200 mt-1 cursor-pointer"
      >
        <span>{{ errorUri }}</span>
      </a>
    </Error>

    <section class="lha-ci-login__shell">
      <div class="lha-ci-login__brand">
        <div class="lha-ci-login__mark">
          <WoodpeckerLogo preserveAspectRatio="xMidYMid meet" class="h-24 w-24" />
        </div>
        <div>
          <p class="lha-ci-kicker">LHA Play</p>
          <h1 class="lha-ci-login__title">CI</h1>
          <p class="lha-ci-login__intro">Build, verify and release LHA Play from one focused delivery workspace.</p>
        </div>
        <div class="lha-ci-login__signals" aria-hidden="true">
          <span>BUILD</span>
          <span>VERIFY</span>
          <span>RELEASE</span>
        </div>
      </div>

      <div class="lha-ci-login__auth">
        <div>
          <p class="lha-ci-kicker">Authentication</p>
          <h2 class="lha-ci-login__auth-title">{{ $t('login_to_woodpecker_with') }}</h2>
          <p class="lha-ci-login__auth-copy">Use your connected source provider to enter the CI workspace.</p>
        </div>

        <div class="lha-ci-login__providers">
          <Button
            v-for="forge in forgesWithNameAndFavicon"
            :key="forge.id"
            :start-icon="forge.type === 'addon' ? 'repo' : forge.type"
            class="lha-ci-login__provider whitespace-normal!"
            @click="authenticate(forge.id)"
          >
            <div class="mr-2 w-4">
              <img
                v-if="forge.favicon && !failedForgeFavicons.has(forge.id)"
                :src="forge.favicon"
                :alt="$t('login_to_woodpecker_with', { forge: forge.name })"
                @error="() => failedForgeFavicons.add(forge.id)"
              />
              <Icon v-else :name="forge.type === 'addon' ? 'repo' : forge.type" />
            </div>

            {{ forge.name }}
          </Button>
        </div>

        <div class="lha-ci-login__footnote">
          <span class="lha-ci-login__status-dot" />
          <span>LHA Play CI</span>
        </div>
      </div>
    </section>
  </main>
</template>

<script lang="ts" setup>
import { computed, onMounted, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRoute } from 'vue-router';

import WoodpeckerLogo from '~/assets/logo.svg?component';
import Button from '~/components/atomic/Button.vue';
import Error from '~/components/atomic/Error.vue';
import Icon from '~/components/atomic/Icon.vue';
import useApiClient from '~/compositions/useApiClient';
import useAuthentication from '~/compositions/useAuthentication';
import { useWPTitle } from '~/compositions/useWPTitle';
import type { Forge } from '~/lib/api/types';

const route = useRoute();
const { authenticate } = useAuthentication();
const i18n = useI18n();
const apiClient = useApiClient();

const forges = ref<Forge[]>([]);

const authErrorMessages = {
  oauth_error: i18n.t('oauth_error'),
  internal_error: i18n.t('internal_error'),
  registration_closed: i18n.t('registration_closed'),
  access_denied: i18n.t('access_denied'),
  invalid_state: i18n.t('invalid_state'),
  org_access_denied: i18n.t('org_access_denied'),
};

const errorMessage = ref<string>();
const errorDescription = ref<string>(route.query.error_description as string);
const errorUri = ref<string>(route.query.error_uri as string);

onMounted(async () => {
  forges.value = (await apiClient.getForges()) ?? [];

  if (route.query.error) {
    const error = route.query.error as keyof typeof authErrorMessages;
    errorMessage.value = authErrorMessages[error] ?? error;
  }
});

useWPTitle(computed(() => [i18n.t('login')]));

const failedForgeFavicons = ref(new Set<number>());

const forgesWithNameAndFavicon = computed(() =>
  forges.value.map((forge) => {
    let name = forge.type.charAt(0).toUpperCase() + forge.type.slice(1);
    let favicon: null | string = null;

    if (forge.url || forge.oauth_host) {
      const url = new URL(forge.oauth_host || forge.url);
      name = url.hostname;
      favicon = `${url.origin}/favicon.ico`;
    }

    return {
      ...forge,
      name,
      favicon,
    };
  }),
);
</script>
