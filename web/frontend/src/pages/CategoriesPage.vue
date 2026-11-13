<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useQueryClient } from '@tanstack/vue-query'
import { FolderTree, Plus, Search } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import Input from '@/components/ui/Input.vue'
import Skeleton from '@/components/ui/Skeleton.vue'
import CategoryEditorDialog from '@/pages/categories/components/CategoryEditorDialog.vue'
import CategoryDeleteDialog from '@/pages/categories/components/CategoryDeleteDialog.vue'
import CategoryTable, { type CategoryRow } from '@/pages/categories/components/CategoryTable.vue'
import { useResourceCategoriesQuery } from '@/composables/useResourceCategoriesQuery'
import {
  usePostApiCategories,
  usePatchApiCategoriesId,
  useDeleteApiCategoriesId,
  getGetApiCategoriesQueryKey,
} from '@/api/generated/categories'
import type { BodyPageCategory } from '@/api/generated/models'

const categoryQueryKey = getGetApiCategoriesQueryKey()
const createMutation = usePostApiCategories()
const updateMutation = usePatchApiCategoriesId()
const deleteMutation = useDeleteApiCategoriesId()
const mutationError = ref('')
const isSaving = computed(() => createMutation.isPending.value || updateMutation.isPending.value)
const isDeleting = computed(() => deleteMutation.isPending.value)
const initialSlug = ref('')
const initialUserId = ref('')
const emit = defineEmits<{ (event: 'navigate', page: string): void }>()
const queryClient = useQueryClient()
const categoriesQuery = useResourceCategoriesQuery()
const searchQuery = ref('')
const collapsedIds = ref(new Set<string>())
const dialogMode = ref<'create' | 'edit' | null>(null)
const editingRow = ref<CategoryRow | null>(null)
const targetParentId = ref<string | null>(null)
const dialogInitialName = ref('')
const pendingDelete = ref<CategoryRow | null>(null)
const categories = computed(() => categoriesQuery.data.value ?? [])
watch(
  categories,
  (items) => {
    collapsedIds.value = new Set(
      items.filter((category) => category.children.length > 0).map((category) => category.id),
    )
  },
  { immediate: true },
)
const categoryCount = computed(() =>
  categories.value.reduce((count, category) => count + 1 + category.children.length, 0),
)

const rows = computed<CategoryRow[]>(() => {
  const query = searchQuery.value.trim().toLocaleLowerCase()
  return categories.value.flatMap((category) => {
    const matchesParent = !query || category.name.toLocaleLowerCase().includes(query)
    const matchingChildren = category.children.filter(
      (child) =>
        !query || child.name.toLocaleLowerCase().includes(query) || category.name.toLocaleLowerCase().includes(query),
    )
    const showChildren = query.length > 0 || !collapsedIds.value.has(category.id)
    const parentRow: CategoryRow = {
      id: category.id,
      name: category.name,
      level: '一级分类',
      parentId: null,
      parentName: '—',
      childCount: category.children.length,
      icon: category.icon,
    }
    const childRows: CategoryRow[] = category.children
      .filter(() => showChildren && (matchesParent || matchingChildren.length > 0))
      .filter((child) => !query || matchingChildren.some((match) => match.id === child.id))
      .map((child) => ({
        id: child.id,
        name: child.name,
        level: '二级分类',
        parentId: category.id,
        parentName: category.name,
        childCount: 0,
        icon: category.icon,
      }))
    return !query || matchesParent || childRows.length > 0 ? [parentRow, ...childRows] : []
  })
})
function openDelete(row: CategoryRow) {
  mutationError.value = ''
  pendingDelete.value = row
}

function retry() {
  void categoriesQuery.refetch()
}

function toggleParent(id: string) {
  const next = new Set(collapsedIds.value)
  if (next.has(id)) next.delete(id)
  else next.add(id)
  collapsedIds.value = next
}

function openCreate(parentId: string | null = null) {
  mutationError.value = ''
  initialSlug.value = ''
  const cached = queryClient.getQueryData<BodyPageCategory>(getGetApiCategoriesQueryKey({ page: 1, page_size: 100 }))
  initialUserId.value = cached?.data?.list?.find((item) => item.id === parentId)?.userId ?? ''
  dialogMode.value = 'create'
  editingRow.value = null
  targetParentId.value = parentId
  dialogInitialName.value = ''
}

function openEdit(row: CategoryRow) {
  mutationError.value = ''
  const cached = queryClient.getQueryData<BodyPageCategory>(getGetApiCategoriesQueryKey({ page: 1, page_size: 100 }))
  const item = cached?.data?.list?.find((item) => item.id === row.id)
  initialSlug.value = item?.slug ?? ''
  initialUserId.value = item?.userId ?? ''
  dialogMode.value = 'edit'
  editingRow.value = row
  targetParentId.value = row.parentId
  dialogInitialName.value = row.name
}

