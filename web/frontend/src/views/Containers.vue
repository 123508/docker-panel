<template>
  <dp-page-header title="容器">
    <template #actions>
      <div class="header-actions">
        <dp-search-input v-model="query"/>
        <dp-button text="筛选" size="medium" variant="outlined"/>
        <dp-button text="+ 创建" size="medium" type="primary" class="create-btn"/>
      </div>
    </template>

    <dp-stat-bar :items="stats" bordered />

    <dp-data-table :bordered="true" :columns="columns">

      <div v-for="c in containers" :key="c.name" class="table-row">
        <div class="td" style="width: 60px; display: flex; justify-content: center">
          <div class="row-icon"></div>
        </div>
        <div class="td" style="width: 220px">
          <div class="name-col">
            <span class="name-text">{{ c.name }}</span>
            <span class="name-id">{{ c.id }}</span>
          </div>
        </div>
        <span class="td td-muted" style="width: 180px">{{ c.image }}</span>
        <div class="td" style="width: 120px">
          <dp-status-badge :status="c.running ? 'running' : 'stopped'">
            {{ c.running ? '运行中' : '已停止' }}
          </dp-status-badge>
        </div>
        <span class="td td-muted-sm" style="width: 140px">{{ c.port }}</span>
        <span class="td td-muted-sm" style="width: 120px">{{ c.created }}</span>
        <div class="td td-actions" style="flex: 1">
          <template v-if="c.running">
            <dp-button text="停止" size="small" type="danger" variant="text"/>
            <dp-button text="重启" size="small" type="info" variant="text"/>
            <dp-button text="日志" size="small" variant="text"/>
          </template>
          <template v-else>
            <dp-button text="启动" size="small" type="primary" variant="text"/>
            <dp-button text="移除" size="small" type="danger" variant="text"/>
            <dp-button text="日志" size="small" variant="text"/>
          </template>
        </div>
      </div>
    </dp-data-table>
  </dp-page-header>
</template>

<script setup lang="ts">
import DpPageHeader from '@/components/dp-page-header.vue'
import DpStatBar from '@/components/dp-stat-bar.vue'
import DpDataTable from '@/components/dp-data-table.vue'
import DpStatusBadge from '@/components/dp-status-badge.vue'
import DpButton from "@/components/dp-button.vue";
import DpSearchInput from "@/components/dp-search-input.vue";
import {ref} from "vue";

const columns = [
  { key: 'icon', width: 60 },
  { key: 'name', label: '名称', width: 220 },
  { key: 'image', label: '镜像', width: 180 },
  { key: 'status', label: '状态', width: 120 },
  { key: 'port', label: '端口', width: 140 },
  { key: 'created', label: '创建时间', width: 120 },
  { key: 'actions', label: '操作', flex: 1 }
]

const containers = [
  { name: 'nginx-web-01', id: 'a3f8d92b', image: 'nginx:latest', running: true, port: '80:80, 443:443', created: '2 小时前' },
  { name: 'postgres-db', id: 'e7b2c41a', image: 'postgres:14', running: true, port: '5432:5432', created: '5 天前' },
  { name: 'redis-cache', id: 'f4c9e3d2', image: 'redis:alpine', running: false, port: '6379:6379', created: '1 天前' }
]

const stats = [
  { label: '总数:', value: '24' },
  { label: '运行中:', value: '18', variant: 'running' },
  { label: '已停止:', value: '6', variant: 'stopped' }
]

const query=ref('')

</script>

<style scoped>
.header-actions {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
}

.create-btn {
  width: 120px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.table-row {
  display: flex;
  align-items: center;
  padding: 16px 20px;
  transition: background 0.15s;
}

.table-row:hover {
  background: var(--table-row-hover);
}

.table-row:not(:last-child) {
  border-bottom: var(--border);
}

.td {
  font-size: 14px;
}

.row-icon {
  width: 32px;
  height: 32px;
  background: var(--bg-card-header);
  border: var(--border);
}

.name-col {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.name-text {
  color: var(--text-primary);
  font-size: 14px;
}

.name-id {
  color: var(--text-dim);
  font-size: 11px;
}

.td-muted {
  color: var(--text-muted);
  font-size: 14px;
}

.td-muted-sm {
  color: var(--text-muted);
  font-size: 11px;
}

.td-actions {
  display: flex;
  gap: 12px;
}

.action-btn {
  font-size: 10px;
  font-weight: 700;
  letter-spacing: 1px;
  color: var(--text-primary);
  background: var(--bg-card-header);
  border: var(--border);
  padding: 6px 10px;
}

.action-btn.muted {
  color: var(--text-secondary);
}

.action-btn.accent {
  color: var(--accent-text);
  background: var(--accent);
  border: none;
}

.action-btn.danger {
  color: var(--color-danger);
}
</style>
