<template>
  <PageHeader title="网络" gap="24px">
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
        <span class="th col-scope">范围</span>
        <span class="th col-subnet">子网</span>
        <span class="th col-gateway">网关</span>
        <span class="th col-containers">容器</span>
        <span class="th" style="flex: 1">操作</span>
      </template>

      <div v-for="n in networks" :key="n.name" class="table-row">
        <span class="td col-name">{{ n.name }}</span>
        <span class="td col-driver td-muted">{{ n.driver }}</span>
        <span class="td col-scope td-muted">{{ n.scope }}</span>
        <span class="td col-subnet td-dim">{{ n.subnet }}</span>
        <span class="td col-gateway td-dim">{{ n.gateway }}</span>
        <div class="td col-containers">
          <CountBadge :count="n.containers" />
        </div>
        <div class="td td-actions" style="flex: 1">
          <button class="link-btn info">查看</button>
          <button class="link-btn" :class="n.removable ? 'danger' : 'disabled'">移除</button>
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

const networks = [
  { name: 'bridge', driver: 'bridge', scope: 'local', subnet: '172.17.0.0/16', gateway: '172.17.0.1', containers: 3, removable: false },
  { name: 'app-network', driver: 'bridge', scope: 'local', subnet: '192.168.1.0/24', gateway: '192.168.1.1', containers: 5, removable: true },
  { name: 'host', driver: 'host', scope: 'local', subnet: '-', gateway: '-', containers: 0, removable: false }
]

const stats = [
  { label: '总数:', value: '5' },
  { label: '自定义:', value: '3' },
  { label: '内置:', value: '2', variant: 'builtin' }
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

.col-name { width: 130px; }
.col-driver { width: 80px; }
.col-scope { width: 70px; }
.col-subnet { width: 140px; }
.col-gateway { width: 120px; }
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

.link-btn.disabled {
  color: var(--text-secondary);
  cursor: default;
}
</style>
