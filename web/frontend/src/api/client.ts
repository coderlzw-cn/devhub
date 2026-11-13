export type QueryValue = string | number | boolean | null | undefined
export type QueryParams = Record<string, QueryValue | readonly QueryValue[]>

export interface ApiRequestOptions extends Omit<RequestInit, 'body'> {
  params?: QueryParams
  body?: BodyInit | Record<string, unknown> | unknown[] | null
}

export class ApiError extends Error {
  readonly status: number
  readonly payload: unknown

  constructor(message: string, status: number, payload?: unknown) {
    super(message)
    this.name = 'ApiError'
    this.status = status
    this.payload = payload
  }
}

const baseUrl = import.meta.env.VITE_API_BASE_URL ?? ''

function buildUrl(path: string, params?: QueryParams) {
  const url = new URL(path, baseUrl || window.location.origin)
  if (params) {
    for (const [key, value] of Object.entries(params)) {
      const values = Array.isArray(value) ? value : [value]
      for (const item of values) {
        if (item !== null && item !== undefined) url.searchParams.append(key, String(item))
      }
    }
  }
  return baseUrl ? url.toString() : `${url.pathname}${url.search}${url.hash}`
}

function isBodyInit(body: ApiRequestOptions['body']): body is BodyInit {
  return (
    typeof body === 'string' ||
    body instanceof Blob ||
    body instanceof FormData ||
    body instanceof URLSearchParams ||
    body instanceof ArrayBuffer ||
    ArrayBuffer.isView(body)
  )
}

async function parseResponse(response: Response): Promise<unknown> {
  if (response.status === 204) return undefined
  const contentType = response.headers.get('content-type') ?? ''
  if (contentType.includes('application/json')) {
    try {
      return await response.json()
    } catch {
      return undefined
    }
  }
  const text = await response.text()
  return text || undefined
}

async function request<T>(path: string, options: ApiRequestOptions = {}): Promise<T> {
  const { params, body, headers, ...init } = options
  const requestHeaders = new Headers(headers)
  requestHeaders.set('Accept', 'application/json')

  let requestBody: BodyInit | null | undefined
  if (body === null || body === undefined) {
    requestBody = body
  } else if (isBodyInit(body)) {
    requestBody = body
  } else {
    requestHeaders.set('Content-Type', 'application/json')
    requestBody = JSON.stringify(body)
  }

  let response: Response
  try {
    response = await fetch(buildUrl(path, params), { ...init, headers: requestHeaders, body: requestBody })
  } catch (error) {
    if (error instanceof DOMException && error.name === 'AbortError') throw error
    throw new ApiError(error instanceof Error ? error.message : 'Network request failed', 0, error)
  }

  const payload = await parseResponse(response)
  if (!response.ok) {
    const message =
      typeof payload === 'object' && payload !== null && 'message' in payload && typeof payload.message === 'string'
        ? payload.message
        : response.statusText || 'Request failed'
    throw new ApiError(message, response.status, payload)
  }

  return payload as T
}

export const apiClient = {
  request,
  get<T>(path: string, options?: Omit<ApiRequestOptions, 'method' | 'body'>) {
    return request<T>(path, { ...options, method: 'GET' })
  },
  post<T>(path: string, body?: ApiRequestOptions['body'], options?: Omit<ApiRequestOptions, 'method' | 'body'>) {
    return request<T>(path, { ...options, method: 'POST', body })
  },
  put<T>(path: string, body?: ApiRequestOptions['body'], options?: Omit<ApiRequestOptions, 'method' | 'body'>) {
    return request<T>(path, { ...options, method: 'PUT', body })
  },
  patch<T>(path: string, body?: ApiRequestOptions['body'], options?: Omit<ApiRequestOptions, 'method' | 'body'>) {
    return request<T>(path, { ...options, method: 'PATCH', body })
  },
  delete<T>(path: string, options?: Omit<ApiRequestOptions, 'method' | 'body'>) {
    return request<T>(path, { ...options, method: 'DELETE' })
  },
}
