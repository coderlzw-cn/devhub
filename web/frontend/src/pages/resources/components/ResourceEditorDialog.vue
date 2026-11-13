<script setup lang="ts">
import { Field, useForm } from 'vee-validate'
import { X } from 'lucide-vue-next'
import type { Resource } from '@/data/resources'
import Button from '@/components/ui/Button.vue'
import Input from '@/components/ui/Input.vue'
import Select from '@/components/ui/Select.vue'

export interface ResourceFormValues {
  name: string
  description: string
  category: string
  url: string
  userId: string
  logoUrl: string
  sortOrder: number
  isPinned: boolean
}

const props = defineProps<{
  mode: 'create' | 'edit'
  resource?: Resource | null
  categoryOptions: { label: string; value: string }[]
  loading?: boolean
  error?: string
  initialCategory?: string
}>()
const emit = defineEmits<{
  close: []
  submit: [values: ResourceFormValues]
}>()
const { handleSubmit } = useForm<ResourceFormValues>({
  initialValues: {
    name: props.resource?.name ?? '',
    description: props.resource?.description ?? '',
    category: props.resource?.categoryId ?? props.initialCategory ?? '',
    url: props.resource?.url ?? '',
    userId: props.resource?.userId ?? '',
    logoUrl: props.resource?.logoUrl ?? '',
    sortOrder: props.resource?.sortOrder ?? 0,
    isPinned: props.resource?.isPinned ?? false,
  },
})
function validateRequired(value: unknown) {
  return String(value ?? '').trim() ? true : '此项为必填项'
}

function validateUrl(value: unknown) {
  try {
    const url = new URL(String(value ?? ''))
    return ['http:', 'https:'].includes(url.protocol) || '请输入有效的网址'
  } catch {
    return '请输入有效的网址'
  }
}

function validateUserId(value: unknown) {
  return (
    props.mode === 'edit' ||
    /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i.test(String(value ?? '').trim()) ||
    '请输入有效的用户 UUID'
  )
}
function validateSortOrder(value: unknown) {
  return (
    (String(value ?? '').trim() !== '' &&
      Number.isSafeInteger(Number(value)) &&
      Number(value) >= -2147483648 &&
      Number(value) <= 2147483647) ||
    '请输入整数'
  )
}
const submit = handleSubmit((values) => {
  if (!props.loading) emit('submit', { ...values, sortOrder: Number(values.sortOrder) })
})
</script>

<template>
  <div
    class="fixed inset-0 z-40 grid place-items-center bg-slate-950/70 p-4 backdrop-blur-sm"
    @click.self="emit('close')"
  >
    <form
      class="max-h-[90vh] w-full max-w-115 overflow-auto rounded-lg border border-slate-700 bg-slate-900 p-5 shadow-2xl light:border-slate-200 light:bg-white"
      role="dialog"
      aria-modal="true"
      aria-labelledby="resource-form-title"
      @submit="submit"
    >
      <div class="mb-5 flex items-start justify-between gap-3">
        <div>
          <h2 id="resource-form-title" class="text-base font-semibold text-slate-100 light:text-slate-800">
            {{ mode === 'edit' ? '编辑资源' : '添加资源' }}
          </h2>
          <p class="mt-1 text-xs text-slate-400">填写资源信息并选择所属分类</p>
        </div>
        <button
          type="button"
          class="grid size-8 place-items-center rounded-md text-slate-400 hover:bg-slate-800 light:hover:bg-slate-100"
          aria-label="关闭"
          :disabled="loading"
          @click="emit('close')"
        >
          <X :size="17" />
        </button>
      </div>
      <div class="grid gap-3">
        <label class="flex flex-col gap-1.5 text-xs text-slate-300 light:text-slate-600">
          名称
          <Field v-slot="{ componentField, errorMessage }" name="name" :rules="validateRequired">
            <Input
              v-bind="componentField"
              maxlength="255"
              placeholder="例如：Vue"
              :error="errorMessage"
              :disabled="loading"
            />
          </Field>
        </label>
        <label class="flex flex-col gap-1.5 text-xs text-slate-300 light:text-slate-600">
          网址
          <Field v-slot="{ componentField, errorMessage }" name="url" :rules="validateUrl">
            <Input
              v-bind="componentField"
              type="url"
              placeholder="https://example.com"
              :error="errorMessage"
              :disabled="loading"
            />
          </Field>
        </label>
        <label class="flex flex-col gap-1.5 text-xs text-slate-300 light:text-slate-600">
          分类
          <Field
            v-slot="{ componentField, errorMessage }"
            name="category"
            :rules="resource?.categoryId ? validateRequired : undefined"
          >
            <Select
              v-bind="componentField"
              :options="resource?.categoryId ? categoryOptions : [{ label: '未分类', value: '' }, ...categoryOptions]"
              placeholder="选择分类"
              block
              :disabled="loading"
            />
            <span v-if="errorMessage" class="text-[10px] text-rose-400" role="alert">{{ errorMessage }}</span>
          </Field>
        </label>
        <label class="flex flex-col gap-1.5 text-xs text-slate-300 light:text-slate-600">
          描述
          <Field v-slot="{ componentField, errorMessage }" name="description">
            <Input
              v-bind="componentField"
              maxlength="4096"
              placeholder="简要描述这个资源"
              :error="errorMessage"
              :disabled="loading"
            />
          </Field>
        </label>
        <label v-if="mode === 'create'" class="flex flex-col gap-1.5 text-xs text-slate-300">
          所属用户 ID
          <Field v-slot="{ componentField, errorMessage }" name="userId" :rules="validateUserId">
            <Input v-bind="componentField" placeholder="用户 UUID" :error="errorMessage" :disabled="loading" />
          </Field>
        </label>
        <label class="flex flex-col gap-1.5 text-xs text-slate-300">
          图标地址
          <Field v-slot="{ componentField }" name="logoUrl">
            <Input v-bind="componentField" placeholder="https://example.com/logo.png" :disabled="loading" />
          </Field>
        </label>
        <label class="flex flex-col gap-1.5 text-xs text-slate-300">
          排序
          <Field v-slot="{ componentField, errorMessage }" name="sortOrder" :rules="validateSortOrder">
            <Input v-bind="componentField" type="number" step="1" :error="errorMessage" :disabled="loading" />
          </Field>
        </label>
        <label class="flex items-center gap-2 text-xs text-slate-300">
          <Field v-slot="{ field }" name="isPinned" type="checkbox" :value="true">
            <input v-bind="field" type="checkbox" :disabled="loading" /> </Field
          >置顶站点
        </label>
      </div>
      <p v-if="error" class="mt-3 text-xs text-rose-400" role="alert">{{ error }}</p>
      <div class="mt-5 flex justify-end gap-2">
        <Button variant="secondary" :disabled="loading" @click="emit('close')">取消</Button>
        <Button type="submit" size="sm" :loading="loading">{{ mode === 'edit' ? '保存修改' : '保存' }}</Button>
      </div>
    </form>
  </div>
</template>
