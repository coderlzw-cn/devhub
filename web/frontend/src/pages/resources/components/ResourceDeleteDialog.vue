<script setup lang="ts">
import { Trash2 } from 'lucide-vue-next'
import type { Resource } from '@/data/resources'
import Button from '@/components/ui/Button.vue'

defineProps<{ resource: Resource; loading?: boolean; error?: string }>()
const emit = defineEmits<{ close: []; confirm: [] }>()
</script>

<template>
  <div
    class="fixed inset-0 z-40 grid place-items-center bg-slate-950/70 p-4 backdrop-blur-sm"
    @click.self="emit('close')"
  >
    <section
      class="w-full max-w-[420px] rounded-lg border border-slate-700 bg-slate-900 p-5 shadow-2xl light:border-slate-200 light:bg-white"
      role="alertdialog"
      aria-modal="true"
      aria-labelledby="resource-delete-title"
    >
      <h2 id="resource-delete-title" class="text-base font-semibold text-slate-100 light:text-slate-800">删除资源</h2>
      <p class="mt-2 text-sm text-slate-400">确定删除“{{ resource.name }}”吗？此操作无法撤销。</p>
      <p v-if="error" role="alert" class="mt-3 text-sm text-rose-400">{{ error }}</p>
      <div class="mt-5 flex justify-end gap-2">
        <Button variant="secondary" size="sm" :disabled="loading" @click="emit('close')">取消</Button>
        <Button variant="danger" size="sm" :loading="loading" @click="emit('confirm')"
          ><Trash2 :size="14" />删除</Button
        >
      </div>
    </section>
  </div>
</template>
