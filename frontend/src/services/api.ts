import router from '@/router'
import { useAuthStore } from '@/stores/useAuthStore'

const BASE_URL = '/api'

interface ApiResponse<T> {
  data: T
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
async function request<T>(method: string, url: string, body?: unknown, params?: Record<string, any>): Promise<ApiResponse<T>> {
  let fullUrl = `${BASE_URL}${url}`

  if (params) {
    const searchParams = new URLSearchParams()
    for (const [key, value] of Object.entries(params)) {
      if (value !== undefined && value !== null) {
        searchParams.append(key, String(value))
      }
    }
    const qs = searchParams.toString()
    if (qs) fullUrl += `?${qs}`
  }

  const options: RequestInit = {
    method,
    credentials: 'include',
    headers: body !== undefined ? { 'Content-Type': 'application/json' } : {},
    body: body !== undefined ? JSON.stringify(body) : undefined,
  }

  const response = await fetch(fullUrl, options)

  if (!response.ok) {
    if (response.status === 401 && !url.includes('/auth/me')) {
      const auth = useAuthStore()
      auth.user = null
      router.push('/landing')
    }
    throw new Error(`${response.status} ${response.statusText}`)
  }

  const contentType = response.headers.get('content-type')
  const data = contentType?.includes('application/json') ? await response.json() : null

  return { data: data as T }
}

const api = {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  get<T>(url: string, config?: { params?: Record<string, any> }): Promise<ApiResponse<T>> {
    return request<T>('GET', url, undefined, config?.params)
  },
  post<T>(url: string, body?: unknown): Promise<ApiResponse<T>> {
    return request<T>('POST', url, body)
  },
  put<T>(url: string, body?: unknown): Promise<ApiResponse<T>> {
    return request<T>('PUT', url, body)
  },
  delete<T>(url: string): Promise<ApiResponse<T>> {
    return request<T>('DELETE', url)
  },
}

export default api
