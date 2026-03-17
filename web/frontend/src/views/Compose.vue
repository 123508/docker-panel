<template>
  <PageHeader title="编排">
    <div class="empty-state">
      <div class="empty-icon">06</div>
      <p class="empty-text">Docker Compose 编排功能即将上线</p>
      <p class="empty-sub">此功能正在开发中，敬请期待。</p>
    </div>
  </PageHeader>
</template>

<script setup lang="ts">
import PageHeader from '@/components/PageHeader.vue';
</script>

<style scoped>
.empty-state {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
}

.empty-icon {
  font-family: var(--font-mono);
  font-size: 48px;
  font-weight: 700;
  color: var(--border-color);
  margin-bottom: 8px;
}

.empty-text {
  font-family: var(--font-mono);
  font-size: 14px;
  color: var(--text-muted);
}

.empty-sub {
  font-family: var(--font-mono);
  font-size: 12px;
  color: var(--text-secondary);
}
</style>
