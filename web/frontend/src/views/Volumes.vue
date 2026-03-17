<template>
  <PageHeader title="卷" gap="24px">
    <template #actions>
      <div class="header-actions">
        <button class="create-btn">+ 创建</button>
      </div>
    </template>

    <StatBar :items="stats" />

    <DataTable :bordered="false" header-bg row-gap fixed-row-height>
      <template #head>
        <span class="th col-name">名称</span>
        <span class="th col-driver">驱动</span>
        <span class="th col-mount">挂载点</span>
        <span class="th col-size">大小</span>
        <span class="th col-containers">容器</span>
        <span class="th" style="flex: 1">操作</span>
      </template>

      <div v-for="v in volumes" :key="v.name" class="table-row">
        <span class="td col-name">{{ v.name }}</span>
        <span class="td col-driver td-muted">{{ v.driver }}</span>
        <span class="td col-mount td-dim">{{ v.mountpoint }}</span>
        <span class="td col-size">{{ v.size }}</span>
        <div class="td col-containers">
          <CountBadge :count="v.containers" />
        </div>
        <div class="td td-actions" style="flex: 1">
          <button class="link-btn info">查看</button>
          <button class="link-btn danger">移除</button>
        </div>
      </div>
    </DataTable>
  </PageHeader>
</template>

<script setup lang="ts">
import PageHeader from '@/components/PageHeader.vue'
import StatBar from '@/components/StatBar.vue'
import DataTable from '@/components/DataTable.vue'
import CountBadge from '@/components/CountBadge.vue'

const volumes = [
  { name: 'postgres_data', driver: 'local', mountpoint: '/var/lib/docker/volumes/postgres_data/_data', size: '1.2 GB', containers: 1 },
  { name: 'redis_cache', driver: 'local', mountpoint: '/var/lib/docker/volumes/redis_cache/_data', size: '512 MB', containers: 1 },
  { name: 'app_logs', driver: 'local', mountpoint: '/var/lib/docker/volumes/app_logs/_data', size: '256 MB', containers: 0 }
]

const stats = [
  { label: '总数:', value: '8' },
  { label: '已使用:', value: '3.2 GB' },
  { label: '未使用:', value: '2', variant: 'dangling' }
]
</script>

<style scoped>
.create-btn {
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 1px;
  color: var(--accent-text);
  background: var(--accent);
  height: 40px;
  width: 140px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.create-btn:hover {
  background: var(--accent-hover);
}

.th {
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 1px;
  color: var(--text-secondary);
}

.table-row {
  display: flex;
  align-items: center;
  gap: 16px;
  transition: background 0.15s;
}

.table-row:hover {
  background: var(--table-row-hover);
}

.td {
  font-size: 12px;
  color: var(--text-primary);
}

.col-name { width: 140px; }
.col-driver { width: 80px; }
.col-mount { width: 320px; }
.col-size { width: 80px; }
.col-containers { width: 100px; }

.td-muted {
  color: var(--text-muted);
}

.td-dim {
  color: var(--text-secondary);
  font-size: 11px;
}

.td-actions {
  display: flex;
  align-items: center;
  gap: 16px;
}

.link-btn {
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 1px;
  padding: 0;
  background: none;
  border: none;
}

.link-btn.info {
  color: var(--color-info);
}

.link-btn.danger {
  color: var(--color-danger);
}
</style>
