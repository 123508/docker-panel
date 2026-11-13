<template>
  <dp-page-header title="系统概览">
    <template #actions>
      <button class="refresh-btn" title="刷新数据" aria-label="刷新数据" @click="refreshData" :disabled="state.loading || state.actionLoading">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <polyline points="23 4 23 10 17 10" />
          <path d="M20.49 15a9 9 0 1 1-2.12-9.36L23 10" />
        </svg>
      </button>
    </template>

    <div class="overview-intro">
      <div><div class="intro-kicker"><span></span> YOUR DOCKER WORKSPACE</div><h2>一切，尽在掌握。</h2><p>在一个地方查看和管理你的容器资源。</p></div>
      <router-link to="/dashboard/containers" class="intro-link">管理容器 <span>↗</span></router-link>
      <div class="intro-art" aria-hidden="true"><i></i><i></i><i></i><i></i><i></i><i></i></div>
    </div>

    <div class="metrics">
      <div v-for="(m, index) in state.metrics" :key="m.label" class="metric-card">
        <div class="metric-icon">{{ ['▦', '◈', '▤', '⌘'][index] }}</div>
        <div class="metric-header">
          <span class="metric-label">{{ m.label }}</span>
          <span class="metric-sub" :class="m.subClass">{{ m.sub }}</span>
        </div>
        <div class="metric-value">{{ m.value }}</div>
      </div>
    </div>

    <div class="section">
      <div class="section-header">
        <div><div class="section-kicker">RECENT ACTIVITY</div><h2 class="section-title">活跃容器</h2><p class="section-description">最近操作或创建的容器</p></div>
        <router-link to="/dashboard/containers" class="view-all">查看全部 →</router-link>
      </div>

      <dp-data-table :bordered="true" :flex="false" height="54vh" :columns="columns">
        <div v-if="!state.loading && !state.recentContainers.length" class="empty-state">暂无容器，创建容器后即可在这里查看。</div>
        <div v-for="c in state.recentContainers" :key="c.fullId" class="table-row">
          <div class="td td-name" data-label="名称" style="width: 280px">
            <div class="container-icon"></div>
            <div class="name-col">
              <span class="name-text" :title="c.name">{{ c.name }}</span>
              <span class="name-id" :title="c.id">{{ c.id }}</span>
            </div>
          </div>
          <span class="td td-image td-ellipsis" data-label="镜像" style="width: 200px" :title="c.image">{{ c.image }}</span>
          <div class="td" data-label="状态" style="width: 120px">
            <dp-status-badge :status="c.running ? 'running' : 'stopped'">
              {{ c.status }}
            </dp-status-badge>
          </div>
          <span class="td td-port td-ellipsis" data-label="端口" style="width: 150px" :title="c.port">{{ c.port }}</span>
          <div class="td td-actions" data-label="操作" style="flex: 1">
            <template v-if="c.running">
              <dp-button :disabled="state.actionLoading" text="停止" size="small" type="danger" variant="text" @click="stopContainer(c.fullId)" />
              <dp-button :disabled="state.actionLoading" text="重启" size="small" type="info" variant="text" @click="restartContainer(c.fullId)" />
            </template>
            <template v-else>
              <dp-button :disabled="state.actionLoading" text="启动" size="small" type="primary" variant="text" @click="startContainer(c.fullId)" />
              <dp-button :disabled="state.actionLoading" text="移除" size="small" type="danger" variant="text" @click="removeContainer(c.fullId)" />
            </template>
          </div>
        </div>
      </dp-data-table>

    </div>
  </dp-page-header>
</template>

<script setup lang="ts">
import DpPageHeader from '@/components/dp-page-header.vue'
import DpDataTable from '@/components/dp-data-table.vue'
import DpStatusBadge from '@/components/dp-status-badge.vue'
import DpButton from '@/components/dp-button.vue'
import { ElMessage } from 'element-plus'
import { DashboardState } from '@/composables/Dashboard'

