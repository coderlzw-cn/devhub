import { computed, ref } from 'vue'
import { defineStore } from 'pinia'

export type ThemeMode = 'system' | 'light' | 'dark'
export type DefaultPage = '首页' | '资源管理' | '分类管理'

interface PreferencesState {
  themeMode: ThemeMode
  defaultPage: DefaultPage
}

const defaults: PreferencesState = {
  themeMode: 'system',
  defaultPage: '首页',
}

export const usePreferencesStore = defineStore(
  'preferences',
  () => {
    const themeMode = ref<ThemeMode>(defaults.themeMode)
    const defaultPage = ref<DefaultPage>(defaults.defaultPage)
    const systemPrefersDark = ref(
      typeof window !== 'undefined' && window.matchMedia('(prefers-color-scheme: dark)').matches,
    )
    const isDark = computed(() => (themeMode.value === 'system' ? systemPrefersDark.value : themeMode.value === 'dark'))

    if (typeof window !== 'undefined') {
      const colorScheme = window.matchMedia('(prefers-color-scheme: dark)')
      colorScheme.addEventListener('change', (event) => {
        systemPrefersDark.value = event.matches
      })
    }

    function setThemeMode(value: ThemeMode) {
      themeMode.value = value
    }
    function setDefaultPage(value: DefaultPage) {
      defaultPage.value = value
    }

    return { themeMode, defaultPage, isDark, setThemeMode, setDefaultPage }
  },
  {
    persist: {
      key: 'devhub-preferences',
      pick: ['themeMode', 'defaultPage'],
      afterHydrate: ({ store }) => {
        if (store.themeMode !== 'system' && store.themeMode !== 'light' && store.themeMode !== 'dark')
          store.themeMode = defaults.themeMode
        if (!['首页', '资源管理', '分类管理'].includes(store.defaultPage))
          store.defaultPage = defaults.defaultPage
      },
    },
    share: {
      enable: true,
      initialize: true,
    },
  },
)
