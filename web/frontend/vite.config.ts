import { defineConfig } from 'vite'
import { fileURLToPath, URL } from 'node:url'
import vue from '@vitejs/plugin-vue'
import tailwindcss from '@tailwindcss/vite'

export default defineConfig({
  root: fileURLToPath(new URL('.', import.meta.url)),
  build: { outDir: fileURLToPath(new URL('../dist', import.meta.url)), emptyOutDir: true },
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
  plugins: [vue(), tailwindcss()],
  server: { host: '0.0.0.0', port: 5173, proxy: { '/api': 'http://127.0.0.1:8080' } },
})
