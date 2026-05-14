<template>
  <dp-page-header title="编排" gap="24px">
    <template #actions>
      <dp-button text="+ 新建编排" size="medium" type="primary" class="create-btn" @click="uploadNew" />
    </template>

    <div class="compose-stats">
      <div v-for="s in state.stats" :key="s.label" class="cstat">
        <span class="cstat-label">{{ s.label }}</span>
        <span class="cstat-value" :class="s.variant">{{ s.value }}</span>
      </div>
    </div>

    <div v-if="state.loading" class="placeholder">加载中...</div>
    <template v-else-if="state.projects.length === 0">
      <div class="empty-state">
        <div class="empty-icon">06</div>
        <p class="empty-text">暂无编排项目</p>
        <p class="empty-sub">点击「+ 新建编排」创建第一个 Docker Compose 项目</p>
      </div>
    </template>
    <template v-else>
      <div class="compose-table-wrap">
        <div class="ct-header">
          <span class="cth-name">项目名称</span>
          <span class="cth-status">状态</span>
          <span class="cth-services">服务</span>
          <span class="cth-path">路径</span>
          <span class="cth-created">创建时间</span>
          <span class="cth-actions">操作</span>
        </div>
        <div v-for="p in state.projects" :key="p.name" class="ct-row">
          <span class="ct-cell ct-name">{{ p.name }}</span>
          <div class="ct-cell ct-status">
            <span
              class="status-dot"
              :class="{ running: p.status === 'running', stopped: p.status === 'stopped' || p.status === 'exited' }"
            ></span>
            <span :class="statusColor(p.status)">{{ statusText(p.status) }}</span>
          </div>
          <div class="ct-cell ct-services-badge">
            <span class="svc-badge">{{ p.services }}</span>
          </div>
          <span class="ct-cell ct-path muted">{{ p.file_path }}</span>
          <span class="ct-cell ct-created muted">{{ formatTime(p.created_at) }}</span>
          <div class="ct-cell ct-actions">
            <dp-button text="查看" size="small" variant="text" type="info" @click="viewProject(p.name)" />
            <dp-button
              v-if="p.status === 'running'"
              text="停止"
              size="small"
              variant="text"
              type="info"
              @click="stopProject(p.name)"
            />
            <dp-button
              v-else
              text="启动"
              size="small"
              variant="text"
              type="primary"
              @click="startProject(p.name)"
            />
            <dp-button text="删除" size="small" variant="text" type="danger" @click="removeProject(p.name)" />
          </div>
        </div>
      </div>
    </template>
  </dp-page-header>
</template>

<script setup lang="ts">
import DpPageHeader from '@/components/dp-page-header.vue'
import DpButton from '@/components/dp-button.vue'
import { ComposeProjectState } from '@/composables/ComposeProjects'

const { state, uploadNew, startProject, stopProject, removeProject, viewProject } = ComposeProjectState()

const statusText = (s: string) => {
  const map: Record<string, string> = { running: '运行中', stopped: '已停止', exited: '已退出', unknown: '未知' }
  return map[s] || s
}

const statusColor = (s: string) => {
  if (s === 'running') return 'text-success'
  if (s === 'stopped' || s === 'exited') return 'text-danger'
  return 'text-muted'
}

const formatTime = (t: string) => {
  if (!t) return '-'
  return new Date(t).toLocaleString()
}
</script>

<style scoped>
.create-btn {
  width: 140px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.compose-stats {
  display: flex;
  gap: 24px;
}

.cstat {
  display: flex;
  gap: 8px;
  align-items: baseline;
}

.cstat-label {
  font-family: var(--font-mono);
  font-size: 11px;
  font-weight: normal;
  color: var(--text-secondary);
}

.cstat-value {
  font-family: var(--font-mono);
  font-size: 14px;
  font-weight: 700;
  color: var(--text-primary);
}

.cstat-value.success { color: var(--color-success); }
.cstat-value.danger { color: var(--color-danger); }

.placeholder {
  padding: 20px;
  border: var(--border);
  background: var(--bg-card);
  color: var(--text-muted);
}

.empty-state {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 80px 0;
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

.compose-table-wrap {
  background: var(--bg-card);
  border: var(--border);
}

.ct-header {
  display: flex;
  align-items: center;
  gap: 16px;
  min-height: 44px;
  padding: 0 20px;
  background: var(--bg-card-header);
}

.cth-name,
.cth-status,
.cth-services,
.cth-path,
.cth-created,
.cth-actions {
  font-family: var(--font-mono);
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 1px;
  color: var(--text-secondary);
}

.cth-name { width: 160px; }
.cth-status { width: 100px; }
.cth-services { width: 60px; }
.cth-path { flex: 1; min-width: 0; }
.cth-created { width: 160px; }
.cth-actions { width: 200px; }

.ct-row {
  display: flex;
  align-items: center;
  gap: 16px;
  min-height: 56px;
  padding: 0 20px;
  border-top: var(--border);
  transition: background 0.15s;
}

.ct-row:hover {
  background: var(--table-row-hover);
}

.ct-cell {
  font-family: var(--font-mono);
  font-size: 12px;
  color: var(--text-primary);
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.ct-name { width: 160px; font-weight: 700; }
.ct-status { width: 100px; display: flex; align-items: center; gap: 6px; }
.ct-services-badge { width: 60px; }
.ct-path { flex: 1; font-size: 11px; }
.ct-created { width: 160px; font-size: 11px; }
.ct-actions { width: 200px; display: flex; gap: 8px; }

.muted { color: var(--text-muted); }
.text-success { color: var(--color-success); }
.text-danger { color: var(--color-danger); }
.text-muted { color: var(--text-muted); }

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--text-muted);
  flex-shrink: 0;
}

.status-dot.running { background: var(--color-success); }
.status-dot.stopped { background: var(--color-danger); }

.svc-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 28px;
  height: 20px;
  padding: 0 6px;
  background: var(--badge-bg);
  color: var(--badge-text);
  font-size: 11px;
  font-weight: 700;
  border-radius: 2px;
}
</style>
