import { getAllCategories, toCategoryTree } from '../src/api/category-tree'
import { afterEach, describe, expect, it, vi } from 'vitest'
import {
  getApiCategories,
  postApiCategories,
  patchApiCategoriesId,
  deleteApiCategoriesId,
} from '../src/api/generated/categories'
import { ApiError } from '../src/api/client'

const fetchMock = vi.fn()
vi.stubGlobal('window', { location: { origin: 'http://localhost' } })
vi.stubGlobal('fetch', fetchMock)
afterEach(() => fetchMock.mockReset())
const ok = (body: unknown) => new Response(JSON.stringify(body), { headers: { 'Content-Type': 'application/json' } })

describe('Orval category client', () => {
  it('serializes pagination and search, forwarding cancellation', async () => {
    fetchMock.mockResolvedValue(ok({ message: 'ok', data: { list: [], pages: 1 } }))
    const signal = new AbortController().signal
    await getApiCategories({ page: 2, page_size: 100, name: '前端' }, { signal })
    const [url, options] = fetchMock.mock.calls[0]!
    expect(new URL(url, 'http://localhost').searchParams.get('name')).toBe('前端')
    expect(url).toContain('page=2')
    expect(options.signal).toBe(signal)
  })
  it('sends create, update and delete to real API paths', async () => {
    fetchMock.mockImplementation(() => Promise.resolve(ok({ message: 'ok' })))
    const data = { userId: 'user', name: '分类', slug: 'category', parentId: 'parent' }
    await postApiCategories(data)
    await patchApiCategoriesId('category-id', { name: '新名称' })
    await deleteApiCategoriesId('category-id')
    expect(fetchMock.mock.calls.map(([url, options]) => [url, options.method])).toEqual([
      ['/api/categories', 'POST'],
      ['/api/categories/category-id', 'PATCH'],
      ['/api/categories/category-id', 'DELETE'],
    ])
    expect(JSON.parse(fetchMock.mock.calls[0]![1].body)).toEqual(data)
  })
  it('propagates backend deletion errors instead of reporting success', async () => {
    fetchMock.mockResolvedValue(
      new Response(JSON.stringify({ message: '分类下仍有站点' }), {
        status: 409,
        headers: { 'Content-Type': 'application/json' },
      }),
    )
    await expect(deleteApiCategoriesId('id')).rejects.toMatchObject({
      status: 409,
      message: '分类下仍有站点',
    } satisfies Partial<ApiError>)
  })
})

it('loads every page and assembles children even when they precede their parents', async () => {
  const parent = { id: 'parent', parentId: null, name: 'Parent' }
  const child = { id: 'child', parentId: 'parent', name: 'Child' }
  fetchMock.mockResolvedValueOnce(ok({ message: 'ok', data: { list: [child], pages: 2, page: 1 } }))
  fetchMock.mockResolvedValueOnce(ok({ message: 'ok', data: { list: [parent], pages: 2, page: 2 } }))
  const response = await getAllCategories()
  expect(fetchMock.mock.calls[1]![0]).toContain('page=2')
  expect(toCategoryTree(response.data.list)).toMatchObject([{ id: 'parent', children: [{ id: 'child' }] }])
})
