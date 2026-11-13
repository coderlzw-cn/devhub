import { apiClient, ApiError } from '@/api/client'

export function orvalFetch<T>(url: string, options: RequestInit = {}): Promise<T> {
  return apiClient.request<T>(url, options)
}

export type ErrorType<T> = ApiError & { payload: T }
