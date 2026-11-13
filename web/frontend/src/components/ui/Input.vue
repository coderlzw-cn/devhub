<script setup lang="ts">
import { computed, ref, useAttrs, useSlots } from 'vue'
import { X } from 'lucide-vue-next'

defineOptions({ name: 'UiInput', inheritAttrs: false })

interface Props {
  modelValue?: string | number
  type?: 'text' | 'email' | 'url' | 'search' | 'password' | 'number' | 'tel'
  appearance?: 'outlined' | 'plain'
  size?: 'sm' | 'md' | 'lg'
  disabled?: boolean
  readonly?: boolean
  clearable?: boolean
  error?: string | boolean
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: '',
  type: 'text',
  appearance: 'outlined',
  size: 'md',
  disabled: false,
  readonly: false,
  clearable: false,
  error: false,
})

const emit = defineEmits<{
  (event: 'update:modelValue', value: string): void
  (event: 'clear'): void
}>()

const attrs = useAttrs()
const slots = useSlots()
const inputRef = ref<HTMLInputElement | null>(null)
const valueText = computed(() => String(props.modelValue ?? ''))
const hasError = computed(() => Boolean(props.error))
const errorId = computed(() => (typeof attrs.id === 'string' ? `${attrs.id}-error` : undefined))
const describedBy = computed(() => {
  return [attrs['aria-describedby'], hasError.value ? errorId.value : undefined].filter(Boolean).join(' ') || undefined
})

const fieldClass = computed(() => [
  'flex w-full min-w-0 items-center gap-2.5 text-slate-400 transition focus-within:border-violet-400 focus-within:ring-2 focus-within:ring-violet-400/20 light:text-slate-500',
  props.appearance === 'outlined'
    ? 'rounded-md border bg-slate-900 px-3 light:bg-white'
    : 'rounded-[9px] border border-slate-700 bg-[#151f30] px-[13px] light:border-slate-200 light:bg-slate-50',
  props.appearance === 'outlined' && !hasError.value ? 'border-slate-600 light:border-slate-200' : '',
  hasError.value ? '!border-rose-500 focus-within:!border-rose-400 focus-within:!ring-rose-400/20' : '',
  props.size === 'sm' ? 'h-8' : props.size === 'lg' ? 'h-11' : 'h-9',
  props.disabled ? 'cursor-not-allowed opacity-50' : '',
])

const controlClass = computed(() => [
  'min-w-0 flex-1 bg-transparent text-xs text-slate-100 outline-none placeholder:text-slate-500 disabled:cursor-not-allowed light:text-slate-800',
  props.appearance === 'outlined' ? 'py-2.5' : 'h-full py-0',
  props.size === 'sm' ? 'text-[11px]' : props.size === 'lg' ? 'text-sm' : '',
])

function updateValue(event: Event) {
  emit('update:modelValue', (event.target as HTMLInputElement).value)
}

function clearValue() {
  emit('update:modelValue', '')
  emit('clear')
  inputRef.value?.focus()
}
</script>

<template>
  <div class="min-w-0" :class="attrs.class">
    <div :class="fieldClass">
      <span
        v-if="slots.leading"
        class="flex shrink-0 items-center"
        :class="appearance === 'outlined' ? '' : 'text-slate-400'"
      >
        <slot name="leading" />
      </span>
      <input
        ref="inputRef"
        v-bind="{ ...attrs, class: undefined, 'aria-describedby': describedBy, 'aria-invalid': hasError || undefined }"
        :value="modelValue"
        :type="type"
        :disabled="disabled"
        :readonly="readonly"
        :class="controlClass"
        @input="updateValue"
      />
      <button
        v-if="clearable && valueText && !disabled && !readonly"
        type="button"
        class="grid size-6 shrink-0 place-items-center rounded-sm text-slate-400 transition hover:bg-slate-700 hover:text-white focus-visible:ring-2 focus-visible:ring-violet-400 focus-visible:outline-none light:hover:bg-slate-100 light:hover:text-slate-800"
        aria-label="清空输入"
        @click="clearValue"
      >
        <X :size="14" />
      </button>
      <span v-if="slots.trailing" class="flex shrink-0 items-center">
        <slot name="trailing" />
      </span>
    </div>
    <p v-if="hasError" :id="errorId" class="mt-1 text-[10px] text-rose-400" role="alert">
      {{ typeof error === 'string' ? error : '输入内容有误' }}
    </p>
  </div>
</template>
