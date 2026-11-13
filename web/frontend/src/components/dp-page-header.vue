<template>
  <div class="page" :style="{ gap: gap }">
    <div class="page-header">
      <div><div class="page-eyebrow">WORKSPACE <span>/</span> DOCKER PANEL</div><h1 class="page-title">{{ title }}</h1></div>
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
  min-height: 100vh;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
  flex-wrap: wrap;
}

.page-title {
  font-family: var(--font-display);
  font-size: 30px;
  font-weight: 700;
  letter-spacing: -0.5px;
  color: var(--text-primary);
}
</style>

<style>
.page-eyebrow { color: var(--text-secondary); font: 700 10px var(--font-mono); letter-spacing: 1.6px; margin-bottom: 10px; }
.page-eyebrow span { color: var(--accent); margin: 0 5px; }
@media (max-width: 700px) { .page { padding: 24px 18px; } .page-title { font-size: 26px; } }
</style>
