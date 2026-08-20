<template>
  <Scaffold v-model:search="search" :go-back="goBack">
    <template #title>
      {{ $t('repo.add') }}
    </template>

    <div class="lha-ci-enable">
      <section class="lha-ci-enable__hero">
        <div>
          <p class="lha-ci-kicker">Repository setup</p>
          <h1 class="lha-ci-enable__title">Connect a repository</h1>
          <p class="lha-ci-enable__intro">Enable a source repository to make its pipelines available in the LHA Play CI workspace.</p>
        </div>
        <div class="lha-ci-enable__summary">
          <span>Available</span>
          <strong>{{ repos?.length ?? 0 }}</strong>
        </div>
      </section>

      <section class="lha-ci-enable__list">
        <template v-if="repos !== undefined && repos.length > 0">
          <template v-for="repo in searchedRepos" :key="repo.forge_remote_id">
            <div v-if="repo.has_forge_name_conflict" class="lha-ci-enable__conflict">
              <ListItem class="lha-ci-enable__item items-center rounded-b-none! border-b-0!">
                <div class="flex w-full items-center gap-3">
                  <div class="min-w-0">
                    <strong class="lha-ci-enable__repo-name">{{ repo.full_name }}</strong>
                    <span class="lha-ci-enable__meta">{{ $t('repo.enable.new_forge_repo') }}</span>
                  </div>
                  <div class="ml-auto flex items-center">
                    <Button :text="$t('repo.enable.conflict')" :title="$t('repo.enable.conflict_desc')" disabled />
                  </div>
                </div>
              </ListItem>

              <ListItem
                :to="{ name: 'repo', params: { repoId: repo.id } }"
                class="lha-ci-enable__item items-center rounded-t-none! border-t-0! opacity-80"
              >
                <div class="min-w-0">
                  <span class="lha-ci-enable__repo-name">{{ repo.full_name }}</span>
                  <span class="lha-ci-enable__meta">{{ $t('repo.enable.stale_wp_repo') }}</span>
                </div>
                <div class="ml-auto" @click.prevent.stop>
                  <Button
                    start-icon="toolbox"
                    :text="$t('repo.settings.actions.actions')"
                    :to="{ name: 'repo-settings-actions', params: { repoId: repo.id } }"
                  />
                </div>
              </ListItem>
            </div>

            <ListItem
              v-else-if="repo.has_no_forge_repo"
              :to="{ name: 'repo', params: { repoId: repo.id } }"
              class="lha-ci-enable__item items-center"
            >
              <span class="lha-ci-enable__repo-name">{{ repo.full_name }}</span>
              <span class="lha-ci-enable__meta ml-auto">{{ $t('repo.enable.forge_repo_missing') }}</span>
            </ListItem>

            <ListItem
              v-else-if="repo.active"
              :to="{ name: 'repo', params: { repoId: repo.id } }"
              class="lha-ci-enable__item items-center"
            >
              <span class="lha-ci-enable__repo-name">{{ repo.full_name }}</span>
              <span class="lha-ci-enable__state ml-auto">{{ $t('repo.enable.enabled') }}</span>
            </ListItem>

            <ListItem
              v-else
              class="lha-ci-enable__item items-center"
              :to="repo.id ? { name: 'repo', params: { repoId: repo.id } } : undefined"
            >
              <span class="lha-ci-enable__repo-name">{{ repo.full_name }}</span>
              <div class="ml-auto flex items-center">
                <Badge v-if="repo.id" class="md:display-unset mr-2 hidden" :value="$t('repo.enable.disabled')" />
                <Button
                  :text="$t('repo.enable.enable')"
                  :is-loading="isActivatingRepo && repoToActivate?.forge_remote_id === repo.forge_remote_id"
                  @click="activateRepo(repo)"
                />
              </div>
            </ListItem>
          </template>
        </template>
        <div v-else-if="loading" class="lha-ci-empty-state">
          <Icon name="spinner" />
        </div>
        <div v-else class="lha-ci-empty-state">
          <strong>{{ $t('no_search_results') }}</strong>
        </div>
      </section>
    </div>
  </Scaffold>
</template>

<script lang="ts" setup>
import { computed, onMounted, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';

import Badge from '~/components/atomic/Badge.vue';
import Button from '~/components/atomic/Button.vue';
import Icon from '~/components/atomic/Icon.vue';
import ListItem from '~/components/atomic/ListItem.vue';
import Scaffold from '~/components/layout/scaffold/Scaffold.vue';
import useApiClient from '~/compositions/useApiClient';
import { useAsyncAction } from '~/compositions/useAsyncAction';
import useNotifications from '~/compositions/useNotifications';
import { useRepoSearch } from '~/compositions/useRepoSearch';
import { useRouteBack } from '~/compositions/useRouteBack';
import { useWPTitle } from '~/compositions/useWPTitle';
import type { Repo } from '~/lib/api/types';

const router = useRouter();
const apiClient = useApiClient();
const notifications = useNotifications();
const repos = ref<Repo[]>();
const repoToActivate = ref<Repo>();
const search = ref('');
const i18n = useI18n();
const loading = ref(false);

const { searchedRepos } = useRepoSearch(repos, search);

onMounted(async () => {
  loading.value = true;
  repos.value = await apiClient.getRepoList({ all: true });
  loading.value = false;
});

const { doSubmit: activateRepo, isLoading: isActivatingRepo } = useAsyncAction(async (repo: Repo) => {
  repoToActivate.value = repo;
  const _repo = await apiClient.activateRepo(repo.forge_remote_id);
  notifications.notify({ title: i18n.t('repo.enable.success'), type: 'success' });
  repoToActivate.value = undefined;
  await router.push({ name: 'repo', params: { repoId: _repo.id } });
});

const goBack = useRouteBack({ name: 'repos' });

useWPTitle(computed(() => [i18n.t('repo.add')]));
</script>
