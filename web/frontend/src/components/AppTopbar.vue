<script setup lang="ts">
import { cn } from '@/utils/cn'
import { computed } from 'vue'
import { Command, Monitor, Moon, Search, Sun } from 'lucide-vue-next'
import Dropdown from '@/components/ui/Dropdown.vue'
import Input from '@/components/ui/Input.vue'
import AppBreadcrumbs from '@/components/AppBreadcrumbs.vue'
import type { RouteLocationRaw } from 'vue-router'

type ThemeMode = 'system' | 'light' | 'dark'

const props = defineProps<{
  themeMode: ThemeMode
  breadcrumbItems: { label: string; to?: RouteLocationRaw }[]
}>()
const emit = defineEmits<{
  (event: 'update:themeMode', value: ThemeMode): void
}>()

const themeItems = [
  { label: '跟随系统', value: 'system', icon: Monitor },
  { label: '浅色模式', value: 'light', icon: Sun },
  { label: '深色模式', value: 'dark', icon: Moon },
]
const themeIcon = computed(() => (props.themeMode === 'system' ? Monitor : props.themeMode === 'light' ? Sun : Moon))

function selectTheme(value: string) {
  if (value === 'system' || value === 'light' || value === 'dark') emit('update:themeMode', value)
}
</script>

<template>
  <header
    class="sticky top-0 z-20 flex h-15.5 items-center gap-2.5 border-b border-slate-800 bg-[#0d1420]/90 px-[13px] text-slate-100 backdrop-blur-2xl sm:h-[74px] sm:gap-4 sm:px-6 lg:px-[42px] light:border-slate-200! light:bg-white/90! light:text-slate-800!"
  >
    <AppBreadcrumbs class="max-w-[43%] shrink" :items="breadcrumbItems" />
    <div class="min-w-0 flex-1"></div>

    <div class="ml-auto flex w-[56%] max-w-115 min-w-0 shrink-0 items-center gap-1.5 sm:gap-2">
      <Input
        appearance="plain"
        class="min-w-0 flex-1"
        readonly
        placeholder="搜索工具、网站、分类..."
      >
        <template #leading><Search :size="17" /></template>
        <template #trailing>
          <kbd class="hidden items-center gap-0.75 rounded-sm px-1.5 py-0.75 text-[10px] sm:flex" aria-hidden="true">
            <Command :size="11" /> K
          </kbd>
        </template>
      </Input>

      <Dropdown
        :items="themeItems"
        align="end"
        placement="bottom"
        fit
        trigger-label="选择主题模式"
        class="shrink-0"
        @select="selectTheme"
      >
        <template #default="{ open }">
          <span
            class="grid size-9 place-items-center rounded-lg text-slate-400 transition hover:bg-slate-800 hover:text-white light:text-slate-600 light:hover:bg-slate-100 light:hover:text-slate-900"
            :class="cn(open && 'bg-slate-800 text-white light:bg-slate-100 light:text-slate-900')"
            :title="
              themeMode === 'system' ? '外观：跟随系统' : themeMode === 'light' ? '外观：浅色模式' : '外观：深色模式'
            "
            aria-label="选择主题模式"
          >
            <component :is="themeIcon" :size="18" />
          </span>
        </template>
      </Dropdown>
    </div>
  </header>
</template>
