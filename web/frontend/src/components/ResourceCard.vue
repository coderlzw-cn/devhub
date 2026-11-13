<script setup lang="ts">
import { ExternalLink, MoreHorizontal, Pencil, Trash2 } from 'lucide-vue-next'
import Dropdown from '@/components/ui/Dropdown.vue'
import type { Resource } from '@/data/resources'
defineProps<{ resource: Resource }>()
const emit = defineEmits<{
  (event: 'open', resource: Resource): void
  (event: 'edit', resource: Resource): void
  (event: 'delete', resource: Resource): void
}>()
function openResource(resource: Resource) {
  emit('open', resource)
}
</script>
<template>
  <article
    class="relative flex min-h-[120px] min-w-0 items-start gap-[15px] rounded-xl border border-slate-700 bg-linear-to-br from-slate-800 to-slate-900 p-4 transition hover:-translate-y-0.5 hover:border-violet-500 hover:shadow-lg sm:p-[19px] light:border-slate-200! light:bg-white! light:bg-none! light:hover:bg-violet-50!"
  >
    <div
      class="grid size-[46px] shrink-0 place-items-center rounded-[11px] text-2xl font-extrabold sm:size-[57px]"
      :style="{
        color: resource.color,
        backgroundColor: `color-mix(in srgb, ${resource.color} 18%, transparent)`,
      }"
    >
      <img v-if="resource.logoUrl" :src="resource.logoUrl" alt="" class="size-full rounded-[11px] object-contain" />
      <span v-else>{{ resource.icon }}</span>
    </div>
    <div class="min-w-0 flex-1 pt-0.5 pr-5">
      <strong class="text-[13px] font-bold text-slate-100 sm:text-[15px] light:text-slate-800">{{
        resource.name
      }}</strong>
      <p
        class="my-[7px] mb-2.5 line-clamp-2 min-h-[17px] text-[11px] leading-normal text-slate-400 light:text-slate-500"
      >
        {{ resource.description }}
      </p>
      <div class="flex flex-wrap gap-1.5">
        <span
          class="rounded-full bg-indigo-950 px-2 py-1 text-[9px] whitespace-nowrap text-indigo-200 light:bg-indigo-50 light:text-indigo-700"
          >{{ resource.category }}</span
        >
      </div>
    </div>
    <div class="absolute inset-y-2.5 right-2.5 flex flex-col items-center justify-between">
      <button
        class="grid size-7 place-items-center rounded-md bg-slate-800 text-indigo-200 transition hover:bg-slate-700 hover:text-white focus-visible:ring-2 focus-visible:ring-violet-500 focus-visible:outline-none light:bg-slate-100 light:text-indigo-700 light:hover:bg-indigo-700 light:hover:text-white light:focus-visible:ring-indigo-600"
        aria-label="打开资源"
        @click="openResource(resource)"
      >
        <ExternalLink :size="17" />
      </button>
      <Dropdown
        :items="[
          { label: '编辑资源', value: 'edit', icon: Pencil },
          { label: '删除资源', value: 'delete', icon: Trash2 },
        ]"
        align="end"
        fit
        trigger-label="资源操作菜单"
        @select="(value) => (value === 'edit' ? emit('edit', resource) : emit('delete', resource))"
      >
        <template #default="{ open }">
          <span
            class="grid size-7 place-items-center rounded-md text-slate-400 transition hover:bg-slate-700 hover:text-white light:hover:bg-slate-100 light:hover:text-slate-800"
            :class="open ? 'bg-slate-700 text-white light:bg-slate-100 light:text-slate-800' : ''"
            aria-hidden="true"
          >
            <MoreHorizontal :size="17" />
          </span>
        </template>
      </Dropdown>
    </div>
  </article>
</template>
