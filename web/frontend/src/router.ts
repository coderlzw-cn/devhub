import { createRouter, createWebHistory } from 'vue-router'
import MainLayout from '@/layouts/MainLayout.vue'
import DevHubLoginPage from '@/pages/DevHubLoginPage.vue'
import DashboardPage from '@/pages/DashboardPage.vue'
import ResourcesPage from '@/pages/ResourcesPage.vue'
import CategoriesPage from '@/pages/CategoriesPage.vue'

const routes = [
  { path: '/login', name: 'login', component: DevHubLoginPage },
  {
    path: '/',
    component: MainLayout,
    children: [
      { path: '', name: 'home', component: DashboardPage },
      { path: 'resources', name: 'resources', component: ResourcesPage },
      { path: 'categories', name: 'categories', component: CategoriesPage },
      { path: 'profile', name: 'profile', component: DashboardPage },
    ],
  },
  { path: '/:pathMatch(.*)*', redirect: '/' },
]

export const router = createRouter({
  history: createWebHistory(),
  routes,
})
