<template>
  <section class="card">
    <h3 v-if="title" class="title">{{ title }}</h3>
    <slot />
  </section>
</template>

<script setup lang="ts">
defineProps<{
  title?: string
}>()
</script>

<style scoped>
.card {
  background: var(--bg-card);
  border: var(--border);
  padding: 16px;
}

.title {
  margin: 0 0 14px;
  font-size: 16px;
  color: var(--text-primary);
}
</style>
