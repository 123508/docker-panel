<template>
  <div class="dp-select" :style="containerStyle">
    <select
      class="dp-select-native"
      :value="modelValue"
      :disabled="disabled"
      v-bind="$attrs"
      @change="handleChange"
    >
      <option v-if="placeholder" value="" disabled hidden>{{ placeholder }}</option>
      <option
        v-for="item in normalizedOptions"
        :key="String(item.value)"
        :value="item.value"
        :disabled="item.disabled"
      >
        {{ item.label }}
      </option>
    </select>
    <span class="dp-select-arrow" aria-hidden="true"></span>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

type SelectOption = string | number | { label: string; value: string | number; disabled?: boolean }

const props = withDefaults(
  defineProps<{
    modelValue?: string | number
    options?: SelectOption[]
    placeholder?: string
    width?: string
    height?: string
    disabled?: boolean
  }>(),
  {
    modelValue: '',
    options: () => [],
    placeholder: '',
    width: '100%',
    height: '40px',
    disabled: false
  }
)

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
  (e: 'change', value: string): void
}>()

const normalizedOptions = computed(() =>
  props.options.map((item) => {
    if (typeof item === 'string' || typeof item === 'number') {
      return { label: String(item), value: String(item), disabled: false }
    }
    return {
      label: item.label,
      value: String(item.value),
      disabled: !!item.disabled
    }
  })
)

const containerStyle = computed(() => ({
  width: props.width,
  height: props.height
}))

const handleChange = (event: Event) => {
  const target = event.target as HTMLSelectElement
  emit('update:modelValue', target.value)
  emit('change', target.value)
}
</script>

<style scoped>
.dp-select {
  position: relative;
  display: inline-flex;
  align-items: center;
}

.dp-select-native {
  width: 100%;
  height: 100%;
  border: var(--border);
  background: var(--select-bg);
  color: var(--text-primary);
  font-family: var(--font-mono);
  font-size: 13px;
  padding: 0 36px 0 14px;
  appearance: none;
  -webkit-appearance: none;
  -moz-appearance: none;
  outline: none;
}

.dp-select-native:hover:not(:disabled) {
  background: var(--select-hover-bg);
}

.dp-select-native:focus {
  border-color: var(--accent);
  box-shadow: inset 0 0 0 1px var(--accent);
}

.dp-select-native:disabled {
  color: var(--text-secondary);
  cursor: not-allowed;
  opacity: 0.8;
}

.dp-select-native option {
  background: var(--select-option-bg);
  color: var(--select-option-text);
}

.dp-select-arrow {
  position: absolute;
  right: 14px;
  width: 8px;
  height: 8px;
  border-right: 2px solid var(--select-arrow);
  border-bottom: 2px solid var(--select-arrow);
  transform: rotate(45deg) translateY(-2px);
  pointer-events: none;
}
</style>
