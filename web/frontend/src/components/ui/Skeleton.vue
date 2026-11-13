<script setup lang="ts">
defineOptions({ name: 'UiSkeleton' })

withDefaults(
  defineProps<{
    variant?: 'text' | 'circle' | 'rect'
    animated?: boolean
  }>(),
  {
    variant: 'rect',
    animated: true,
  },
)
</script>

<template>
  <div
    aria-hidden="true"
    :class="[
      'bg-slate-700/70 light:bg-slate-200/80',
      animated ? 'animate-pulse' : '',
      variant === 'text' ? 'h-3 rounded-sm' : '',
      variant === 'circle' ? 'aspect-square rounded-full' : '',
      variant === 'rect' ? 'rounded-md' : '',
    ]"
  />
</template>
