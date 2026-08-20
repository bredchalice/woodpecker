// cSpell:ignore tseslint
// @ts-check

import antfu from '@antfu/eslint-config';
import js from '@eslint/js';
import vueI18n from '@intlify/eslint-plugin-vue-i18n';
import eslintPromise from 'eslint-plugin-promise';
import eslintPluginVueScopedCSS from 'eslint-plugin-vue-scoped-css';

export default antfu(
  {
    stylistic: false,
    typescript: { tsconfigPath: './tsconfig.json' },
    vue: true,
    jsonc: false,
    yaml: false,
  },
  js.configs.recommended,
  eslintPromise.configs['flat/recommended'],
  ...eslintPluginVueScopedCSS.configs['flat/recommended'],
  ...vueI18n.configs['flat/recommended'],
  {
    rules: {
      'import/order': 'off',
      'sort-imports': 'off',
      'perfectionist/sort-imports': 'off',
      'perfectionist/sort-named-imports': 'off',
      'promise/prefer-await-to-callbacks': 'error',
      'vue-scoped-css/no-parsing-error': 'off',
      '@intlify/vue-i18n/no-raw-text': ['error', { attributes: { '/.+/': ['label'] } }],
      '@intlify/vue-i18n/key-format-style': ['error', 'snake_case'],
      '@intlify/vue-i18n/no-duplicate-keys-in-locale': 'error',
      '@intlify/vue-i18n/no-dynamic-keys': 'error',
      '@intlify/vue-i18n/no-deprecated-i18n-component': 'error',
      '@intlify/vue-i18n/no-deprecated-tc': 'error',
      '@intlify/vue-i18n/no-i18n-t-path-prop': 'error',
      '@intlify/vue-i18n/no-missing-keys-in-other-locales': 'off',
      '@intlify/vue-i18n/valid-message-syntax': 'error',
      '@intlify/vue-i18n/no-missing-keys': 'error',
      '@intlify/vue-i18n/no-unknown-locale': 'error',
      '@intlify/vue-i18n/no-unused-keys': ['error', { extensions: ['.ts', '.vue'] }],
      '@intlify/vue-i18n/prefer-sfc-lang-attr': 'error',
      '@intlify/vue-i18n/no-html-messages': 'error',
      '@intlify/vue-i18n/prefer-linked-key-with-paren': 'error',
      '@intlify/vue-i18n/sfc-locale-attr': 'error',
    },
    settings: { 'vue-i18n': { localeDir: './src/assets/locales/en.json', messageSyntaxVersion: '^9.0.0' } },
  },
  {
    files: ['**/*.vue'],
    rules: {
      'vue/multi-word-component-names': 'off',
      'vue/html-self-closing': ['error', { html: { void: 'always', normal: 'always', component: 'always' }, svg: 'always', math: 'always' }],
      'vue/html-indent': 'off',
      'vue/block-order': ['error', { order: ['template', 'script', 'style'] }],
      'vue/singleline-html-element-content-newline': ['off'],
      'no-useless-assignment': ['off'],
    },
  },
  {
    files: [
      'src/components/layout/header/Navbar.vue',
      'src/components/repo/RepoItem.vue',
      'src/components/repo/pipeline/PipelineItem.vue',
      'src/views/Login.vue',
      'src/views/RepoAdd.vue',
      'src/views/Repos.vue',
      'src/views/user/UserWrapper.vue',
      'src/views/admin/AdminSettingsWrapper.vue',
      'src/views/admin/AdminAgents.vue',
      'src/views/admin/AdminQueue.vue',
      'src/views/org/settings/OrgSettingsWrapper.vue',
      'src/views/repo/settings/RepoSettings.vue',
      'src/views/repo/RepoBranches.vue',
      'src/views/repo/RepoBranch.vue',
      'src/views/repo/RepoPullRequests.vue',
      'src/views/repo/RepoPullRequest.vue',
      'src/views/repo/RepoManualPipeline.vue',
      'src/views/repo/RepoPipelines.vue',
      'src/views/repo/pipeline/Pipeline.vue',
      'src/views/repo/pipeline/PipelineWrapper.vue',
    ],
    rules: { '@intlify/vue-i18n/no-raw-text': 'off' },
  },
  {
    ignores: [
      'dist', 'coverage/', 'package.json', 'tsconfig.eslint.json', 'tsconfig.json',
      'src/assets/locales/**/*', '!src/assets/locales/en.json', 'components.d.ts',
    ],
  },
);
