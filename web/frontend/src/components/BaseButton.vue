<template>
  <button
      class="button"
      :class="[sizeClass,typeClass,variantClass,$attrs.class]"
      @click="!isDisabled && click()"
      :data-disabled="isDisabled"
  >
    <span v-if="loading" class="loading-dot">...</span>
    <span v-else>{{ text }}</span>
  </button>
</template>

<script setup lang="ts">
import { computed } from 'vue'

type Size = 'small' | 'medium' | 'large'
type ButtonType = 'primary' | 'default' | 'danger' | 'info'
type Variant = 'solid' | 'text'

const props = withDefaults(
    defineProps<{
      text?: string
      size?: Size
      type?: ButtonType
      variant?: Variant
      disabled?: boolean
      loading?: boolean
    }>(),
    {
      text: 'button',
      size: 'small',
      type: 'primary',
      variant: 'solid',
      disabled: false,
      loading: false
    }
)

const emit = defineEmits<{
  (e: 'click'): void
}>()

function click() {
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
.button-solid {
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

/* solid 模式 */
.button-solid.button-primary {
  background: var(--accent);
  color: var(--accent-text);
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

.button-large.button-solid.button-primary:hover:not([data-disabled="true"]) {
  background: var(--accent-hover);
}

</style>