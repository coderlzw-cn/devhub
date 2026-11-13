import { computed } from 'vue'
import { useGetApiSites, getGetApiSitesQueryKey } from '@/api/generated/sites'
import { getAllSites, toResource } from '@/api/resources'
import { useResourceCategoriesQuery } from './useResourceCategoriesQuery'
import type { BodyPageSite } from '@/api/generated/site-models'

export const resourcesQueryKey = getGetApiSitesQueryKey()

export function useResourcesQuery() {
  const categories = useResourceCategoriesQuery()
  return useGetApiSites(undefined, {
    query: {
      queryFn: ({ signal }) => getAllSites(signal),
      select: computed(() => {
        const all = (categories.data.value ?? []).flatMap((parent) => [parent, ...parent.children])
        return (response: BodyPageSite) => (response.data?.list ?? []).map((site) => toResource(site, all))
      }),
    },
  })
}
