<template>
  <PageHeader title="镜像" gap="24px">
    <template #actions>
      <div class="header-actions">
        <div class="search-box">
          <span class="search-text">搜索镜像...</span>
        </div>
        <button class="pull-btn">拉取镜像</button>
      </div>
    </template>

    <StatBar :items="stats" />

    <DataTable :bordered="false" header-bg row-gap fixed-row-height>
      <template #head>
        <span class="th" style="width: 40px"></span>
        <span class="th col-repo">仓库</span>
        <span class="th col-tag">标签</span>
        <span class="th col-id">镜像 ID</span>
        <span class="th col-size">大小</span>
        <span class="th col-created">创建时间</span>
        <span class="th" style="flex: 1">操作</span>
      </template>

      <div v-for="img in images" :key="img.repo" class="table-row">
        <div class="td" style="width: 40px; display: flex; align-items: center">
          <div class="row-icon" :style="{ background: img.color }"></div>
        </div>
        <span class="td col-repo">{{ img.repo }}</span>
        <span class="td col-tag td-muted">{{ img.tag }}</span>
        <span class="td col-id td-dim">{{ img.id }}</span>
        <span class="td col-size">{{ img.size }}</span>
        <span class="td col-created td-muted-sm">{{ img.created }}</span>
        <div class="td td-actions" style="flex: 1">
          <button class="link-btn accent">运行</button>
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

const images = [
  { repo: 'nginx', tag: 'latest', id: 'a72860cb95fd', size: '187 MB', created: '2 天前', color: '#3B82F6' },
  { repo: 'postgres', tag: '16-alpine', id: 'b4d181a07f80', size: '432 MB', created: '5 天前', color: '#22C55E' },
  { repo: 'redis', tag: '7-alpine', id: '3c41ce05add9', size: '41 MB', created: '1 周前', color: '#EF4444' },
  { repo: 'node', tag: '20-alpine', id: 'c7b5a7e3f2d1', size: '124 MB', created: '2 周前', color: '#FACC15' }
]

const stats = [
  { label: '总数:', value: '42' },
  { label: '大小:', value: '12.4 GB' },
  { label: '悬空:', value: '3', variant: 'dangling' }
]
</script>

<style scoped>
.header-actions {
  display: flex;
  gap: 12px;
}

.search-box {
  display: flex;
  align-items: center;
  background: var(--bg-card-header);
  height: 40px;
  padding: 0 16px;
  width: 200px;
}

.search-text {
  font-size: 11px;
  color: var(--text-dim);
}

.pull-btn {
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 1px;
  color: var(--accent-text);
  background: var(--accent);
  height: 40px;
  width: 120px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.pull-btn:hover {
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

.col-repo { width: 120px; }
.col-tag { width: 100px; }
.col-id { width: 120px; }
.col-size { width: 80px; }
.col-created { width: 80px; }

.row-icon {
  width: 32px;
  height: 32px;
  flex-shrink: 0;
}

.td-muted {
  color: var(--text-muted);
}

.td-dim {
  color: var(--text-secondary);
  font-size: 11px;
}

.td-muted-sm {
  color: var(--text-muted);
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

.link-btn.accent {
  color: var(--accent);
}

.link-btn.danger {
  color: var(--color-danger);
}
</style>
