<template>
  <Scaffold v-if="org" enable-tabs :go-back="goBack">
    <template #title>
      <span><router-link :to="{ name: 'org' }" class="hover:underline">{{ org.name }}</router-link> / {{ $t('settings') }}</span>
    </template>

    <Tab icon="secret" :to="{ name: 'org-settings-secrets' }" :title="$t('secrets.secrets')" />
    <Tab icon="docker" :to="{ name: 'org-settings-registries' }" :title="$t('registries.registries')" />
    <Tab v-if="userRegisteredAgents" icon="agent" :to="{ name: 'org-settings-agents' }" :title="$t('admin.settings.agents.agents')" />

    <div class="lha-ci-shell-intro">
      <div class="lha-ci-shell-intro__copy">
        <p class="lha-ci-kicker">Organization scope</p>
        <h2 class="lha-ci-shell-intro__title">Shared CI configuration</h2>
        <p class="lha-ci-shell-intro__text">Manage credentials, registries and build agents shared by repositories in this organization.</p>
      </div>
      <div class="lha-ci-shell-intro__context">
        <span>Organization</span>
        <strong>{{ org.name }}</strong>
      </div>
    </div>

    <router-view />
  </Scaffold>
</template>

<script lang="ts" setup>
import { onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';

import Scaffold from '~/components/layout/scaffold/Scaffold.vue';
import Tab from '~/components/layout/scaffold/Tab.vue';
import useConfig from '~/compositions/useConfig';
import { requiredInject } from '~/compositions/useInjectProvide';
import useNotifications from '~/compositions/useNotifications';
import { useRouteBack } from '~/compositions/useRouteBack';

const notifications = useNotifications();
const router = useRouter();
const i18n = useI18n();
const { userRegisteredAgents } = useConfig();
const org = requiredInject('org');
const orgPermissions = requiredInject('org-permissions');

onMounted(async () => {
  if (!orgPermissions.value?.admin) {
    notifications.notify({ type: 'error', title: i18n.t('org.settings.not_allowed') });
    await router.replace({ name: 'home' });
  }
});

const goBack = useRouteBack({ name: 'org' });
</script>
