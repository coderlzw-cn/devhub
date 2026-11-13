<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useQueryClient } from '@tanstack/vue-query'
import { Plus } from 'lucide-vue-next'
import ResourceCard from '@/components/ResourceCard.vue'
import Button from '@/components/ui/Button.vue'
import Input from '@/components/ui/Input.vue'
import Select from '@/components/ui/Select.vue'
import Pagination from '@/components/ui/Pagination.vue'
import type { Resource } from '@/data/resources'
import { usePostApiSites, usePatchApiSitesId, useDeleteApiSitesId, getGetApiSitesQueryKey } from '@/api/generated/sites'
import ResourceEditorDialog, { type ResourceFormValues } from '@/pages/resources/components/ResourceEditorDialog.vue'
import ResourceDeleteDialog from '@/pages/resources/components/ResourceDeleteDialog.vue'
import { useResourcesQuery } from '@/composables/useResourcesQuery'
import { useResourceCategoriesQuery } from '@/composables/useResourceCategoriesQuery'
import Skeleton from '@/components/ui/Skeleton.vue'
import { findCategory, getCategoryNames } from '@/data/category-tree'
const props = defineProps<{
  activeCategory: string
}>()
const resourcesQuery = useResourcesQuery()
const categoriesQuery = useResourceCategoriesQuery()
const queryClient = useQueryClient()
const resources = computed(() => resourcesQuery.data.value ?? [])
const isPending = computed(() => resourcesQuery.isPending.value || categoriesQuery.isPending.value)
const isError = computed(() => resourcesQuery.isError.value || categoriesQuery.isError.value)
const isFetching = computed(() => resourcesQuery.isFetching.value || categoriesQuery.isFetching.value)
const emit = defineEmits<{
  (event: 'open', resource: Resource): void
  (event: 'category', category: string): void
}>()
const sortBy = ref('默认排序')
const search = ref('')
const currentPage = ref(1)
const pageSize = ref(9)
const editorMode = ref<'create' | 'edit' | null>(null)
const editingResource = ref<Resource | null>(null)
const pendingDelete = ref<Resource | null>(null)
const submitError = ref('')
const createMutation = usePostApiSites()
const updateMutation = usePatchApiSitesId()
const deleteMutation = useDeleteApiSitesId()
const isSubmitting = computed(() => createMutation.isPending.value || updateMutation.isPending.value)
const deleteError = ref('')
const selectedParent = computed(() =>
  categoriesQuery.data.value?.find(
    (category) =>
      category.name === props.activeCategory || category.children.some((child) => child.name === props.activeCategory),
  ),
)
const categoryOptions = computed(() => {
  const parent = selectedParent.value
  if (!parent) return [{ label: '全部分类', value: '全部' }]
  if (parent.children.length === 0) return [{ label: parent.name, value: parent.name }]
  return [
    { label: `全部 ${parent.name}`, value: parent.name },
    ...parent.children.map((child) => ({ label: child.name, value: child.name })),
  ]
})
const selectedCategory = computed(() => {
  if (selectedParent.value?.children.some((child) => child.name === props.activeCategory)) return props.activeCategory
  return selectedParent.value?.name ?? '全部'
})
const resourceCategoryOptions = computed(() => {
  const categories = categoriesQuery.data.value ?? []
  return categories.flatMap((category) => [
    { label: category.name, value: category.id },
    ...category.children.map((child) => ({ label: `${category.name} / ${child.name}`, value: child.id })),
  ])
})
function openAddDialog() {
  editorMode.value = 'create'
  editingResource.value = null
  submitError.value = ''
}
function openEditDialog(resource: Resource) {
  editorMode.value = 'edit'
  editingResource.value = resource
  submitError.value = ''
}
function closeEditor() {
  if (isSubmitting.value) return
  editorMode.value = null
  editingResource.value = null
}
async function invalidateSites() {
  await queryClient.invalidateQueries({ queryKey: getGetApiSitesQueryKey() })
}
async function submitResource(values: ResourceFormValues) {
  if (isSubmitting.value) return
  submitError.value = ''
  try {
    const data = {
      title: values.name.trim(),
      description: values.description.trim(),
      url: new URL(values.url).toString(),
      categoryId: values.category || undefined,
      logoUrl: values.logoUrl.trim(),
      isPinned: values.isPinned,
      sortOrder: values.sortOrder,
    }
    if (editingResource.value) {
      await updateMutation.mutateAsync({
        id: editingResource.value.id,
        data,
      })
    } else {
      await createMutation.mutateAsync({ data: { ...data, userId: values.userId.trim() } })
    }
    editorMode.value = null
    editingResource.value = null
    currentPage.value = 1
    await invalidateSites()
  } catch (error) {
    submitError.value = error instanceof Error ? error.message : '站点保存失败，请重试'
  }
}
function openDeleteDialog(resource: Resource) {
  deleteError.value = ''
  pendingDelete.value = resource
}
async function deleteResource() {
  if (!pendingDelete.value || deleteMutation.isPending.value) return
  deleteError.value = ''
  try {
    await deleteMutation.mutateAsync({ id: pendingDelete.value.id })
    pendingDelete.value = null
    await invalidateSites()
  } catch (error) {
    deleteError.value = error instanceof Error ? error.message : '站点删除失败，请重试'
  }
}
function retryQueries() {
  void Promise.all([resourcesQuery.refetch(), categoriesQuery.refetch()])
}
const filtered = computed(() => {
  const selectedCategory =
    props.activeCategory === '全部' ? undefined : findCategory(categoriesQuery.data.value ?? [], props.activeCategory)
  const selectedCategoryNames = selectedCategory
    ? 'icon' in selectedCategory
      ? getCategoryNames(selectedCategory)
      : [selectedCategory.name]
    : undefined
  const matches = resources.value.filter((resource) => {
    const matchesCategory = !selectedCategoryNames || selectedCategoryNames.includes(resource.category)
    const query = search.value.trim().toLocaleLowerCase()
    const matchesSearch =
      !query ||
      [resource.name, resource.url, resource.description].some((value) => value.toLocaleLowerCase().includes(query))
    return matchesCategory && matchesSearch
  })
  return sortBy.value === '名称' ? [...matches].sort((a, b) => a.name.localeCompare(b.name)) : matches
})
watch([filtered, pageSize], () => {
  currentPage.value = Math.min(currentPage.value, Math.max(1, Math.ceil(filtered.value.length / pageSize.value)))
})
watch([() => props.activeCategory, search, sortBy], () => {
  currentPage.value = 1
})
const visibleResources = computed(() =>
  filtered.value.slice((currentPage.value - 1) * pageSize.value, currentPage.value * pageSize.value),
)
</script>
<template>
  <div
    v-if="isPending"
    class="mx-auto grid min-h-[50vh] max-w-[1600px] content-start gap-5 px-3.5 pt-5 sm:px-6 sm:pt-7 lg:px-[42px]"
    aria-label="资源列表加载中"
  >
    <Skeleton class="h-8 w-48" />
    <div class="grid grid-cols-1 gap-2 sm:grid-cols-2 xl:grid-cols-3">
      <Skeleton v-for="item in 9" :key="item" class="h-36.25 rounded-xl" />
    </div>
  </div>
  <div v-else-if="isError" class="grid min-h-[50vh] place-content-center gap-3 text-center">
    <p class="text-sm text-rose-300">资源或分类加载失败，请检查连接后重试。</p>
    <button
      class="mx-auto rounded-md border border-slate-700 px-3 py-1.5 text-xs text-slate-200 disabled:opacity-50"
      :disabled="isFetching"
      @click="retryQueries"
    >
      {{ isFetching ? '正在重试…' : '重试' }}
    </button>
  </div>
  <div v-else>
    <div class="mx-auto max-w-[1600px] px-3.25 pt-3.25 pb-10 sm:px-6 sm:pt-5 sm:pb-12.5 lg:px-[42px]">
      <div class="flex flex-col items-start justify-between gap-3 pb-4 sm:flex-row sm:items-center">
        <div class="flex flex-wrap items-center gap-x-3 gap-y-2 sm:gap-4.5">
          <Input
            v-model="search"
            type="search"
            placeholder="搜索站点名称、网址或描述"
            aria-label="搜索站点"
            class="w-full sm:w-60"
          />
          <label
            v-if="selectedParent"
            class="flex items-center gap-2 text-[10px] whitespace-nowrap text-slate-300 sm:text-xs light:text-slate-600"
            >二级分类：<Select
              :model-value="selectedCategory"
              :options="categoryOptions"
              size="sm"
              class="min-w-[130px] sm:min-w-[170px]"
              @update:model-value="emit('category', String($event))"
          /></label>
          <label
            v-else
            class="flex items-center gap-2 text-[10px] whitespace-nowrap text-slate-300 sm:text-xs light:text-slate-600"
            >分类：<Select
              :model-value="'全部'"
              :options="[{ label: '全部分类', value: '全部' }]"
              size="sm"
              class="min-w-[130px] sm:min-w-[170px]"
          /></label>
          <label
            class="flex items-center gap-2 text-[10px] whitespace-nowrap text-slate-300 sm:text-xs light:text-slate-600"
            >排序：<Select v-model="sortBy" :options="['默认排序', '名称']" size="sm" class="min-w-25 sm:min-w-32.5"
          /></label>
        </div>
        <Button class="self-end sm:self-auto" size="sm" :disabled="isPending" @click="openAddDialog">
          <Plus :size="15" />添加资源
        </Button>
      </div>
      <div class="mb-3 hidden gap-3 text-[11px] text-slate-400">
        <span>共 {{ filtered.length }} 个资源</span
        ><span v-if="activeCategory !== '全部'">分类：{{ activeCategory }}</span>
      </div>
      <div class="grid grid-cols-1 gap-2 sm:grid-cols-2 xl:grid-cols-3">
        <ResourceCard
          v-for="resource in visibleResources"
          :key="resource.id"
          :resource="resource"
          @open="emit('open', $event)"
          @edit="openEditDialog"
          @delete="openDeleteDialog"
        />
        <div
          v-if="visibleResources.length === 0"
          class="col-span-full grid min-h-48 place-items-center rounded-xl border border-dashed border-slate-700 bg-slate-900/40 px-6 text-center text-sm text-slate-400"
        >
          没有找到匹配的资源，请调整筛选条件。
        </div>
      </div>
      <Pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :total="filtered.length"
        :page-size-options="[9, 18, 36]"
      />
      <ResourceEditorDialog
        v-if="editorMode"
        :key="`${editorMode}-${editingResource?.id ?? 'new'}`"
        :mode="editorMode"
        :resource="editingResource"
        :category-options="resourceCategoryOptions"
        :initial-category="
          resourceCategoryOptions.find(
            (option) => option.label === activeCategory || option.label.endsWith(' / ' + activeCategory),
          )?.value ?? ''
        "
        :loading="isSubmitting"
        :error="submitError"
        @close="closeEditor"
        @submit="submitResource"
      />
      <ResourceDeleteDialog
        v-if="pendingDelete"
        :resource="pendingDelete"
        :loading="deleteMutation.isPending.value"
        :error="deleteError"
        @close="!deleteMutation.isPending.value && (pendingDelete = null)"
        @confirm="deleteResource"
      />
    </div>
  </div>
</template>
