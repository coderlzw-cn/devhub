import type { AnyResourceCategory, ResourceCategory } from '@/api/resource-categories'

export function getCategoryNames(category: ResourceCategory): string[] {
  return [category.name, ...category.children.map((child) => child.name)]
}

export function findCategory(categories: ResourceCategory[], name: string): AnyResourceCategory | undefined {
  for (const category of categories) {
    if (category.name === name) return category
    const childMatch = category.children.find((child) => child.name === name)
    if (childMatch) return childMatch
  }
  return undefined
}
