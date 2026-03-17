<template>
  <span class="status-badge" :class="status">
    <span class="status-dot"></span>
    <slot />
  </span>
</template>

<script setup lang="ts">
  defineProps<{
    status: 'running' | 'stopped'
  }>()
</script>

<style scoped>
.status-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: 10px;
  font-weight: 700;
  letter-spacing: 1px;
  padding: 3px 6px;
}

.status-badge.running {
  color: var(--color-success);
}

.status-badge.stopped {
  color: var(--color-danger);
}

.status-dot {
  width: 6px;
  height: 6px;
  border-radius: 0;
  background: currentColor;
}
</style>
