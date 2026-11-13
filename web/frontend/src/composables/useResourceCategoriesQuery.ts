import { useGetApiCategories } from '@/api/generated/categories'
import { getAllCategories, toCategoryTree } from '@/api/category-tree'

export function useResourceCategoriesQuery() {
  return useGetApiCategories(
    { page: 1, page_size: 100 },
    {
      query: {
        queryFn: ({ signal }) => getAllCategories(signal),
        select: (response) => toCategoryTree(response.data?.list ?? []),
      },
    },
  )
}
