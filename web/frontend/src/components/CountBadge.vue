<template>
  <span class="container-badge" :class="count > 0 ? 'active' : 'empty'">
    {{ count }}
  </span>
</template>

<script setup lang="ts">
  defineProps<{
    count: number
  }>()
</script>

<style scoped>
.container-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 20px;
  font-size: 11px;
  font-weight: 700;
}

.container-badge.active {
  background: rgba(34, 197, 94, 0.125);
  color: var(--color-success);
}

.container-badge.empty {
  background: rgba(113, 113, 122, 0.125);
  color: var(--text-secondary);
}
</style>
