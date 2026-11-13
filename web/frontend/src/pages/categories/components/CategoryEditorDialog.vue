<script setup lang="ts">
import { Field, useForm } from 'vee-validate'
import { X } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import Input from '@/components/ui/Input.vue'
import type { ResourceCategory } from '@/api/resource-categories'

const props = defineProps<{
  mode: 'create' | 'edit'
  initialName: string
  initialSlug: string
  initialUserId: string
  busy: boolean
  error: string
  targetParentName?: string
  categories: ResourceCategory[]
  editingId?: string
}>()
const emit = defineEmits<{ close: []; save: [input: { name: string; slug: string; userId: string }] }>()
const { handleSubmit } = useForm<{ name: string; slug: string; userId: string }>({
  initialValues: { name: props.initialName, slug: props.initialSlug, userId: props.initialUserId },
})

function validateName(input: unknown) {
  const name = String(input ?? '').trim()
  if (!name) return '请输入分类名称'
  const value = name.toLocaleLowerCase()
  const duplicate = props.categories.some((category) => {
    const parentMatches = category.name.toLocaleLowerCase() === value && category.id !== props.editingId
    const childMatches = category.children.some(
      (child) => child.name.toLocaleLowerCase() === value && child.id !== props.editingId,
    )
    return parentMatches || childMatches
  })
  return duplicate ? '分类名称已存在' : true
}

function required(value: unknown) {
  return String(value ?? '').trim() ? true : '此项为必填项'
}
function validateUserId(value: unknown) {
  return (
    props.mode === 'edit' ||
    /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i.test(String(value ?? '').trim()) ||
    '请输入有效的用户 UUID'
  )
}
const submit = handleSubmit(({ name, slug, userId }) => {
  if (!props.busy) emit('save', { name: name.trim(), slug: slug.trim(), userId: userId.trim() })
})
</script>

<template>
  <div
    class="fixed inset-0 z-40 grid place-items-center bg-slate-950/70 p-4 backdrop-blur-sm"
    @click.self="emit('close')"
  >
    <form
      class="w-full max-w-[420px] rounded-lg border border-slate-700 bg-slate-900 p-5 shadow-2xl light:border-slate-200 light:bg-white"
      role="dialog"
      aria-modal="true"
      aria-labelledby="category-form-title"
      @submit="submit"
    >
      <div class="mb-5 flex items-start justify-between gap-3">
        <div>
          <h2 id="category-form-title" class="text-base font-semibold text-slate-100 light:text-slate-800">
            {{ mode === 'edit' ? '编辑分类' : targetParentName ? '新增二级分类' : '新增一级分类' }}
          </h2>
          <p v-if="targetParentName && mode === 'create'" class="mt-1 text-xs text-slate-400">
            上级分类：{{ targetParentName }}
          </p>
        </div>
        <button
          type="button"
          class="grid size-8 place-items-center rounded-md text-slate-400 hover:bg-slate-800 light:hover:bg-slate-100"
          aria-label="关闭"
          @click="emit('close')"
        >
          <X :size="17" />
        </button>
      </div>
      <label class="flex flex-col gap-2 text-xs font-medium text-slate-300 light:text-slate-600">
        分类名称
        <Field v-slot="{ componentField, errorMessage }" name="name" :rules="validateName">
          <Input v-bind="componentField" autofocus maxlength="40" placeholder="输入分类名称" :error="errorMessage" />
        </Field>
      </label>
      <label class="mt-4 flex flex-col gap-2 text-xs text-slate-300">
        分类标识（slug）
        <Field v-slot="{ componentField, errorMessage }" name="slug" :rules="required">
          <Input
            v-bind="componentField"
            maxlength="255"
            placeholder="例如 frontend"
            :error="errorMessage"
            :disabled="busy"
          />
        </Field>
      </label>
      <label v-if="mode === 'create'" class="mt-4 flex flex-col gap-2 text-xs text-slate-300">
        所属用户 ID
        <Field v-slot="{ componentField, errorMessage }" name="userId" :rules="validateUserId">
          <Input v-bind="componentField" placeholder="用户 UUID" :error="errorMessage" :disabled="busy" />
        </Field>
      </label>
      <p v-if="error" role="alert" class="mt-3 text-sm text-rose-400">{{ error }}</p>
      <div class="mt-5 flex justify-end gap-2">
        <Button variant="secondary" size="sm" :disabled="busy" @click="emit('close')">取消</Button>
        <Button type="submit" size="sm" :loading="busy">保存</Button>
      </div>
    </form>
  </div>
</template>
