import { useQuery } from '@tanstack/vue-query'
import type { ProfileData } from '@/data/profile'

function isProfileData(value: unknown): value is ProfileData {
  return (
    typeof value === 'object' &&
    value !== null &&
    'name' in value &&
    typeof value.name === 'string' &&
    'email' in value &&
    typeof value.email === 'string' &&
    'bio' in value &&
    typeof value.bio === 'string' &&
    'avatar' in value &&
    typeof value.avatar === 'string' &&
    'joinedAt' in value &&
    typeof value.joinedAt === 'string'
  )
}

function getPersistedProfile(): ProfileData | undefined {
  const saved = localStorage.getItem('devhub-local-profile-v2')
  if (!saved) return undefined

  try {
    const parsed: unknown = JSON.parse(saved)
    if (isProfileData(parsed)) return parsed
  } catch {
    return undefined
  }
  return undefined
}

export function useProfileQuery() {
  return useQuery({
    queryKey: profileQueryKey,
    queryFn: () => getPersistedProfile() ?? { name: '', email: '', bio: '', avatar: '', joinedAt: '' },
  })
}

export const profileQueryKey = ['profile'] as const
