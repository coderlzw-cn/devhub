<script setup lang="ts">
import { cn } from '@/utils/cn'
import { computed, onMounted, ref } from 'vue'
import { storeToRefs } from 'pinia'
import { RouterView, useRoute, useRouter } from 'vue-router'
import { Check } from 'lucide-vue-next'
import AppSidebar from '@/components/AppSidebar.vue'
import AppTopbar from '@/components/AppTopbar.vue'
import ProfilePage from '@/pages/ProfilePage.vue'
import type { Resource } from '@/data/resources'
import type { RouteLocationRaw } from 'vue-router'
import { useResourceCategoriesQuery } from '@/composables/useResourceCategoriesQuery'
import { usePreferencesStore, type DefaultPage, type ThemeMode } from '@/stores/preferences'

const route = useRoute()
const router = useRouter()
const preferences = usePreferencesStore()
const categoriesQuery = useResourceCategoriesQuery()
const { themeMode, defaultPage, isDark } = storeToRefs(preferences)
const activeCategory = computed(() => (typeof route.query.category === 'string' ? route.query.category : '全部'))
const activePage = computed(() => {
  if (route.name === 'categories') return '分类管理'
  if (route.name === 'resources') return '资源管理'
  return route.name === 'profile' ? '个人中心' : '首页'
})
const isProfileRoute = computed(() => route.name === 'profile')
const breadcrumbItems = computed<{ label: string; to?: RouteLocationRaw }[]>(() => {
  const home = { label: '首页', to: { name: 'home' } }
  if (route.name === 'home') return [{ label: '首页' }]
  if (route.name === 'categories') return [home, { label: '分类管理' }]
  if (route.name === 'profile') return [home, { label: '个人中心' }]

  const resources = { label: '资源管理', to: { name: 'resources' } }
  if (activeCategory.value === '全部') return [home, { label: resources.label }]

  const parent = categoriesQuery.data.value?.find(
    (category) => category.name === activeCategory.value || category.children.some((child) => child.name === activeCategory.value),
  )
  if (!parent) return [home, resources, { label: activeCategory.value }]
  if (parent.name === activeCategory.value) return [home, resources, { label: parent.name }]

  return [
    home,
    resources,
    { label: parent.name, to: { name: 'resources', query: { category: parent.name } } },
    { label: activeCategory.value },
  ]
})
const toast = ref('')

function navigate(page: string) {
  if (page === '设置' || page === '个人中心') {
    void router.push({ name: 'profile' })
    return
  }
  if (page === '首页') void router.push({ name: 'home' })
  else if (page === '分类管理') void router.push({ name: 'categories' })
  else if (page === '资源管理') void router.push({ name: 'resources' })
  else void router.push({ name: 'resources', query: { category: page } })
}

function updatePreference(key: 'themeMode' | 'defaultPage', value: string) {
  if (key === 'themeMode' && (value === 'system' || value === 'light' || value === 'dark'))
    preferences.setThemeMode(value as ThemeMode)
  if (key === 'defaultPage' && ['首页', '资源管理', '分类管理'].includes(value))
    preferences.setDefaultPage(value as DefaultPage)
}

function openResource(resource: Resource) {
  window.open(resource.url, '_blank', 'noopener,noreferrer')
}

function showToast(message: string) {
  toast.value = message
  window.setTimeout(() => {
    toast.value = ''
  }, 1800)
}

onMounted(() => {
  if (route.name === 'home' && defaultPage.value !== '首页') navigate(defaultPage.value)
})
</script>

<template>
  <div
    class="min-h-screen bg-[#0b0f19] text-slate-100"
    :class="cn(!isDark && 'light bg-[#f5f7fb] text-slate-800')"
  >
    <Transition name="route-page" mode="out-in">
      <ProfilePage
        v-if="isProfileRoute"
        key="profile"
        :theme-mode="themeMode"
        :default-page="defaultPage"
        :is-dark="isDark"
        @update-preference="updatePreference"
        @navigate="navigate"
        @notify="showToast"
      />
      <div v-else key="application" class="flex min-h-screen">
      <AppSidebar :active-page="activePage" :on-navigate="navigate" />
      <main
        class="ml-[66px] min-h-screen w-[calc(100%-66px)] min-w-0 sm:ml-[214px] sm:w-[calc(100%-214px)] lg:ml-[244px] lg:w-[calc(100%-244px)]"
      >
        <AppTopbar
          :theme-mode="themeMode"
          :breadcrumb-items="breadcrumbItems"
          @update:theme-mode="updatePreference('themeMode', $event)"
        />
        <RouterView v-slot="{ Component }">
          <Transition name="route-page" mode="out-in">
            <component
              :is="Component"
              v-if="route.name === 'resources'"
              :key="`resources-${route.name}`"
              :active-category="activeCategory"
              @category="navigate"
              @open="openResource"
            /><component :is="Component" v-else :key="`page-${String(route.name)}`" @navigate="navigate" />
          </Transition>
        </RouterView>
      </main>
      </div>
    </Transition>
    <Transition
      enter-active-class="transition duration-200"
      enter-from-class="translate-y-2 opacity-0"
      enter-to-class="translate-y-0 opacity-100"
      leave-active-class="transition duration-200"
      leave-from-class="translate-y-0 opacity-100"
      leave-to-class="translate-y-2 opacity-0"
    >
      <div
        v-if="toast"
        class="fixed bottom-[88px] left-1/2 z-30 flex -translate-x-1/2 items-center gap-2 rounded-lg border border-emerald-800 bg-emerald-950 px-4 py-[11px] text-[11px] whitespace-nowrap text-emerald-200 shadow-xl"
      >
        <Check :size="16" />{{ toast }}
      </div>
    </Transition>
  </div>
</template>
