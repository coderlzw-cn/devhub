<script setup lang="ts">
import { computed } from 'vue'
import { ChevronLeft, ChevronRight } from 'lucide-vue-next'
import Select from '@/components/ui/Select.vue'

defineOptions({ name: 'UiPagination' })

type PageItem = number | 'ellipsis-start' | 'ellipsis-end'

const props = withDefaults(
  defineProps<{
    currentPage: number
    pageSize: number
    total: number
    pageSizeOptions?: number[]
  }>(),
  { pageSizeOptions: () => [9, 18, 36] },
)

const emit = defineEmits<{
  (event: 'update:currentPage', page: number): void
  (event: 'update:pageSize', pageSize: number): void
}>()

const pageCount = computed(() => Math.max(1, Math.ceil(props.total / props.pageSize)))
const rangeStart = computed(() => (props.total === 0 ? 0 : (props.currentPage - 1) * props.pageSize + 1))
const rangeEnd = computed(() => Math.min(props.currentPage * props.pageSize, props.total))
const pageItems = computed<PageItem[]>(() => {
  const count = pageCount.value
  if (count <= 5) return Array.from({ length: count }, (_, index) => index + 1)

  const pages = [...new Set([1, count, props.currentPage - 1, props.currentPage, props.currentPage + 1])]
    .filter((page) => page >= 1 && page <= count)
    .sort((left, right) => left - right)
  const items: PageItem[] = []
  for (let index = 0; index < pages.length; index += 1) {
    const page = pages[index]
    const previousPage = pages[index - 1]
    if (previousPage !== undefined) {
      const gap = page - previousPage
      if (gap === 2) items.push(previousPage + 1)
      else if (gap > 2) items.push(index === 1 ? 'ellipsis-start' : 'ellipsis-end')
    }
    items.push(page)
  }
  return items
})

function changePage(page: number) {
  if (page < 1 || page > pageCount.value || page === props.currentPage) return
  emit('update:currentPage', page)
}

function changePageSize(value: string | number) {
  const nextSize = Number(value)
  if (!props.pageSizeOptions.includes(nextSize) || nextSize === props.pageSize) return
  emit('update:pageSize', nextSize)
  emit('update:currentPage', 1)
}
</script>

<template>
  <nav class="flex flex-wrap items-center justify-between gap-3 px-1 pt-5 text-[11px] text-slate-400 light:text-slate-500" aria-label="分页">
    <span class="w-full sm:w-auto">显示 {{ rangeStart }}-{{ rangeEnd }} / {{ total }} 条</span>
    <div class="order-3 flex items-center gap-1 sm:order-0">
      <button
        class="grid size-8 place-items-center rounded-md border border-slate-700 bg-slate-900 text-slate-300 transition hover:border-violet-400 hover:text-white disabled:cursor-not-allowed disabled:opacity-35 light:border-slate-200 light:bg-white light:text-slate-600 light:hover:text-violet-700"
        :disabled="currentPage <= 1"
        aria-label="上一页"
        title="上一页"
        @click="changePage(currentPage - 1)"
      >
        <ChevronLeft :size="16" />
      </button>
      <template v-for="item in pageItems" :key="item">
        <span v-if="typeof item !== 'number'" class="grid size-8 place-items-center text-slate-500" aria-hidden="true">…</span>
        <button
          v-else
          class="grid size-8 place-items-center rounded-md border border-slate-700 bg-slate-900 text-xs text-slate-300 transition hover:border-violet-400 light:border-slate-200 light:bg-white light:text-slate-600"
          :class="item === currentPage ? 'border-violet-400 bg-violet-600 text-white! hover:text-white' : ''"
          :aria-label="`第 ${item} 页`"
          :aria-current="item === currentPage ? 'page' : undefined"
          @click="changePage(item)"
        >
          {{ item }}
        </button>
      </template>
      <button
        class="grid size-8 place-items-center rounded-md border border-slate-700 bg-slate-900 text-slate-300 transition hover:border-violet-400 hover:text-white disabled:cursor-not-allowed disabled:opacity-35 light:border-slate-200 light:bg-white light:text-slate-600 light:hover:text-violet-700"
        :disabled="currentPage >= pageCount"
        aria-label="下一页"
        title="下一页"
        @click="changePage(currentPage + 1)"
      >
        <ChevronRight :size="16" />
      </button>
    </div>
    <label class="ml-auto flex items-center gap-2 whitespace-nowrap sm:ml-0">
      每页
      <Select
        :model-value="pageSize"
        :options="pageSizeOptions"
        size="sm"
        class="min-w-[68px]"
        :aria-label="'每页条数'"
        @update:model-value="changePageSize"
      />
      条
    </label>
  </nav>
</template>
