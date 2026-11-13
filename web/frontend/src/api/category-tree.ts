import { getApiCategories } from './generated/categories'
import type { Category } from './generated/models'
import type { ResourceCategory } from './resource-categories'

export function toCategoryTree(items: Category[]): ResourceCategory[] {
  return items
    .filter((item) => !item.parentId || !items.some((parent) => parent.id === item.parentId))
    .map((item) => ({
      id: item.id,
      name: item.name,
      icon: item.icon || 'FolderTree',
      children: items
        .filter((child) => child.parentId === item.id)
        .map((child) => ({ id: child.id, name: child.name, children: [] })),
    }))
}

export async function getAllCategories(signal?: AbortSignal) {
  const first = await getApiCategories({ page: 1, page_size: 100 }, { signal })
  if (!first.data) throw new Error(first.message || '分类加载失败')
  const items: Category[] = [...(first.data.list ?? [])]
  for (let page = 2; page <= first.data.pages; page++) {
    const next = await getApiCategories({ page, page_size: 100 }, { signal })
    if (!next.data) throw new Error(next.message || '分类加载失败')
    items.push(...(next.data.list ?? []))
  }
  return { ...first, data: { ...first.data, list: items } }
}
