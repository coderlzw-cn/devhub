<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'

defineOptions({ name: 'UiDropdown' })

export interface DropdownItem {
  label: string
  value: string
  icon?: object
  disabled?: boolean
  separatorBefore?: boolean
}

const props = withDefaults(
  defineProps<{
    items: readonly DropdownItem[]
    align?: 'start' | 'end'
    placement?: 'top' | 'bottom'
    fit?: boolean
    triggerLabel?: string
  }>(),
  {
    align: 'start',
    placement: 'top',
    fit: false,
    triggerLabel: '打开菜单',
  },
)

const emit = defineEmits<{ (event: 'select', value: string): void }>()
const rootRef = ref<HTMLElement | null>(null)
const triggerRef = ref<HTMLButtonElement | null>(null)
const isOpen = ref(false)
const activeIndex = ref(-1)
const menuId = `base-dropdown-${Math.random().toString(36).slice(2)}`

const menuClass = computed(() => (props.align === 'end' ? 'right-0' : 'left-0'))
const placementClass = computed(() => (props.placement === 'bottom' ? 'top-full mt-2' : 'bottom-full mb-2'))

function enabledIndex(start: number, direction: 1 | -1) {
  for (let index = start; index >= 0 && index < props.items.length; index += direction) {
    if (!props.items[index].disabled) return index
  }
  return -1
}

function openMenu(index = enabledIndex(0, 1)) {
  isOpen.value = true
  activeIndex.value = index
}

function closeMenu(returnFocus = false) {
  isOpen.value = false
  activeIndex.value = -1
  if (returnFocus) triggerRef.value?.focus()
}

function selectItem(index: number) {
  const item = props.items[index]
  if (!item || item.disabled) return
  emit('select', item.value)
  closeMenu(true)
}

function handleKeydown(event: KeyboardEvent) {
  if (event.key === 'Escape') {
    if (isOpen.value) {
      event.preventDefault()
      closeMenu(true)
    }
    return
  }
  if (!isOpen.value && (event.key === 'ArrowDown' || event.key === 'Enter' || event.key === ' ')) {
    event.preventDefault()
    openMenu()
    return
  }
  if (!isOpen.value) return
  if (event.key === 'ArrowDown' || event.key === 'ArrowUp') {
    event.preventDefault()
    const direction = event.key === 'ArrowDown' ? 1 : -1
    const next = enabledIndex(activeIndex.value + direction, direction)
    if (next >= 0) activeIndex.value = next
  } else if (event.key === 'Home' || event.key === 'End') {
    event.preventDefault()
    activeIndex.value = enabledIndex(event.key === 'Home' ? 0 : props.items.length - 1, event.key === 'Home' ? 1 : -1)
  } else if (event.key === 'Enter') {
    event.preventDefault()
    selectItem(activeIndex.value)
  }
}

function handleDocumentClick(event: MouseEvent) {
  if (rootRef.value && !rootRef.value.contains(event.target as Node)) closeMenu()
}

onMounted(() => document.addEventListener('mousedown', handleDocumentClick))
onBeforeUnmount(() => document.removeEventListener('mousedown', handleDocumentClick))
</script>

<template>
  <div ref="rootRef" :class="['relative', fit ? 'w-fit' : 'w-full']">
    <button
      ref="triggerRef"
      type="button"
      class="w-full text-left outline-none focus-visible:ring-2 focus-visible:ring-violet-400/60"
      aria-haspopup="menu"
      :aria-label="triggerLabel"
      :aria-expanded="isOpen"
      :aria-controls="menuId"
      @click="isOpen ? closeMenu() : openMenu()"
      @keydown="handleKeydown"
    >
      <slot :open="isOpen" />
    </button>
    <div
      v-if="isOpen"
      :id="menuId"
      role="menu"
      :class="[
        'absolute z-30 min-w-47.5 rounded-lg border border-slate-700 bg-slate-900 p-1.5 shadow-xl light:border-slate-200 light:bg-white',
        placementClass,
        menuClass,
      ]"
    >
      <template v-for="(item, index) in items" :key="item.value">
        <div v-if="item.separatorBefore" class="my-1.5 h-px bg-slate-700 light:bg-slate-200" role="separator" />
        <button
          type="button"
          role="menuitem"
          :disabled="item.disabled"
          :tabindex="index === activeIndex ? 0 : -1"
          :class="[
            'flex h-9 w-full items-center gap-2.5 rounded-md px-2.5 text-xs transition disabled:cursor-not-allowed disabled:opacity-40',
            index === activeIndex
              ? 'bg-slate-800 text-white light:bg-slate-100 light:text-slate-900'
              : 'text-slate-300 hover:bg-slate-800 hover:text-white light:text-slate-700 light:hover:bg-slate-100',
          ]"
          @mouseenter="!item.disabled && (activeIndex = index)"
          @click="selectItem(index)"
        >
          <component :is="item.icon" v-if="item.icon" :size="15" class="shrink-0 text-slate-400" aria-hidden="true" />
          <span>{{ item.label }}</span>
        </button>
      </template>
    </div>
  </div>
</template>