const { state, loadData, startContainer, stopContainer, restartContainer, removeContainer } = DashboardState()
const refreshData = async () => {
  if (await loadData()) ElMessage.success('系统概览已更新')
}
const columns = [
  { key: 'name', label: '名称', width: 280 },
  { key: 'image', label: '镜像', width: 200 },
  { key: 'status', label: '状态', width: 120 },
  { key: 'port', label: '端口', width: 150 },
  { key: 'actions', label: '操作', width: 120 }
]
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
  min-width: 0;
  flex: 1;
}

.name-text {
  color: var(--text-primary);
  font-size: 14px;
  display: block;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.name-id {
  color: var(--text-dim);
  font-size: 11px;
  display: block;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* 内容超长时截断显示省略号，鼠标悬停可见完整文本 */
.td-ellipsis {
  display: inline-block;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  vertical-align: middle;
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

<style scoped>
.overview-intro { min-height: 180px; position: relative; overflow: hidden; padding: 30px 34px; border-radius: 16px; background: var(--intro-gradient); color: #fff; display: flex; align-items: center; justify-content: space-between; gap: 20px; }
.intro-kicker { color: var(--intro-kicker); font: 700 10px var(--font-mono); letter-spacing: 1.5px; margin-bottom: 14px; }
.intro-kicker span { display: inline-block; width: 7px; height: 7px; margin-right: 7px; background: #6fdbd0; border-radius: 50%; }
.overview-intro h2 { font-size: 25px; letter-spacing: -.5px; margin-bottom: 8px; }
.overview-intro p { color: var(--intro-text); font-size: 13px; }
.intro-link { position: relative; z-index: 1; flex-shrink: 0; padding: 11px 15px; border: 1px solid #ffffff66; border-radius: 8px; font-size: 12px; font-weight: 700; background: #ffffff13; }
.intro-link:hover { background: #ffffff2b; }
.intro-link span { margin-left: 12px; }
.intro-art { position: absolute; right: 20%; top: -26px; width: 250px; height: 250px; transform: rotate(-18deg); opacity: .13; display: grid; grid-template-columns: repeat(3, 1fr); gap: 9px; }
.intro-art i { border: 2px solid #fff; border-radius: 15px; }
.metrics { display: grid; grid-template-columns: repeat(4, minmax(0, 1fr)); }
.metric-card { border-radius: 14px; padding: 20px 22px; gap: 10px; box-shadow: 0 8px 24px #07152508; }
.metric-icon { width: 35px; height: 35px; display: grid; place-items: center; border-radius: 10px; background: var(--create-accent-soft); color: var(--accent); font-size: 20px; margin-bottom: 4px; }
.metric-label { font-family: var(--font-display), sans-serif; font-size: 13px; letter-spacing: 0; }
.metric-value { font-size: 30px; line-height: 1; }
.metric-sub { font-size: 11px; }
.section-kicker { color: var(--text-dim); font: 700 10px var(--font-mono); letter-spacing: 1.3px; margin-bottom: 5px; }
.section-description { font-size: 12px; color: var(--text-secondary); margin-top: 4px; }
.section-header { margin-bottom: 18px; }
.view-all { letter-spacing: 0; font-family: var(--font-display), sans-serif; font-size: 12px; }
.refresh-btn { border-radius: 10px; }
.container-icon { border-radius: 9px; background: var(--create-accent-soft); position: relative; }
.container-icon::after { content: '▦'; color: var(--accent); position: absolute; inset: 0; display: grid; place-items: center; font-size: 18px; }
.empty-state { color: var(--text-secondary); padding: 52px 20px; text-align: center; font-size: 13px; }
@media (max-width: 1000px) { .metrics { grid-template-columns: repeat(2, minmax(0, 1fr)); } }
@media (max-width: 550px) { .overview-intro { padding: 24px; align-items: flex-start; flex-direction: column; } .overview-intro h2 { font-size: 22px; } .intro-art { right: -40px; } .metric-card { padding: 16px; } }
</style>
