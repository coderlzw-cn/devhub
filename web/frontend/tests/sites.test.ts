import { afterEach, describe, expect, it, vi } from 'vitest'
import { getApiSites, postApiSites, patchApiSitesId, deleteApiSitesId } from '../src/api/generated/sites'
import { getAllSites, toResource } from '../src/api/resources'
import type { Site } from '../src/api/generated/site-models'

const fetchMock = vi.fn()
vi.stubGlobal('window', { location: { origin: 'http://localhost' } })
vi.stubGlobal('fetch', fetchMock)
afterEach(() => fetchMock.mockReset())
const ok = (body: unknown) => new Response(JSON.stringify(body), { headers: { 'Content-Type': 'application/json' } })

describe('Orval site client', () => {
  it('serializes pagination, category and false-valued filters and passes cancellation', async () => {
    fetchMock.mockResolvedValue(ok({ message: 'ok', data: { list: [], pages: 1 } }))
    const signal = new AbortController().signal
    await getApiSites(
      { page: 2, page_size: 100, title: '开发', category_id: 'category-id', is_pinned: false, is_archived: false },
      { signal },
    )
    const [url, options] = fetchMock.mock.calls[0]!
    const params = new URL(url, 'http://localhost').searchParams
    expect(params.get('title')).toBe('开发')
    expect(params.get('category_id')).toBe('category-id')
    expect(params.get('is_pinned')).toBe('false')
    expect(params.get('is_archived')).toBe('false')
    expect(params.get('page')).toBe('2')
    expect(options.signal).toBe(signal)
  })
  it('persists create and edit fields and deletes by the real string ID', async () => {
    fetchMock.mockImplementation(() => Promise.resolve(ok({ message: 'ok' })))
    const input = {
      userId: 'owner-id',
      title: '站点',
      url: 'https://example.com',
      categoryId: 'category-id',
      logoUrl: 'https://example.com/icon.png',
      sortOrder: 12,
      isPinned: true,
    }
    await postApiSites(input)
    await patchApiSitesId('site-id', { title: '新名称', isPinned: false, isArchived: true })
    await deleteApiSitesId('site-id')
    expect(fetchMock.mock.calls.map(([url, options]) => [url, options.method])).toEqual([
      ['/api/sites', 'POST'],
      ['/api/sites/site-id', 'PATCH'],
      ['/api/sites/site-id', 'DELETE'],
    ])
    expect(JSON.parse(fetchMock.mock.calls[0]![1].body)).toEqual(input)
    expect(JSON.parse(fetchMock.mock.calls[1]![1].body)).toMatchObject({ isPinned: false, isArchived: true })
  })
  it('loads all pages including null lists without reading legacy localStorage data', async () => {
    fetchMock.mockResolvedValueOnce(ok({ message: 'ok', data: { list: [{ id: 'first' }], pages: 3, total: 2 } }))
    fetchMock.mockResolvedValueOnce(ok({ message: 'ok', data: { list: null, pages: 3 } }))
    fetchMock.mockResolvedValueOnce(ok({ message: 'ok', data: { list: [{ id: 'last' }], pages: 3 } }))
    const result = await getAllSites()
    expect(result.data.list.map((site) => site.id)).toEqual(['first', 'last'])
    expect(fetchMock.mock.calls[2]![0]).toContain('page=3')
  })
  it('surfaces API errors and missing page data instead of silently showing an empty list', async () => {
    fetchMock.mockResolvedValueOnce(
      new Response(JSON.stringify({ message: '没有权限' }), {
        status: 403,
        headers: { 'Content-Type': 'application/json' },
      }),
    )
    await expect(deleteApiSitesId('id')).rejects.toMatchObject({ status: 403, message: '没有权限' })
    fetchMock.mockResolvedValueOnce(ok({ message: '站点不可用' }))
    await expect(getAllSites()).rejects.toThrow('站点不可用')
  })
  it('preserves backend identifiers and metadata while resolving category names', () => {
    const site = {
      id: 'site-id',
      title: '站点',
      categoryId: 'category-id',
      userId: 'owner-id',
      isPinned: true,
      isArchived: false,
      logoUrl: 'https://example.com/logo.png',
    } as Site
    expect(toResource(site, [{ id: 'category-id', name: '分类' }])).toMatchObject({
      ...site,
      name: '站点',
      category: '分类',
    })
    expect(toResource({ ...site, categoryId: null }, []).category).toBe('未分类')
  })
})
