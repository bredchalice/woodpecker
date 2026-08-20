<template>
  <nav class="lha-ci-navbar flex border-b px-4 py-3 font-bold">
    <div class="flex min-w-0 items-center gap-2">
      <router-link :to="{ name: 'home' }" class="lha-ci-brand flex min-w-0 items-center gap-3 pr-2">
        <span class="lha-ci-brand__mark">
          <WoodpeckerLogo class="h-7 w-7 shrink-0" />
        </span>
        <span class="hidden min-w-0 flex-col sm:flex">
          <strong class="lha-ci-brand__name">LHA Play CI</strong>
          <span class="lha-ci-brand__version" :title="version?.current">{{ version?.current }}</span>
        </span>
      </router-link>

      <div class="lha-ci-navbar__divider hidden h-7 w-px sm:block" />

      <router-link v-if="user" :to="{ name: 'repos' }" class="navbar-clickable navbar-link">
        <span class="flex md:hidden">{{ $t('repos') }}</span>
        <span class="hidden md:flex">{{ $t('repositories.title') }}</span>
      </router-link>
      <a href="https://woodpecker-ci.org/" target="_blank" class="navbar-clickable navbar-link hidden md:flex">{{ $t('docs') }}</a>
      <a v-if="enableSwagger" :href="apiUrl" target="_blank" class="navbar-clickable navbar-link hidden md:flex">{{ $t('api') }}</a>
    </div>

    <div class="ml-auto flex items-center gap-1.5">
      <IconButton
        v-if="user?.admin"
        class="navbar-icon relative"
        :title="$t('settings')"
        :to="{ name: 'admin-settings' }"
      >
        <Icon name="settings" />
        <div v-if="version?.needsUpdate" class="bg-wp-error-100 absolute top-2 right-2 h-2.5 w-2.5 rounded-full" />
      </IconButton>

      <ActivePipelines v-if="user" class="navbar-icon p-1.5!" />
      <IconButton v-if="user" :to="{ name: 'user' }" :title="$t('user.settings.settings')" class="navbar-icon p-1.5!">
        <img v-if="user.avatar_url" class="rounded-md" :src="user.avatar_url" />
      </IconButton>
      <Button v-else :text="$t('login')" :to="{ name: 'login' }" class="navbar-login" @click="saveRedirect" />
    </div>
  </nav>
</template>

<script lang="ts" setup>
import { useRoute } from 'vue-router';

import WoodpeckerLogo from '~/assets/logo.svg?component';
import Button from '~/components/atomic/Button.vue';
import Icon from '~/components/atomic/Icon.vue';
import IconButton from '~/components/atomic/IconButton.vue';
import useAuthentication from '~/compositions/useAuthentication';
import useConfig from '~/compositions/useConfig';
import useUserConfig from '~/compositions/useUserConfig';
import { useVersion } from '~/compositions/useVersion';

import ActivePipelines from './ActivePipelines.vue';

const version = useVersion();
const config = useConfig();
const userConfig = useUserConfig();
const route = useRoute();
const { user } = useAuthentication();
const apiUrl = `${config.rootPath ?? ''}/swagger/index.html`;
const { enableSwagger } = config;

function saveRedirect() {
  userConfig.setUserConfig('redirectUrl', route.fullPath);
}
</script>

<style scoped>
.lha-ci-navbar {
  position: relative;
  z-index: 30;
  border-color: var(--lha-ci-border);
  background: color-mix(in srgb, var(--lha-ci-surface) 94%, var(--lha-ci-accent) 6%);
  color: var(--wp-text-200);
  box-shadow: 0 1px 0 rgba(0, 0, 0, 0.03);
}

.lha-ci-navbar::after {
  position: absolute;
  right: 0;
  bottom: -1px;
  left: 0;
  height: 1px;
  background: linear-gradient(90deg, transparent, var(--lha-ci-accent), transparent);
  content: '';
  opacity: 0.45;
}

.lha-ci-brand {
  color: var(--wp-text-200);
  text-decoration: none;
}

.lha-ci-brand__mark {
  display: flex;
  width: 2.5rem;
  height: 2.5rem;
  align-items: center;
  justify-content: center;
  border: 1px solid color-mix(in srgb, var(--lha-ci-accent) 28%, var(--lha-ci-border));
  border-radius: 0.75rem;
  background: var(--lha-ci-accent-soft);
}

.lha-ci-brand__name {
  line-height: 1.05;
  letter-spacing: -0.025em;
}

.lha-ci-brand__version {
  max-width: 13rem;
  margin-top: 0.18rem;
  overflow: hidden;
  color: var(--wp-text-alt-100);
  font-size: 0.62rem;
  font-weight: 650;
  line-height: 1;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.lha-ci-navbar__divider {
  background: var(--lha-ci-border);
}

.navbar-icon {
  width: 2.5rem;
  height: 2.5rem;
  padding: 0.55rem;
  border-radius: 0.7rem;
  color: var(--wp-text-100);
}

.navbar-icon:hover,
.navbar-link:hover {
  background: var(--lha-ci-accent-soft);
  color: var(--wp-text-200);
}

.navbar-icon :deep(svg) {
  width: 100%;
  height: 100%;
}

.navbar-link {
  margin-block: -0.25rem;
  padding: 0.6rem 0.75rem;
  border-radius: 0.65rem;
  color: var(--wp-text-100);
  transition: background 140ms ease, color 140ms ease;
}

.navbar-login {
  border-color: var(--lha-ci-accent) !important;
  background: var(--lha-ci-accent-soft) !important;
  color: var(--wp-text-200) !important;
}
</style>
