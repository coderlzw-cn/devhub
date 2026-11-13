import type { Site } from '@/api/generated/site-models'

export type Resource = Site & {
  name: string
  category: string
  icon: string
  color: string
}
