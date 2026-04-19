<template>
  <div class="field">
    <label>
      {{ label }}
      <span v-if="required" class="required">*</span>
    </label>
    <slot />
  </div>
</template>

<script setup lang="ts">
defineProps<{
  label: string
  required?: boolean
}>()
</script>

<style scoped>
.field {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

label {
  color: var(--text-muted);
  font-size: 12px;
}

.required {
  color: #ef4444;
  margin-left: 2px;
}
</style>