async function saveCategory(input: { name: string; slug: string; userId: string }) {
  if (isSaving.value) return
  mutationError.value = ''
  try {
    if (dialogMode.value === 'edit' && editingRow.value) {
      await updateMutation.mutateAsync({ id: editingRow.value.id, data: { name: input.name, slug: input.slug } })
    } else {
      await createMutation.mutateAsync({
        data: { ...input, ...(targetParentId.value ? { parentId: targetParentId.value } : {}) },
      })
    }
    dialogMode.value = null
    await queryClient.invalidateQueries({ queryKey: categoryQueryKey })
  } catch (error) {
    mutationError.value = error instanceof Error ? error.message : '保存失败，请重试'
  }
}

async function deleteCategory() {
  const row = pendingDelete.value
  if (!row || isDeleting.value) return
  mutationError.value = ''
  if (row.childCount) {
    mutationError.value = '请先删除该分类下的二级分类。'
    return
  }
  try {
    await deleteMutation.mutateAsync({ id: row.id })
    pendingDelete.value = null
    await queryClient.invalidateQueries({ queryKey: categoryQueryKey })
  } catch (error) {
    mutationError.value = error instanceof Error ? error.message : '删除失败，请重试'
  }
}
</script>

<template>
  <div
    v-if="categoriesQuery.isPending.value"
    class="mx-auto grid min-h-[50vh] content-start gap-4 px-3.5 pt-5 sm:px-6 sm:pt-7 lg:px-[42px]"
    aria-label="分类加载中"
  >
    <Skeleton class="h-8 w-48" />
    <Skeleton class="h-80 w-full rounded-lg" />
  </div>
  <div v-else-if="categoriesQuery.isError.value" class="grid min-h-[50vh] place-content-center gap-3 text-center">
    <p class="text-sm text-rose-300">分类加载失败，请检查连接后重试。</p>
    <button
      class="mx-auto rounded-md border border-slate-700 px-3 py-1.5 text-xs text-slate-200 disabled:opacity-50"
      :disabled="categoriesQuery.isFetching.value"
      @click="retry"
    >
      {{ categoriesQuery.isFetching.value ? '正在重试…' : '重试' }}
    </button>
  </div>
  <div v-else>
    <div class="mx-auto px-3.5 pt-4 pb-10 sm:px-6 sm:pt-7 lg:px-10.5">
      <div class="mb-5 flex flex-wrap items-end justify-between gap-4">
        <div class="flex items-center gap-3">
          <FolderTree class="text-violet-400" :size="24" />
          <div>
            <h1 class="text-xl font-bold">分类管理</h1>
            <p class="mt-1 text-xs text-slate-400">共 {{ categoryCount }} 个分类，最多两层</p>
          </div>
        </div>
        <div class="flex w-full flex-col gap-2 sm:w-auto sm:flex-row">
          <Input
            v-model="searchQuery"
            type="search"
            placeholder="搜索分类或上级分类"
            class="w-full sm:w-62.5"
            aria-label="搜索分类"
          >
            <template #leading><Search :size="15" /></template>
          </Input>
          <Button size="md" class="px-3 text-[11px] font-medium" @click="openCreate()"
            ><Plus :size="14" />新增分类</Button
          >
        </div>
      </div>

      <CategoryTable
        :rows="rows"
        :collapsed-ids="collapsedIds"
        :search-query="searchQuery"
        @toggle="toggleParent"
        @add-child="openCreate"
        @edit="openEdit"
        @delete="openDelete"
        @navigate="emit('navigate', $event)"
      />

      <CategoryEditorDialog
        v-if="dialogMode"
        :key="`${dialogMode}-${editingRow?.id ?? targetParentId ?? 'new'}`"
        :mode="dialogMode"
        :initial-name="dialogInitialName"
        :initial-slug="initialSlug"
        :initial-user-id="initialUserId"
        :busy="isSaving"
        :error="mutationError"
        :target-parent-name="categories.find((category) => category.id === targetParentId)?.name"
        :categories="categories"
        :editing-id="editingRow?.id"
        @close="!isSaving && (dialogMode = null)"
        @save="saveCategory"
      />
      <CategoryDeleteDialog
        v-if="pendingDelete"
        :name="pendingDelete.name"
        :busy="isDeleting"
        :error="mutationError"
        :is-parent="pendingDelete.level === '一级分类'"
        :child-count="pendingDelete.childCount"
        @close="!isDeleting && (pendingDelete = null)"
        @confirm="deleteCategory"
      />
    </div>
  </div>
</template>
