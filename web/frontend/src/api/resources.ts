import { getApiSites } from '@/api/generated/sites'
import { getAllCategories } from '@/api/category-tree'
import type { Category } from '@/api/generated/models'
import type { Site } from '@/api/generated/site-models'
import type { Resource } from '@/data/resources'

export async function getAllSites(signal?: AbortSignal) {
  const first = await getApiSites({ page: 1, page_size: 100 }, { signal })
  if (!first.data) throw new Error(first.message || '站点加载失败')
  const sites = [...(first.data.list ?? [])]
  for (let page = 2; page <= first.data.pages; page++) {
    const next = await getApiSites({ page, page_size: 100 }, { signal })
    if (!next.data) throw new Error(next.message || '站点加载失败')
    sites.push(...(next.data.list ?? []))
  }
  return { ...first, data: { ...first.data, list: sites } }
}

export function toResource(site: Site, categories: Pick<Category, 'id' | 'name'>[]): Resource {
  return {
    ...site,
    name: site.title,
    category: categories.find((item) => item.id === site.categoryId)?.name ?? '未分类',
    icon: site.title.slice(0, 2),
    color: '#8b5cf6',
  }
}

export async function getResources(signal?: AbortSignal): Promise<Resource[]> {
  const [sites, categories] = await Promise.all([getAllSites(signal), getAllCategories(signal)])
  return sites.data.list.map((site) => toResource(site, categories.data.list))
}
