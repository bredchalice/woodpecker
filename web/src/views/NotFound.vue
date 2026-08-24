<template>
  <main class="lha-ci-state-page">
    <section class="lha-ci-state-card">
      <p class="lha-ci-kicker">LHA Play CI</p>
      <div class="lha-ci-state-code">404</div>
      <h1 class="lha-ci-state-title">{{ $t('not_found.not_found') }}</h1>
      <router-link class="lha-ci-state-link" replace :to="{ name: 'home' }">
        {{ $t('not_found.back_home') }}
      </router-link>
    </section>
  </main>
</template>
