<script setup lang="ts">
import { computed } from 'vue'
import { FlexRender, tableFeatures, useTable, type ColumnDef } from '@tanstack/vue-table'
import { ArrowUpRight, ChevronDown, ChevronRight, Pencil, Plus, Trash2 } from 'lucide-vue-next'
import { resolveIcon } from '@/utils/icon-map'

export interface CategoryRow {
  id: string
  name: string
  level: '一级分类' | '二级分类'
  parentId: string | null
  parentName: string
  childCount: number
  icon?: string | null
}

const props = defineProps<{
  rows: CategoryRow[]
  collapsedIds: Set<string>
  searchQuery: string
}>()
const emit = defineEmits<{
  toggle: [id: string]
  addChild: [id: string]
  edit: [row: CategoryRow]
  delete: [row: CategoryRow]
  navigate: [page: string]
}>()
const features = tableFeatures({})
const columns: ColumnDef<typeof features, CategoryRow>[] = [
  { accessorKey: 'name', header: '分类名称' },
  { accessorKey: 'parentName', header: '上级分类' },
  { accessorKey: 'childCount', header: '子分类数' },
  { id: 'actions', accessorKey: 'name', header: '操作' },
]
const tableData = computed(() => props.rows)
const table = useTable({ features, columns, data: tableData })
</script>

<template>
  <div class="overflow-x-auto rounded-lg border border-slate-800 bg-slate-900/50 light:border-slate-200 light:bg-white">
    <table class="w-full min-w-[700px] border-collapse text-left text-xs">
      <thead class="border-b border-slate-800 bg-slate-900/80 light:border-slate-200 light:bg-slate-50">
        <tr v-for="headerGroup in table.getHeaderGroups()" :key="headerGroup.id">
          <th
            v-for="header in headerGroup.headers"
            :key="header.id"
            class="h-10 px-4 font-medium whitespace-nowrap text-slate-400 light:text-slate-500"
          >
            <FlexRender v-if="!header.isPlaceholder" :header="header" />
          </th>
        </tr>
      </thead>
      <tbody class="divide-y divide-slate-800 light:divide-slate-100">
        <tr
          v-for="row in table.getRowModel().rows"
          :key="row.id"
          class="h-[54px] transition hover:bg-slate-800/60 light:hover:bg-slate-50"
        >
          <td v-for="cell in row.getAllCells()" :key="cell.id" class="px-4 py-2.5 text-slate-300 light:text-slate-700">
            <div v-if="cell.column.id === 'actions'" class="flex items-center justify-end gap-1">
              <button
                v-if="row.original.level === '一级分类'"
                class="grid size-8 place-items-center rounded-md text-slate-400 hover:bg-slate-700 hover:text-violet-300 light:hover:bg-slate-100"
                :aria-label="`为 ${row.original.name} 添加二级分类`"
                title="添加二级分类"
                @click="emit('addChild', row.original.id)"
              >
                <Plus :size="15" />
              </button>
              <button
                class="grid size-8 place-items-center rounded-md text-slate-400 hover:bg-slate-700 hover:text-violet-300 light:hover:bg-slate-100"
                :aria-label="`编辑 ${row.original.name}`"
                title="编辑分类"
                @click="emit('edit', row.original)"
              >
                <Pencil :size="14" />
              </button>
              <button
                class="grid size-8 place-items-center rounded-md text-slate-400 hover:bg-rose-950 hover:text-rose-300 light:hover:bg-rose-50 light:hover:text-rose-700"
                :aria-label="`删除 ${row.original.name}`"
                title="删除分类"
                @click="emit('delete', row.original)"
              >
                <Trash2 :size="14" />
              </button>
              <button
                class="ml-1 inline-flex items-center gap-1 px-2 whitespace-nowrap text-violet-300 hover:text-violet-200 light:text-violet-700 light:hover:text-violet-900"
                :aria-label="`查看 ${row.original.name} 下的资源`"
                title="查看资源"
                @click="emit('navigate', row.original.name)"
              >
                <ArrowUpRight :size="14" />
              </button>
            </div>
            <span v-else-if="cell.column.id === 'name'" class="flex items-center gap-2.5">
              <button
                v-if="row.original.level === '一级分类' && row.original.childCount > 0"
                class="grid size-6 shrink-0 place-items-center rounded-sm text-slate-400 hover:bg-slate-700 hover:text-white light:hover:bg-slate-100 light:hover:text-slate-800"
                :aria-label="collapsedIds.has(row.original.id) ? '展开子分类' : '折叠子分类'"
                :aria-expanded="!collapsedIds.has(row.original.id)"
                :title="collapsedIds.has(row.original.id) ? '展开子分类' : '折叠子分类'"
                @click="emit('toggle', row.original.id)"
              >
                <ChevronRight v-if="collapsedIds.has(row.original.id)" :size="15" />
                <ChevronDown v-else :size="15" />
              </button>
              <span v-else class="w-6 shrink-0"></span>
              <component
                :is="resolveIcon(row.original.icon)"
                :size="15"
                :class="row.original.level === '一级分类' ? 'text-violet-400' : 'text-slate-500'"
              />
              <span
                :class="row.original.level === '一级分类' ? 'font-semibold text-slate-100 light:text-slate-800' : ''"
              >
                {{ row.original.name }}
              </span>
            </span>
            <span v-else-if="cell.column.id === 'parentName'">{{ row.original.parentName }}</span>
            <span v-else-if="cell.column.id === 'childCount'">
              {{ row.original.level === '一级分类' ? row.original.childCount : '—' }}
            </span>
            <FlexRender v-else :cell="cell" />
          </td>
        </tr>
        <tr v-if="table.getRowModel().rows.length === 0">
          <td colspan="4" class="h-32 text-center text-sm text-slate-400">
            {{ searchQuery ? '没有匹配的分类' : '暂无分类' }}
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
