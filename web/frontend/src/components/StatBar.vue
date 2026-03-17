<template>
  <div class="stats-bar" :class="{ bordered }">
    <div v-for="(item, i) in items" :key="i" class="stat-item">
      <span class="stat-label">{{ item.label }}</span>
      <span class="stat-value" :class="item.variant">{{ item.value }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  items: Array<{ label: string; value: string; variant?: string }>
  bordered?: boolean
}>()
</script>

<style scoped>
.stats-bar {
  display: flex;
  gap: 24px;
}

.stats-bar.bordered {
  gap: 12px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.stats-bar.bordered .stat-item {
  background: var(--bg-card);
  border: var(--border);
  padding: 8px 12px;
}

.stat-label {
  font-size: 11px;
  letter-spacing: 1px;
  color: var(--text-secondary);
  font-weight: 700;
}

.stat-value {
  font-size: 14px;
  font-weight: 700;
  color: var(--text-primary);
}

.stats-bar.bordered .stat-value {
  font-size: 11px;
}

.stat-value.running { color: var(--color-success); }
.stat-value.stopped { color: var(--color-danger); }
.stat-value.dangling { color: var(--accent); }
.stat-value.builtin { color: var(--text-muted); }
</style>
