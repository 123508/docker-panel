<template>
  <div class="page" :style="{ gap: gap }">
    <div class="page-header">
      <h1 class="page-title">{{ title }}</h1>
      <slot name="actions"></slot>
    </div>
    <slot></slot>
  </div>
</template>

<script setup lang="ts">
  withDefaults(defineProps<{
    title: string
    gap?: string
  }>(), {
    gap: 'var(--section-gap)'
  })
</script>

<style>
/* Global page layout styles */
.page {
  padding: var(--page-padding-y) var(--page-padding-x);
  display: flex;
  flex-direction: column;
  height: 100vh;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.page-title {
  font-family: var(--font-display);
  font-size: 40px;
  font-weight: 700;
  letter-spacing: -1px;
  color: var(--text-primary);
}
</style>
