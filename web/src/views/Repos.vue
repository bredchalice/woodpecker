<template>
  <Scaffold v-model:search="search">
    <template #title>
      {{ $t('repositories.title') }}
    </template>

    <template #headerActions>
      <Button :to="{ name: 'repo-add' }" start-icon="plus" :text="$t('repo.add')" />
      <Button start-icon="refresh" :is-loading="isRefreshing" :text="$t('repo.refresh')" @click="refreshRepositories" />
    </template>

    <div class="lha-ci-repos">
      <section v-if="search === ''" class="lha-ci-repos__hero">
        <div>
          <p class="lha-ci-kicker">LHA Play CI</p>
          <h1 class="lha-ci-repos__title">Delivery workspace</h1>
          <p class="lha-ci-repos__intro">Follow active repositories, recent builds and release activity from one operational view.</p>
        </div>
        <div class="lha-ci-repos__stats">
          <div>
            <span>Repositories</span>
            <strong>{{ repos.length }}</strong>
          </div>
          <div>
            <span>Recent</span>
            <strong>{{ reposLastAccess.length }}</strong>
          </div>
          <div>
            <span>Visible</span>
            <strong>{{ reposLastActivity.length }}</strong>
          </div>
        </div>
      </section>

      <Transition name="fade" mode="out-in">
        <div v-if="search === '' && repos.length > 0" class="lha-ci-repos__sections">
          <section v-if="reposLastAccess.length > 0 && repos.length > 4" class="lha-ci-repos__section">
            <div class="lha-ci-section-heading">
              <div>
                <p class="lha-ci-kicker">Resume work</p>
                <h2>{{ $t('repositories.last.title') }}</h2>
                <span>{{ $t('repositories.last.desc') }}</span>
              </div>
            </div>
            <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
              <RepoItem v-for="repo in reposLastAccess" :key="repo.id" :repo="repo" />
            </div>
          </section>

          <section class="lha-ci-repos__section">
            <div class="lha-ci-section-heading">
              <div>
                <p class="lha-ci-kicker">Build activity</p>
                <h2>{{ $t('repositories.all.title') }}</h2>
                <span>{{ $t('repositories.all.desc') }}</span>
              </div>
            </div>
            <div class="flex flex-col gap-3">
              <RepoItem v-for="repo in reposLastActivity" :key="repo.id" :repo="repo" />
            </div>
          </section>
        </div>
        <section v-else class="lha-ci-repos__section">
          <div v-if="reposLastActivity.length > 0" class="flex flex-col gap-3">
            <RepoItem v-for="repo in reposLastActivity" :key="repo.id" :repo="repo" />
          </div>
          <div v-else class="lha-ci-empty-state">
            <p class="lha-ci-kicker">Search</p>
            <strong>{{ $t('no_search_results') }}</strong>
          </div>
        </section>
      </Transition>
    </div>
  </Scaffold>
</template>

<script lang="ts" setup>
import { computed, onMounted, ref } from 'vue';
import { useI18n } from 'vue-i18n';

import Button from '~/components/atomic/Button.vue';
import Scaffold from '~/components/layout/scaffold/Scaffold.vue';
import RepoItem from '~/components/repo/RepoItem.vue';
import { useAsyncAction } from '~/compositions/useAsyncAction';
import useRepos from '~/compositions/useRepos';
import { useRepoSearch } from '~/compositions/useRepoSearch';
import { useWPTitle } from '~/compositions/useWPTitle';
import { useRepoStore } from '~/store/repos';

const repoStore = useRepoStore();

const { sortReposByLastAccess, sortReposByLastActivity, repoWithLastPipeline } = useRepos();
const repos = computed(() => Object.values(repoStore.ownedRepos).map((r) => repoWithLastPipeline(r)));

const reposLastAccess = computed(() => sortReposByLastAccess(repos.value || []).slice(0, 4));

const search = ref('');
const { searchedRepos } = useRepoSearch(repos, search);
const reposLastActivity = computed(() => sortReposByLastActivity(searchedRepos.value || []));

const { doSubmit: refreshRepositories, isLoading: isRefreshing } = useAsyncAction(async () => {
  await repoStore.refreshRepos();
  await repoStore.loadRepos();
});

onMounted(async () => {
  await repoStore.loadRepos();
});

const { t } = useI18n();
useWPTitle(computed(() => [t('repositories.title')]));
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: all 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
