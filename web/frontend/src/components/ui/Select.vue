<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, useId, useSlots } from 'vue'
import { Check, ChevronDown } from 'lucide-vue-next'

defineOptions({ name: 'UiSelect', inheritAttrs: false })

type SelectValue = string | number
type SelectOption = SelectValue | { label: string; value: SelectValue; disabled?: boolean }

interface Props {
  modelValue?: SelectValue | null
  options: readonly SelectOption[]
  placeholder?: string
  disabled?: boolean
  size?: 'sm' | 'md'
  block?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: null,
  placeholder: '请选择',
  disabled: false,
  size: 'md',
  block: false,
})

const emit = defineEmits<{ (event: 'update:modelValue', value: SelectValue): void }>()
const slots = useSlots()
const rootRef = ref<HTMLElement | null>(null)
const triggerRef = ref<HTMLButtonElement | null>(null)
const isOpen = ref(false)
const activeIndex = ref(-1)
const listboxId = `base-select-${useId()}`

const normalizedOptions = computed(() =>
  props.options.map((option) =>
    typeof option === 'string' || typeof option === 'number' ? { label: String(option), value: option } : option,
  ),
)
const selectedIndex = computed(() =>
  normalizedOptions.value.findIndex((option) => Object.is(option.value, props.modelValue)),
)
const selectedOption = computed(() =>
  selectedIndex.value >= 0 ? normalizedOptions.value[selectedIndex.value] : undefined,
)
const hasLeading = computed(() => Boolean(slots.leading))

const triggerClass = computed(() => [
  'flex w-full items-center justify-between gap-2 rounded-lg border border-slate-700 bg-slate-900 px-2.5 text-left text-xs text-slate-200 outline-none transition focus:border-violet-500 focus:ring-2 focus:ring-violet-400/20 disabled:cursor-not-allowed disabled:opacity-50 light:border-slate-200 light:bg-white light:text-slate-700',
  props.size === 'sm' ? 'h-8 text-[11px]' : 'h-9',
  hasLeading.value ? 'pl-9' : '',
])

function openMenu() {
  if (props.disabled) return
  isOpen.value = true
  activeIndex.value = selectedIndex.value >= 0 ? selectedIndex.value : firstEnabledIndex(0, 1)
}

function closeMenu() {
  isOpen.value = false
  activeIndex.value = -1
}

function firstEnabledIndex(start: number, direction: 1 | -1) {
  const options = normalizedOptions.value
  for (let index = start; index >= 0 && index < options.length; index += direction) {
    if (!options[index].disabled) return index
  }
  return -1
}

function selectOption(index: number) {
  const option = normalizedOptions.value[index]
  if (!option || option.disabled) return
  emit('update:modelValue', option.value)
  closeMenu()
  triggerRef.value?.focus()
}

function handleKeydown(event: KeyboardEvent) {
  if (props.disabled) return
  if (event.key === 'Escape') {
    if (isOpen.value) {
      event.preventDefault()
      closeMenu()
    }
    return
  }
  if (event.key === 'Enter' || event.key === ' ') {
    event.preventDefault()
    if (!isOpen.value) openMenu()
    else if (activeIndex.value >= 0) selectOption(activeIndex.value)
    return
  }
  if (event.key !== 'ArrowDown' && event.key !== 'ArrowUp' && event.key !== 'Home' && event.key !== 'End') return
  event.preventDefault()
  if (!isOpen.value) {
    openMenu()
    return
  }
  const direction = event.key === 'ArrowUp' ? -1 : 1
  const start =
    event.key === 'Home' ? 0 : event.key === 'End' ? normalizedOptions.value.length - 1 : activeIndex.value + direction
  const nextIndex = firstEnabledIndex(start, direction)
  if (nextIndex >= 0) activeIndex.value = nextIndex
}

function handleDocumentClick(event: MouseEvent) {
  if (rootRef.value && !rootRef.value.contains(event.target as Node)) closeMenu()
}

onMounted(() => document.addEventListener('mousedown', handleDocumentClick))
onBeforeUnmount(() => document.removeEventListener('mousedown', handleDocumentClick))
</script>

<template>
  <div ref="rootRef" :class="['relative', block ? 'w-full' : '']">
    <button
      ref="triggerRef"
      v-bind="$attrs"
      type="button"
      :class="triggerClass"
      :disabled="disabled"
      role="combobox"
      :aria-expanded="isOpen"
      :aria-controls="listboxId"
      :aria-haspopup="'listbox'"
      @click="isOpen ? closeMenu() : openMenu()"
      @keydown="handleKeydown"
    >
      <span v-if="hasLeading" class="pointer-events-none absolute top-1/2 left-3 -translate-y-1/2 text-slate-400"
        ><slot name="leading"
      /></span>
      <span class="min-w-0 flex-1 truncate" :class="selectedOption ? '' : 'text-slate-500'">{{
        selectedOption?.label ?? placeholder
      }}</span>
      <ChevronDown
        :size="16"
        class="shrink-0 text-slate-400 transition"
        :class="isOpen ? 'rotate-180' : ''"
        aria-hidden="true"
      />
    </button>

    <ul
      v-if="isOpen"
      :id="listboxId"
      role="listbox"
      class="absolute top-full left-0 z-30 mt-1 max-h-60 w-full min-w-full overflow-auto rounded-lg border border-slate-700 bg-slate-900 p-1 shadow-xl light:border-slate-200 light:bg-white"
    >
      <li
        v-for="(option, index) in normalizedOptions"
        :id="`${listboxId}-option-${index}`"
        :key="`${String(option.value)}-${index}`"
        role="option"
        :aria-selected="selectedIndex === index"
        :class="[
          'flex min-h-8 cursor-pointer items-center justify-between rounded-md px-2.5 text-xs transition',
          option.disabled
            ? 'cursor-not-allowed opacity-40'
            : index === activeIndex || index === selectedIndex
              ? 'bg-violet-600/20 text-violet-200 light:bg-violet-50 light:text-violet-700'
              : 'text-slate-300 hover:bg-slate-800 light:text-slate-600 light:hover:bg-slate-100',
        ]"
        @mousedown.prevent
        @click="selectOption(index)"
      >
        <span class="truncate">{{ option.label }}</span>
        <Check v-if="selectedIndex === index" :size="14" aria-hidden="true" />
      </li>
    </ul>
  </div>
</template>
