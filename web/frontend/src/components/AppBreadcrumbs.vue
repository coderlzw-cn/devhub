<script setup lang="ts">
import { ChevronRight } from 'lucide-vue-next'
import { RouterLink, type RouteLocationRaw } from 'vue-router'

defineProps<{
  items: { label: string; to?: RouteLocationRaw }[]
}>()
</script>

<template>
  <nav class="min-w-0" aria-label="面包屑">
    <ol class="flex min-w-0 items-center gap-1.5 text-xs sm:gap-2">
      <li v-for="(item, index) in items" :key="`${item.label}-${index}`" class="flex min-w-0 items-center gap-1.5 sm:gap-2">
        <ChevronRight v-if="index > 0" class="shrink-0 text-slate-600 light:text-slate-400" :size="14" aria-hidden="true" />
        <RouterLink
          v-if="item.to"
          :to="item.to"
          class="truncate text-slate-400 transition hover:text-slate-100 light:text-slate-500 light:hover:text-slate-800"
        >
          {{ item.label }}
        </RouterLink>
        <span
          v-else
          class="truncate font-medium text-slate-100 light:text-slate-800"
          aria-current="page"
        >
          {{ item.label }}
        </span>
      </li>
    </ol>
  </nav>
</template>
