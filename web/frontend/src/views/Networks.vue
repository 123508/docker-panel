<template>
  <dp-page-header title="网络" gap="24px">
    <template #actions>
      <div class="header-actions">
        <dp-button text="+ 创建" size="large" type="primary" class="btn" @click="createNetwork" />
      </div>
    </template>

    <dp-stat-bar :items="state.stats" />

    <dp-data-table :bordered="false" header-bg row-gap fixed-row-height :columns="columns">
      <div v-for="n in pagedNetworks" :key="n.name" class="table-row">
        <span class="td col-name">{{ n.name }}</span>
        <span class="td col-driver td-muted">{{ n.driver }}</span>
        <span class="td col-scope td-muted">{{ n.scope }}</span>
        <span class="td col-subnet td-dim">{{ n.subnet }}</span>
        <span class="td col-gateway td-dim">{{ n.gateway }}</span>
        <div class="td col-containers">
          <dp-count-badge :count="n.containers" />
        </div>
        <div class="td td-actions" style="flex: 1">
          <dp-button text="查看" size="small" variant="text" type="info" @click="inspectNetwork(n.id)" />
          <dp-button
            text="移除"
            size="small"
            variant="text"
            type="danger"
            :disabled="n.removable"
            @click="removeNetwork(n.id)"
          />
        </div>
      </div>
    </dp-data-table>

    <dp-pagination
      :page="state.page"
      :total="state.networks.length"
      :page-size="state.pageSize"
      @change="state.page = $event"
    />
  </dp-page-header>
</template>

<script setup lang="ts">
import DpPageHeader from '@/components/dp-page-header.vue'
import DpStatBar from '@/components/dp-stat-bar.vue'
import DpDataTable from '@/components/dp-data-table.vue'
import DpCountBadge from '@/components/dp-count-badge.vue'
import DpButton from '@/components/dp-button.vue'
import DpPagination from '@/components/dp-pagination.vue'
import { NetworkState } from '@/composables/Networks'

const { state, pagedNetworks, createNetwork, inspectNetwork, removeNetwork } = NetworkState()

const columns = [
  { key: 'name', label: '名称', width: 130 },
  { key: 'driver', label: '驱动', width: 80 },
  { key: 'scope', label: '范围', width: 70 },
  { key: 'subnet', label: '子网', width: 140 },
  { key: 'gateway', label: '网关', width: 120 },
  { key: 'containers', label: '容器', width: 100 },
  { key: 'actions', label: '操作', flex: 1 }
]
</script>

<style scoped>
.btn {
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
