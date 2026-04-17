<template>
  <button
      class="button"
      :class="[sizeClass,typeClass,variantClass]"
      v-bind="$attrs"
      @click="click()"
      :disabled="isDisabled"
      :data-disabled="isDisabled"
  >
    <span v-if="loading" class="loading-dot">...</span>
    <span v-else>{{ text }}</span>
  </button>
</template>

<script setup lang="ts">
import { computed } from 'vue'

defineOptions({
  name:'dp-button'
})

type Size = 'small' | 'medium' | 'large'
type ButtonType = 'primary' | 'default' | 'danger' | 'info'
type Variant = 'filled' | 'text' | 'outlined'

interface BaseButtonProps {
  text?: string
  size?: Size
  type?: ButtonType
  variant?: Variant
  disabled?: boolean
  loading?: boolean
}

const props = withDefaults(
    defineProps<BaseButtonProps>(),
    {
      text: 'button',
      size: 'small',
      type: 'default',
      variant: 'filled',
      disabled: false,
      loading: false
    }
)

const emit = defineEmits<{
  (e: 'click'): void
}>()

function click() {
  if (isDisabled.value) return
  emit('click')
}

const sizeClass = computed(() => {
  return `button-${props.size}`
})

const typeClass = computed(() => {
  return `button-${props.type}`
})

const variantClass = computed(() => {
  return `button-${props.variant}`
})

const isDisabled = computed(() => {
  return props.disabled || props.loading
})

</script>

<style scoped>

.button {
  border: none;
  cursor: pointer;
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 1px;
  transition: background-color 0.2s ease;
}

.button-small {
  height: 24px;
}

.button-medium {
  height: 40px;
  padding: 0 16px;
}

.button-large {
  height: 40px;
  padding: 0 20px;
}

/* 文本按钮 */
.button-text {
  background: none;
  padding: 0;
}

/* 实体按钮 */
.button-filled {
  padding: 0 12px;
}

.button-outlined {
  background: var(--bg-card);
  border: var(--border);
  padding: 0 12px;
}


/* 模式 */

/* text 模式 */
.button-text.button-info {
  color: var(--color-info);
}

.button-text.button-danger {
  color: var(--color-danger);
}

.button-text.button-primary {
  color: var(--button-success);
}

/* filled 模式 */
.button-filled.button-primary {
  background: var(--accent);
  color: var(--accent-text);
}

.button-outlined {
  color: var(--text-primary);
}

/* 状态机 */

.button[data-disabled="true"] {
  color: var(--text-secondary);
  cursor: not-allowed;
  background: none;
  opacity: 0.7;
}

/* disabled 统一兜底（最高优先级） */
.button[data-disabled="true"] {
  color: var(--text-secondary);
  cursor: default;
  background: none;
}

/* 动画效果 */

.loading-dot {
  letter-spacing: 2px;
  animation: blink 1s infinite;
}

@keyframes blink {
  0% { opacity: 0.3 }
  50% { opacity: 1 }
  100% { opacity: 0.3 }
}

.button-small.button-text:hover:not([data-disabled="true"]) {
  text-shadow: 0 0 4px currentColor;
}

.button-large.button-filled.button-primary:hover:not([data-disabled="true"]) {
  background: var(--accent-hover);
}

.button-outlined:hover:not([data-disabled="true"]) {
  background: var(--table-row-hover);
}

</style>