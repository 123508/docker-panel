<template>
  <dp-page-header title="系统概览">
    <template #actions>
      <button class="refresh-btn" @click="loadData" :disabled="state.loading || state.actionLoading">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <polyline points="23 4 23 10 17 10" />
          <path d="M20.49 15a9 9 0 1 1-2.12-9.36L23 10" />
        </svg>
      </button>
    </template>

    <div class="metrics">
      <div v-for="m in state.metrics" :key="m.label" class="metric-card">
        <div class="metric-header">
          <span class="metric-label">{{ m.label }}</span>
          <span class="metric-sub" :class="m.subClass">{{ m.sub }}</span>
        </div>
        <div class="metric-value">{{ m.value }}</div>
      </div>
    </div>

    <div class="section">
      <div class="section-header">
        <h2 class="section-title">活跃容器</h2>
        <router-link to="/containers" class="view-all">查看全部 →</router-link>
      </div>

      <dp-data-table :bordered="true" :flex="false" height="54vh">
        <template #head>
          <span class="th" style="width: 280px">名称</span>
          <span class="th" style="width: 200px">镜像</span>
          <span class="th" style="width: 120px">状态</span>
          <span class="th" style="width: 150px">端口</span>
          <span class="th" style="flex: 1">操作</span>
        </template>

        <div v-for="c in pagedContainers" :key="c.fullId" class="table-row">
          <div class="td td-name" style="width: 280px">
            <div class="container-icon"></div>
            <div class="name-col">
              <span class="name-text">{{ c.name }}</span>
              <span class="name-id">{{ c.id }}</span>
            </div>
          </div>
          <span class="td td-image" style="width: 200px">{{ c.image }}</span>
          <div class="td" style="width: 120px">
            <dp-status-badge :status="c.running ? 'running' : 'stopped'">
              {{ c.status }}
            </dp-status-badge>
          </div>
          <span class="td td-port" style="width: 150px">{{ c.port }}</span>
          <div class="td td-actions" style="flex: 1">
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
    </div>
  </dp-page-header>
</template>

<script setup lang="ts">
import DpPageHeader from '@/components/dp-page-header.vue'
import DpDataTable from '@/components/dp-data-table.vue'
import DpStatusBadge from '@/components/dp-status-badge.vue'
import DpButton from '@/components/dp-button.vue'
import DpPagination from '@/components/dp-pagination.vue'
import { DashboardState } from '@/composables/Dashboard'

const { state, pagedContainers, loadData, startContainer, stopContainer, restartContainer, removeContainer } = DashboardState()
</script>

<style scoped>
.refresh-btn {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-card);
  border: var(--border);
  color: var(--text-secondary);
  transition: color 0.2s;
}

.refresh-btn:hover {
  color: var(--text-primary);
}

.metrics {
  display: flex;
  gap: 16px;
}

.metric-card {
  flex: 1;
  background: var(--bg-card);
  border: var(--border);
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.metric-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.metric-label {
  font-family: var(--font-mono);
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 1px;
  color: var(--text-secondary);
  text-transform: uppercase;
}

.metric-sub {
  font-family: var(--font-mono);
  font-size: 11px;
  color: var(--text-muted);
}

.metric-sub.success {
  color: var(--color-success);
}

.metric-value {
  font-family: var(--font-display);
  font-size: 32px;
  font-weight: 700;
  color: var(--text-primary);
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.section-title {
  font-family: var(--font-display);
  font-size: 18px;
  font-weight: 700;
  color: var(--text-primary);
}

.view-all {
  font-family: var(--font-mono);
  font-size: 11px;
  color: var(--accent);
  letter-spacing: 1px;
  font-weight: 700;
}

.view-all:hover {
  text-decoration: underline;
}

.th {
  font-family: var(--font-mono);
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 1px;
  color: var(--text-secondary);
}

.table-row {
  display: flex;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid transparent;
  transition: background 0.15s;
}

.table-row:hover {
  background: var(--table-row-hover);
}

.table-row:not(:last-child) {
  border-bottom-color: var(--border-color);
}

.td {
  font-family: var(--font-mono);
  font-size: 14px;
}

.td-name {
  display: flex;
  align-items: center;
  gap: 12px;
}

.container-icon {
  width: 32px;
  height: 32px;
  background: var(--bg-card-header);
  border: var(--border);
  flex-shrink: 0;
}

.name-col {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.name-text {
  color: var(--text-primary);
  font-size: 14px;
}

.name-id {
  color: var(--text-dim);
  font-size: 11px;
}

.td-image {
  color: var(--text-muted);
}

.td-port {
  color: var(--text-muted);
  font-size: 11px;
}

.td-actions {
  display: flex;
  gap: 12px;
}
</style>
