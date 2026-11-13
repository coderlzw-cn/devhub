import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { VueQueryPlugin, QueryClient } from '@tanstack/vue-query'
import persistedState from 'pinia-plugin-persistedstate'
import { PiniaSharedState } from 'pinia-shared-state'
import App from '@/App.vue'
import { router } from '@/router'
import '@/style.css'

async function bootstrap() {
  // Older development builds may have left an MSW service worker controlling this origin.
  if ('serviceWorker' in navigator) {
    const registrations = await navigator.serviceWorker.getRegistrations()
    await Promise.all(
      registrations
        .filter((registration) =>
          [registration.active, registration.waiting, registration.installing].some(
            (worker) => worker && new URL(worker.scriptURL).pathname.endsWith('/mockServiceWorker.js'),
          ),
        )
        .map((registration) => registration.unregister()),
    )
  }

  const pinia = createPinia()
  pinia.use(persistedState)
  pinia.use(PiniaSharedState({ enable: false, type: 'localstorage' }))

  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        staleTime: 60_000,
        retry: 1,
      },
    },
  })

  createApp(App).use(pinia).use(router).use(VueQueryPlugin, { queryClient }).mount('#app')
}

void bootstrap()
