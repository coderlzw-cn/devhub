import { defineConfig } from 'orval'

export default defineConfig({
  categories: {
    input: { target: './openapi/backend.json', filters: { tags: ['分类'] } },
    output: {
      target: './src/api/generated/categories.ts',
      schemas: './src/api/generated/models',
      client: 'vue-query',
      httpClient: 'fetch',
      override: {
        mutator: { path: './src/api/orval-fetch.ts', name: 'orvalFetch' },
        query: { version: 5 },
        fetch: { includeHttpResponseReturnType: false },
      },
    },
  },
  sites: {
    input: { target: './openapi/backend.json', filters: { tags: ['站点'] } },
    output: {
      target: './src/api/generated/sites.ts',
      schemas: './src/api/generated/site-models',
      client: 'vue-query',
      httpClient: 'fetch',
      override: {
        mutator: { path: './src/api/orval-fetch.ts', name: 'orvalFetch' },
        query: { version: 5 },
        fetch: { includeHttpResponseReturnType: false },
      },
    },
  },
})
