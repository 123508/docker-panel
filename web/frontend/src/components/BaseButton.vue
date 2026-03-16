<template>
  <button
      class="button"
      :class="[
        sizeClass,
        typeClass]"
      @click="click">
    {{ text }}
  </button>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { PropType } from 'vue'

type Size = 'small' | 'medium' | 'large'
type ButtonType = 'primary' | 'success' | 'danger' | 'warning'

const props = defineProps({
  text: {
    type: String,
    default: 'button'
  },
  size:{
    type:String as PropType<Size>,
    default:'medium',
  },
  type:{
    type: String as PropType<ButtonType>,
    default:'primary'
  }
})

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

</script>

<style scoped>

.button {
  border: none;
  border-radius: 6px;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.2s ease;
}

.button-small {
  height: 28px;
  padding: 0 10px;
  font-size: 12px;
}

.button-medium {
  height: 36px;
  padding: 0 14px;
  font-size: 14px;
}

.button-large {
  height: 44px;
  padding: 0 18px;
  font-size: 16px;
}

.button-primary {
  background: var(--button-primary);
  color: var(--button-primary-text);
}

.button-primary:hover {
  background: var(--button-primary-hover);
}


.button-success {
  background: var(--button-success);
  color: var(--button-success-text);
}

.button-success:hover {
  background: var(--button-success-hover);
}


.button-danger {
  background: var(--button-danger);
  color: var(--button-danger-text);
}

.button-danger:hover {
  background: var(--button-danger-hover);
}


.button-warning {
  background: var(--button-warning);
  color: var(--button-warning-text);
}

.button-warning:hover {
  background: var(--button-warning-hover);
}

</style>