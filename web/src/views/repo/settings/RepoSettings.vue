<template>
  <Scaffold enable-tabs :go-back="goBack">
    <template #title>
      <span>
        <router-link :to="{ name: 'org', params: { orgId: repo.org_id } }" class="hover:underline">{{ repo!.owner }}</router-link>
        /
        <router-link :to="{ name: 'repo' }" class="hover:underline">{{ repo!.name }}</router-link>
        /
        {{ $t('settings') }}
      </span>
    </template>

    <Tab icon="settings-outline" :to="{ name: 'repo-settings' }" :title="$t('repo.settings.general.general')" />
    <Tab icon="secret" :to="{ name: 'repo-settings-secrets' }" :title="$t('secrets.secrets')" />
    <Tab icon="docker" :to="{ name: 'repo-settings-registries' }" :title="$t('registries.registries')" />
    <Tab icon="cron" :to="{ name: 'repo-settings-crons' }" :title="$t('repo.settings.crons.crons')" />
    <Tab icon="tag" :to="{ name: 'repo-settings-badge' }" :title="$t('repo.settings.badge.badge')" />
    <Tab icon="puzzle" :to="{ name: 'repo-settings-extensions' }" :title="$t('extensions')" />
    <Tab icon="toolbox" :to="{ name: 'repo-settings-actions' }" :title="$t('repo.settings.actions.actions')" />

    <div class="lha-ci-shell-intro">
      <div class="lha-ci-shell-intro__copy">
        <p class="lha-ci-kicker">Repository configuration</p>
        <h2 class="lha-ci-shell-intro__title">Build and delivery settings</h2>
        <p class="lha-ci-shell-intro__text">Manage execution, secrets, registries, schedules, extensions and repository-level CI behavior.</p>
      </div>
      <div class="lha-ci-shell-intro__context">
        <span>Repository</span>
        <strong>{{ repo.full_name }}</strong>
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
import { requiredInject } from '~/compositions/useInjectProvide';
import useNotifications from '~/compositions/useNotifications';
import { useRouteBack } from '~/compositions/useRouteBack';

const notifications = useNotifications();
const router = useRouter();
const i18n = useI18n();
const repoPermissions = requiredInject('repo-permissions');
const repo = requiredInject('repo');

onMounted(async () => {
  if (!repoPermissions.value.admin) {
    notifications.notify({ type: 'error', title: i18n.t('repo.settings.not_allowed') });
    await router.replace({ name: 'home' });
  }
});

const goBack = useRouteBack({ name: 'repo' });
</script>
