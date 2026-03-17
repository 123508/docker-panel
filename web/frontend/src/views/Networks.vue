<template>
  <PageHeader title="网络" gap="24px">
    <template #actions>
      <div class="header-actions">
        <BaseButton text="+ 创建" size="large" type="primary" class="btn"/>
      </div>
    </template>

    <StatBar :items="stats" />

    <DataTable :bordered="false" header-bg row-gap fixed-row-height :columns="columns">

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
          <BaseButton text="查看" size="small" variant="text" type="info" />
          <BaseButton
              text="移除"
              size="small"
              variant="text"
              type="danger"
              :disabled="n.removable"
          />
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
import BaseButton from "@/components/BaseButton.vue";

const columns = [
  { key: 'name', label: '名称', width: 130 },
  { key: 'driver', label: '驱动', width: 80 },
  { key: 'scope', label: '范围', width: 70 },
  { key: 'subnet', label: '子网', width: 140 },
  { key: 'gateway', label: '网关', width: 120 },
  { key: 'containers', label: '容器', width: 100 },
  { key: 'actions', label: '操作', flex: 1 }
]

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

.btn{
  width: 140px;
  display: flex;
  align-items: center;
  justify-content: center;
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

</style>
