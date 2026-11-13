<script setup lang="ts">
import { ArrowUpRight, Layers3, Star, ChevronRight } from 'lucide-vue-next'
import { computed } from 'vue'
import { useResourcesQuery } from '@/composables/useResourcesQuery'
import { useResourceCategoriesQuery } from '@/composables/useResourceCategoriesQuery'
import Skeleton from '@/components/ui/Skeleton.vue'
import heroImage from '@/assets/devhub-hero.png'
import { resolveIcon } from '@/utils/icon-map'
const resourcesQuery = useResourcesQuery()
const categoriesQuery = useResourceCategoriesQuery()
const resources = computed(() => resourcesQuery.data.value ?? [])
const isPending = computed(() => resourcesQuery.isPending.value || categoriesQuery.isPending.value)
const isError = computed(() => resourcesQuery.isError.value || categoriesQuery.isError.value)
const isFetching = computed(() => resourcesQuery.isFetching.value || categoriesQuery.isFetching.value)
const emit = defineEmits<{
  (event: 'navigate', page: string): void
}>()
const categories = computed(() =>
  (categoriesQuery.data.value ?? []).map((category) => ({
    label: category.name,
    count: resources.value.filter((resource) =>
      [category.id, ...category.children.map((child) => child.id)].includes(resource.categoryId ?? ''),
    ).length,
    tone: 'blue',
    icon: category.icon,
  })),
)
function retryQueries() {
  void Promise.all([resourcesQuery.refetch(), categoriesQuery.refetch()])
}
const categoryToneClass: Record<string, string> = {
  purple: 'text-violet-300',
  cyan: 'text-cyan-300',
  orange: 'text-amber-300',
  blue: 'text-blue-300',
  violet: 'text-sky-300',
  sky: 'text-violet-300',
  indigo: 'text-cyan-300',
  green: 'text-indigo-300',
}
</script>
<template>
  <div
    v-if="isPending"
    class="mx-auto grid min-h-[50vh] max-w-[1600px] content-start gap-5 px-3.5 pt-5 sm:px-6 sm:pt-7 lg:px-[42px]"
    aria-label="首页加载中"
  >
    <Skeleton class="h-60 w-full rounded-xl" />
    <div class="grid grid-cols-2 gap-3 sm:grid-cols-4">
      <Skeleton v-for="item in 8" :key="item" class="h-20 rounded-lg" />
    </div>
    <div class="grid grid-cols-2 gap-3 sm:grid-cols-3">
      <Skeleton v-for="item in 3" :key="item" class="h-36 rounded-xl" />
    </div>
  </div>
  <div v-else-if="isError" class="grid min-h-[50vh] place-content-center gap-3 text-center">
    <p class="text-sm text-rose-300">首页数据加载失败，请检查连接后重试。</p>
    <button
      class="mx-auto rounded-md border border-slate-700 px-3 py-1.5 text-xs text-slate-200 disabled:opacity-50"
      :disabled="isFetching"
      @click="retryQueries"
    >
      {{ isFetching ? '正在重试…' : '重试' }}
    </button>
  </div>
  <div v-else>
    <div class="mx-auto max-w-[1600px] px-[13px] pt-[17px] pb-10 sm:px-6 sm:pt-[29px] lg:px-[42px]">
      <section class="relative h-[250px] overflow-hidden rounded-[15px] border border-blue-900 sm:h-[236px]">
        <img
          :src="heroImage"
          alt=""
          class="absolute inset-0 size-full object-cover object-[62%_center] sm:object-[center_52%]"
        />
        <div class="relative z-10 p-4 sm:px-[30px] sm:py-[31px]">
          <h1 class="mb-2 text-[25px] font-extrabold tracking-[-1px] text-slate-50 sm:text-[34px]">开发者资源中心</h1>
          <p class="text-xs text-blue-100 sm:text-base">整理和管理你的开发资源</p>
          <div class="mt-6 flex flex-wrap gap-3 sm:mt-[31px] sm:gap-[34px]">
            <div
              class="grid grid-cols-[27px_auto] grid-rows-[22px_17px] items-center gap-x-2 sm:grid-cols-[32px_auto] sm:gap-x-[9px]"
            >
              <Layers3
                class="row-span-2 size-[26px] rounded-lg bg-sky-900/40 p-1.5 text-sky-400 sm:size-[31px]"
              /><strong class="text-base leading-none text-white sm:text-[21px]">{{ resources.length }}</strong
              ><span class="text-[11px] text-blue-100/70">全部资源</span>
            </div>
          </div>
        </div>
      </section>
      <section class="mt-6">
        <div class="mb-[15px] flex items-center justify-between">
          <div class="flex items-center gap-2.5">
            <span class="h-[17px] w-[3px] rounded-sm bg-violet-400"></span>
            <h2 class="text-[15px] font-bold">资源分类</h2>
          </div>
          <button
            class="flex items-center gap-1 text-[10px] text-violet-300 hover:text-violet-200"
            @click="emit('navigate', '资源管理')"
          >
            查看全部 <ArrowUpRight :size="15" />
          </button>
        </div>
        <div class="grid grid-cols-2 gap-2 sm:grid-cols-4">
          <button
            v-for="item in categories"
            :key="item.label"
            class="flex h-[69px] min-w-0 items-center gap-1.5 rounded-[10px] border border-slate-800 bg-linear-to-br from-slate-800 to-slate-900 px-2 text-left transition hover:-translate-y-0.5 hover:border-violet-500 sm:h-[87px] sm:gap-3 sm:px-3.5 light:border-slate-200! light:bg-white! light:bg-none! light:hover:bg-violet-50!"
            @click="emit('navigate', item.label)"
          >
            <component
              :is="resolveIcon(item.icon)"
              class="size-8 shrink-0 sm:size-12"
              :class="categoryToneClass[item.tone]"
              :aria-label="`${item.label} 图标`"
              aria-hidden="true"
            />
            <div class="flex min-w-0 flex-1 flex-col gap-1.5">
              <strong class="truncate text-[10px] text-slate-100 sm:text-sm light:text-slate-800!">{{
                item.label
              }}</strong
              ><span class="text-[9px] text-slate-400 sm:text-[11px] light:text-slate-500!"
                >{{ item.count }} 个资源</span
              >
            </div>
            <ChevronRight class="hidden text-slate-500 sm:block light:text-slate-400!" :size="16" />
          </button>
        </div>
      </section>
      <section class="mt-[19px]">
        <div
          class="min-w-0 rounded-[14px] border border-slate-800 bg-[#121b2a] p-3.5 sm:p-4 light:border-slate-200 light:bg-white"
        >
          <div class="mb-3 flex items-center justify-between">
            <div class="flex items-center gap-2.5">
              <span class="h-4.25 w-0.75 rounded-sm bg-amber-400"></span>
              <h2 class="text-[15px] font-bold text-slate-100 light:text-slate-800">置顶站点</h2>
            </div>
            <button
              class="flex items-center gap-1 text-[10px] text-violet-300 light:text-violet-700"
              @click="emit('navigate', '资源管理')"
            >
              全部站点 <ArrowUpRight :size="14" />
            </button>
          </div>
          <p v-if="!resources.some((resource) => resource.isPinned)" class="text-sm text-slate-400">暂无置顶站点</p>
          <div class="grid grid-cols-2 gap-2 sm:grid-cols-3">
            <button
              v-for="recommendation in resources.filter((resource) => resource.isPinned)"
              :key="recommendation.name"
              class="relative flex min-h-34.75 min-w-0 flex-col items-start gap-[5px] rounded-xl border border-slate-800 bg-linear-to-br from-slate-800 to-slate-900 p-[13px] text-left transition hover:-translate-y-px hover:border-violet-500 light:border-slate-200! light:bg-white! light:bg-none! light:hover:bg-violet-50!"
              @click="emit('navigate', '资源管理')"
            >
              <div
                class="grid size-9.75 place-items-center rounded-[9px] bg-slate-950 text-[17px] font-extrabold text-white light:bg-slate-100 light:text-slate-800"
              >
                {{ recommendation.name.slice(0, 1) }}
              </div>
              <strong class="text-[13px] text-slate-100 light:text-slate-800">{{ recommendation.name }}</strong
              ><span class="line-clamp-2 min-h-6.75 text-[10px] leading-[1.35] text-slate-400 light:text-slate-500">{{
                recommendation.description
              }}</span
              ><small
                class="rounded-sm bg-blue-950 px-1.5 py-[3px] text-[9px] text-blue-200 light:bg-blue-50 light:text-blue-700"
                >置顶</small
              ><Star class="absolute top-3 right-3 text-blue-200 light:text-amber-500" :size="15" />
            </button>
          </div>
        </div>
      </section>
    </div>
  </div>
</template>
