<script setup lang="ts">
import { computed } from 'vue'
import { LoaderCircle } from 'lucide-vue-next'

defineOptions({ name: 'UiButton' })

type ButtonVariant = 'primary' | 'secondary' | 'ghost' | 'danger'
type ButtonSize = 'sm' | 'md' | 'lg'

interface Props {
  variant?: ButtonVariant
  size?: ButtonSize
  type?: 'button' | 'submit' | 'reset'
  loading?: boolean
  disabled?: boolean
  block?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'primary',
  size: 'md',
  type: 'button',
  loading: false,
  disabled: false,
  block: false,
})

const variantClasses: Record<ButtonVariant, string> = {
  primary:
    'border border-violet-300/30 bg-gradient-to-br from-violet-500 to-indigo-700 text-white shadow-[0_6px_20px_#6552e844] hover:-translate-y-px hover:brightness-110',
  secondary:
    'border border-slate-600 bg-slate-800 text-slate-300 hover:border-violet-500 hover:text-white light:border-slate-300 light:bg-white light:text-slate-700 light:hover:border-violet-400 light:hover:text-violet-700',
  ghost:
    'text-slate-400 hover:bg-slate-800 hover:text-white light:text-slate-600 light:hover:bg-slate-100 light:hover:text-slate-900',
  danger: 'border border-rose-400/30 bg-rose-600 text-white hover:bg-rose-500',
}

const sizeClasses: Record<ButtonSize, string> = {
  sm: 'h-8 gap-1.5 rounded-md px-3 text-[10px]',
  md: 'h-9 gap-2 rounded-lg px-3.5 text-xs',
  lg: 'h-11 gap-2 rounded-[10px] px-4 text-sm',
}

const buttonClass = computed(() => [
  'inline-flex items-center justify-center whitespace-nowrap font-semibold outline-none transition focus-visible:ring-2 focus-visible:ring-violet-400/60 focus-visible:ring-offset-2 focus-visible:ring-offset-slate-950 disabled:cursor-not-allowed disabled:opacity-50 disabled:hover:translate-y-0 disabled:hover:brightness-100',
  variantClasses[props.variant],
  sizeClasses[props.size],
  props.block ? 'w-full' : '',
])
</script>

<template>
  <button :type="type" :disabled="disabled || loading" :aria-busy="loading || undefined" :class="buttonClass">
    <LoaderCircle v-if="loading" class="size-4 animate-spin" aria-hidden="true" />
    <slot />
  </button>
</template>
