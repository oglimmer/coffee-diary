import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

// Mock the router and auth store so the 401 branch can be asserted in isolation.
// vi.mock is hoisted, so shared state must go through vi.hoisted.
const { pushMock, authState } = vi.hoisted(() => ({
  pushMock: vi.fn(),
  authState: { user: 'someone' as unknown },
}))
vi.mock('@/router', () => ({
  default: { push: pushMock },
}))
vi.mock('@/stores/useAuthStore', () => ({
  useAuthStore: () => authState,
}))

import api from './api'

function mockResponse(
  body: unknown,
  { status = 200, statusText = 'OK', contentType = 'application/json' } = {},
) {
  return {
    ok: status >= 200 && status < 300,
    status,
    statusText,
    headers: { get: (h: string) => (h === 'content-type' ? contentType : null) },
    json: async () => body,
  } as unknown as Response
}

describe('api request builder', () => {
  let fetchMock: ReturnType<typeof vi.fn>

  beforeEach(() => {
    fetchMock = vi.fn()
    vi.stubGlobal('fetch', fetchMock)
    pushMock.mockClear()
    authState.user = 'someone'
  })

  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('prefixes the base url and defaults to no body/content-type on GET', async () => {
    fetchMock.mockResolvedValue(mockResponse({ ok: true }))

    const res = await api.get('/coffees')

    expect(fetchMock).toHaveBeenCalledTimes(1)
    const [url, options] = fetchMock.mock.calls[0]
    expect(url).toBe('/api/coffees')
    expect(options.method).toBe('GET')
    expect(options.credentials).toBe('include')
    expect(options.body).toBeUndefined()
    expect(options.headers).toEqual({})
    expect(res.data).toEqual({ ok: true })
  })

  it('serialises query params and skips null/undefined values', async () => {
    fetchMock.mockResolvedValue(mockResponse({ content: [] }))

    await api.get('/diary-entries', {
      params: { page: 0, size: 20, coffeeId: null, sieveId: undefined, minRating: 3 },
    })

    const url = fetchMock.mock.calls[0][0] as string
    expect(url).toBe('/api/diary-entries?page=0&size=20&minRating=3')
  })

  it('sends a JSON body with content-type header on POST', async () => {
    fetchMock.mockResolvedValue(mockResponse({ id: 1 }))

    const payload = { name: 'Ethiopia' }
    await api.post('/coffees', payload)

    const [url, options] = fetchMock.mock.calls[0]
    expect(url).toBe('/api/coffees')
    expect(options.method).toBe('POST')
    expect(options.headers).toEqual({ 'Content-Type': 'application/json' })
    expect(options.body).toBe(JSON.stringify(payload))
  })

  it('returns null data when the response is not JSON (e.g. 204 delete)', async () => {
    fetchMock.mockResolvedValue(mockResponse(null, { contentType: 'text/plain' }))

    const res = await api.delete('/coffees/1')

    expect(res.data).toBeNull()
  })

  it('throws with status and statusText on a non-ok response', async () => {
    fetchMock.mockResolvedValue(mockResponse(null, { status: 500, statusText: 'Server Error' }))

    await expect(api.get('/coffees')).rejects.toThrow('500 Server Error')
  })

  it('clears the user and redirects to /landing on a 401', async () => {
    fetchMock.mockResolvedValue(mockResponse(null, { status: 401, statusText: 'Unauthorized' }))

    await expect(api.get('/diary-entries')).rejects.toThrow('401 Unauthorized')
    expect(authState.user).toBeNull()
    expect(pushMock).toHaveBeenCalledWith('/landing')
  })

  it('does not redirect on a 401 from /auth/me (session probe)', async () => {
    fetchMock.mockResolvedValue(mockResponse(null, { status: 401, statusText: 'Unauthorized' }))

    await expect(api.get('/auth/me')).rejects.toThrow('401 Unauthorized')
    expect(pushMock).not.toHaveBeenCalled()
  })
})
