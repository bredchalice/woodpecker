<template>
  <Scaffold enable-tabs class="lha-ci-nav-shell">
    <template #title>{{ $t('user.settings.settings') }}</template>
    <template #headerActions><Button :text="$t('logout')" :to="`${address}/logout`" /></template>

    <Tab icon="settings-outline" :to="{ name: 'user' }" :title="$t('user.settings.general.general')" />
    <Tab icon="secret" :to="{ name: 'user-secrets' }" :title="$t('secrets.secrets')" />
    <Tab icon="docker" :to="{ name: 'user-registries' }" :title="$t('registries.registries')" />
    <Tab icon="console" :to="{ name: 'user-cli-and-api' }" :title="$t('user.settings.cli_and_api.cli_and_api')" />
    <Tab v-if="userRegisteredAgents" icon="agent" :to="{ name: 'user-agents' }" :title="$t('admin.settings.agents.agents')" />

    <div class="lha-ci-shell-intro">
      <div class="lha-ci-shell-intro__copy">
        <p class="lha-ci-kicker">Developer profile</p>
        <h2 class="lha-ci-shell-intro__title">Personal CI configuration</h2>
        <p class="lha-ci-shell-intro__text">Manage your account, personal secrets, registry credentials, CLI access and registered build agents.</p>
      </div>
      <div class="lha-ci-shell-intro__context">
        <span>Product</span>
        <strong>LHA Play CI</strong>
      </div>
    </div>

    <router-view />
  </Scaffold>
</template>

<script lang="ts" setup>
import Button from '~/components/atomic/Button.vue';
import Scaffold from '~/components/layout/scaffold/Scaffold.vue';
import Tab from '~/components/layout/scaffold/Tab.vue';
import useConfig from '~/compositions/useConfig';

const { userRegisteredAgents } = useConfig();
const address = `${window.location.protocol}//${window.location.host}${useConfig().rootPath}`;
</script>
