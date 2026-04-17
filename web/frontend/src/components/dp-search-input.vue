<template>
  <div class="search-box">
    <!-- 输入框 -->
    <input
        class="input"
        :value="modelValue"
        :placeholder="placeholder"
        @input="onInput"
        @keydown.enter="onSearch"
    />

    <!-- 清空 -->
    <span
        v-if="clearable && modelValue"
        class="clear"
        @click="onClear"
    >
      ✕
    </span>
  </div>
</template>

<script setup lang="ts">

defineOptions({
  name:'dp-search-input'
})

const props = withDefaults(
    defineProps<{
      modelValue: string
      placeholder?: string
      clearable?: boolean
    }>(),
    {
      placeholder: '搜索...',
      clearable: true
    }
)

const emit = defineEmits<{
  (e: 'update:modelValue', v: string): void
  (e: 'search', v: string): void
  (e: 'clear'): void
}>()

function onInput(e: Event) {
  const value = (e.target as HTMLInputElement).value
  emit('update:modelValue', value)
}

function onSearch() {
  emit('search', props.modelValue)
}

function onClear() {
  emit('update:modelValue', '')
  emit('clear')
}
</script>

<style scoped>
/* ✅ 完全复用你原来的 search-box 设计 */
.search-box {
  display: flex;
  align-items: center;

  background: var(--bg-card);
  border: var(--border);

  padding: 10px 14px;
  width: 240px;

  gap: 8px;
}

/* input */
.input {
  flex: 1;
  border: none;
  outline: none;
  background: transparent;

  font-size: 11px;
  color: var(--text-primary);
}

/* placeholder 风格保持一致 */
.input::placeholder {
  color: var(--text-dim);
}

/* 清空按钮 */
.clear {
  font-size: 11px;
  color: var(--text-dim);
  cursor: pointer;
}

.clear:hover {
  color: var(--text-primary);
}

/* focus 状态（增强但不破坏风格） */
.search-box:focus-within {
  border-color: var(--accent);
}
</style>