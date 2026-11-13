<script setup lang="ts">
import { cn } from '@/utils/cn'
import { ChevronDown, Code2, Home, FolderTree, Layers3, Settings, UserRound } from 'lucide-vue-next'
import Dropdown from '@/components/ui/Dropdown.vue'
import { computed } from 'vue'
import { useResourceCategoriesQuery } from '@/composables/useResourceCategoriesQuery'
import { resolveIcon } from '@/utils/icon-map'

defineProps<{ activePage: string; onNavigate: (page: string) => void }>()
const categoriesQuery = useResourceCategoriesQuery()
const items = computed(() => categoriesQuery.data.value ?? [])
const accountMenuItems = [
  { label: '个人中心', value: '个人中心', icon: UserRound },
  { label: '设置', value: '设置', icon: Settings },
]
</script>

<template>
  <aside
    class="fixed inset-y-0 left-0 z-10 flex min-h-screen w-[66px] flex-col border-r border-slate-800 bg-[#101725] px-2 py-[18px] text-slate-100 sm:w-[214px] sm:px-3.5 sm:py-[22px] lg:w-[244px] light:border-slate-200! light:bg-white! light:text-slate-800!"
  >
    <div class="flex items-center gap-[11px] px-1.5 pb-[26px] sm:pb-[31px]">
      <div
        class="grid size-9 shrink-0 place-items-center rounded-[11px] bg-linear-to-br from-violet-400 to-indigo-700 text-white shadow-[0_5px_20px_#6755f540]"
      >
        <Code2 :size="21" :stroke-width="2.5" />
      </div>
      <div class="hidden sm:block">
        <div class="text-[17px] font-bold tracking-[-.4px]">DevHub</div>
        <div class="mt-0.5 text-[9px] tracking-[.25px] text-slate-400 light:text-slate-500">Developer Resource Hub</div>
      </div>
    </div>
    <nav class="flex flex-col gap-1">
      <button
        class="flex h-[39px] w-full items-center justify-center gap-3 rounded-lg px-0 text-[13px] text-slate-400 transition hover:bg-slate-800 hover:text-indigo-50 sm:justify-start sm:px-3 light:text-slate-600 light:hover:bg-slate-100 light:hover:text-slate-900"
        :class="
          activePage === '首页'
            ? 'bg-linear-to-r from-indigo-700 to-indigo-950 text-white! shadow-[inset_0_0_0_1px_#8d7eff30]'
            : ''
        "
        title="首页"
        @click="onNavigate('首页')"
      >
        <Home :size="17" /><span class="hidden flex-1 text-left sm:inline">首页</span>
      </button>
      <button
        class="flex h-[39px] w-full items-center justify-center gap-3 rounded-lg px-0 text-[13px] text-slate-400 transition hover:bg-slate-800 hover:text-indigo-50 sm:justify-start sm:px-3 light:text-slate-600 light:hover:bg-slate-100 light:hover:text-slate-900"
        :class="activePage === '分类管理' ? 'bg-linear-to-r from-indigo-700 to-indigo-950 text-white!' : ''"
        title="分类管理"
        @click="onNavigate('分类管理')"
      >
        <FolderTree :size="17" /><span class="hidden flex-1 text-left sm:inline">分类管理</span>
      </button>
      <button
        class="flex h-[39px] w-full items-center justify-center gap-3 rounded-lg px-0 text-[13px] text-slate-400 transition hover:bg-slate-800 hover:text-indigo-50 sm:justify-start sm:px-3 light:text-slate-600 light:hover:bg-slate-100 light:hover:text-slate-900"
        :class="
          activePage === '资源管理'
            ? 'bg-linear-to-r from-indigo-700 to-indigo-950 text-white! shadow-[inset_0_0_0_1px_#8d7eff30]'
            : ''
        "
        title="资源管理"
        @click="onNavigate('资源管理')"
      >
        <Layers3 :size="17" /><span class="hidden flex-1 text-left sm:inline">资源管理</span>
      </button>
      <div
        class="hidden items-center justify-between px-3 pt-[21px] pb-2 text-[11px] tracking-[.4px] text-slate-500 sm:flex light:text-slate-400"
      >
        <span>资源分类</span><ChevronDown :size="14" />
      </div>
      <button
        v-for="item in items"
        :key="item.id"
        class="flex h-[39px] w-full items-center justify-center gap-3 rounded-lg px-0 text-[13px] text-slate-400 transition hover:bg-slate-800 hover:text-indigo-50 sm:justify-start sm:px-3 sm:pl-[18px] light:text-slate-600 light:hover:bg-slate-100 light:hover:text-slate-900"
        :class="
          activePage === item.name
            ? 'bg-linear-to-r from-indigo-700 to-indigo-950 text-white! shadow-[inset_0_0_0_1px_#8d7eff30]'
            : ''
        "
        :title="item.name"
        @click="onNavigate(item.name)"
      >
        <component :is="resolveIcon(item.icon)" class="size-[17px] shrink-0" aria-hidden="true" />
        <span class="hidden flex-1 truncate text-left font-medium sm:inline">
          {{ item.name }}
        </span>
      </button>
    </nav>
    <div class="mt-auto hidden sm:block">
      <Dropdown :items="accountMenuItems" align="start" @select="onNavigate">
        <template #default="{ open }">
          <div
            class="flex w-full items-center gap-2.5 border-t border-slate-700 px-[7px] pt-[15px] text-left text-slate-400 transition hover:text-slate-200 light:border-slate-200! light:text-slate-500! light:hover:text-slate-800!"
            :class="cn(open && 'text-white light:text-slate-900!')"
          >
            <div
              class="grid size-[31px] shrink-0 place-items-center rounded-full bg-linear-to-br from-orange-200 to-purple-500 font-bold text-white"
            >
              L
            </div>
            <div class="min-w-0 flex-1">
              <strong class="block truncate text-[11px] text-slate-200 light:text-slate-700!">梁志伟</strong
              ><span class="mt-[3px] block truncate text-[10px] text-slate-500 light:text-slate-500!"
                >个人工作空间</span
              >
            </div>
            <ChevronDown class="shrink-0 transition" :class="cn(open && 'rotate-180')" :size="15" />
          </div>
        </template>
      </Dropdown>
    </div>
  </aside>
</template>
