<template>
  <dp-page-header title="容器">
    <template #actions>
      <dp-container-header-actions class="header-actions">
        <dp-search-input v-model="state.query" />
        <dp-button text="刷新" size="medium" variant="outlined" @click="refreshList" />
        <dp-button text="+ 创建" size="medium" type="primary" class="create-btn" @click="goCreate" />
      </dp-container-header-actions>
    </template>

    <dp-stat-bar :items="state.stats" bordered />

    <dp-data-table :bordered="true" :columns="columns">
      <div v-for="c in pagedContainers" :key="c.fullId" class="table-row">
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
        <span class="td td-muted-sm" style="width: 140px">{{ c.created }}</span>
        <div class="td td-actions" style="flex: 1">
          <dp-button text="详情" size="small" type="info" variant="text" @click="goDetail(c.fullId)" />
          <template v-if="c.running">
            <dp-button text="停止" size="small" type="danger" variant="text" @click="stopContainer(c.fullId)" />
            <dp-button text="重启" size="small" type="info" variant="text" @click="restartContainer(c.fullId)" />
          </template>
          <template v-else>
            <dp-button text="启动" size="small" type="primary" variant="text" @click="startContainer(c.fullId)" />
            <dp-button text="移除" size="small" type="danger" variant="text" @click="removeContainer(c.fullId)" />
          </template>
        </div>
      </div>
    </dp-data-table>

    <dp-pagination
      :page="state.page"
      :total="state.containers.length"
      :page-size="state.pageSize"
      @change="state.page = $event"
    />
  </dp-page-header>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import DpPageHeader from '@/components/dp-page-header.vue'
import DpStatBar from '@/components/dp-stat-bar.vue'
import DpDataTable from '@/components/dp-data-table.vue'
import DpStatusBadge from '@/components/dp-status-badge.vue'
import DpButton from '@/components/dp-button.vue'
import DpSearchInput from '@/components/dp-search-input.vue'
import DpPagination from '@/components/dp-pagination.vue'
import DpContainerHeaderActions from '@/components/container/dp-container-header-actions.vue'
import { ContainerState } from '@/composables/Containers'

const router = useRouter()
const { state, pagedContainers, loadData, startContainer, stopContainer, restartContainer, removeContainer } = ContainerState()

const goCreate = () => {
  router.push('/containers/create')
}

const goDetail = (id: string) => {
  router.push(`/containers/${id}`)
}

const refreshList = async () => {
  await loadData()
}

const columns = [
  { key: 'icon', width: 60 },
  { key: 'name', label: '名称', width: 220 },
  { key: 'image', label: '镜像', width: 180 },
  { key: 'status', label: '状态', width: 120 },
  { key: 'port', label: '端口', width: 140 },
  { key: 'created', label: '创建时间', width: 140 },
  { key: 'actions', label: '操作', flex: 1 }
]
</script>

<style scoped>
.header-actions {
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
</style>
