import { getAllCategories, toCategoryTree } from '@/api/category-tree'

export interface ResourceCategory {
  id: string
  name: string
  icon: string
  children: ResourceCategoryLeaf[]
}

export interface ResourceCategoryLeaf {
  id: string
  name: string
  children: []
}

export type AnyResourceCategory = ResourceCategory | ResourceCategoryLeaf

export async function getResourceCategories(signal?: AbortSignal) {
  const response = await getAllCategories(signal)
  return toCategoryTree(response.data.list)
}
