<template>
  <div class="page">
    <div class="page-header">
      <h1 class="page-title">系统概览</h1>
      <button class="refresh-btn">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <polyline points="23 4 23 10 17 10" />
          <path d="M20.49 15a9 9 0 1 1-2.12-9.36L23 10" />
        </svg>
      </button>
    </div>

    <div class="metrics">
      <div v-for="m in metrics" :key="m.label" class="metric-card">
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

      <div class="table-wrap">
        <div class="table-head">
          <span class="th" style="width: 280px">名称</span>
          <span class="th" style="width: 200px">镜像</span>
          <span class="th" style="width: 120px">状态</span>
          <span class="th" style="width: 150px">端口</span>
          <span class="th" style="flex: 1">操作</span>
        </div>
        <div v-for="c in containers" :key="c.name" class="table-row">
          <div class="td td-name" style="width: 280px">
            <div class="container-icon"></div>
            <div class="name-col">
              <span class="name-text">{{ c.name }}</span>
              <span class="name-id">{{ c.id }}</span>
            </div>
          </div>
          <span class="td td-image" style="width: 200px">{{ c.image }}</span>
          <div class="td" style="width: 120px">
            <span class="status-badge" :class="c.status === '运行中' ? 'running' : 'stopped'">
              <span class="status-dot"></span>
              {{ c.status }}
            </span>
          </div>
          <span class="td td-port" style="width: 150px">{{ c.port }}</span>
          <div class="td td-actions" style="flex: 1">
            <button class="action-btn">停止</button>
            <button class="action-btn muted">重启</button>
            <button class="action-btn muted">日志</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const metrics = [
  { label: '容器', value: '24', sub: '18 运行中', subClass: 'success' },
  { label: '镜像', value: '42', sub: '12.4 GB', subClass: 'muted' },
  { label: '卷', value: '8', sub: '3.2 GB 已使用', subClass: 'muted' },
  { label: '网络', value: '5', sub: '3 自定义', subClass: 'muted' }
]

const containers = [
  { name: 'nginx-web-01', id: 'a3f8d92b', image: 'nginx:latest', status: '运行中', port: '80:80, 443:443' },
  { name: 'postgres-db', id: 'e7b2c41a', image: 'postgres:14', status: '运行中', port: '5432:5432' },
  { name: 'redis-cache', id: 'f4c9e3d2', image: 'redis:alpine', status: '已停止', port: '6379:6379' }
]
</script>

<style scoped>
.page {
  padding: var(--page-padding-y) var(--page-padding-x);
  display: flex;
  flex-direction: column;
  gap: var(--section-gap);
  min-height: 100vh;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.page-title {
  font-family: var(--font-display);
  font-size: 40px;
  font-weight: 700;
  letter-spacing: -1px;
  color: var(--text-primary);
}

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

.table-wrap {
  background: var(--bg-card);
  border: var(--border);
  overflow: hidden;
}

.table-head {
  display: flex;
  align-items: center;
  padding: 0 20px;
  height: 44px;
  border-bottom: var(--border);
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

.status-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: 10px;
  font-weight: 700;
  letter-spacing: 1px;
  padding: 3px 6px;
}

.status-badge.running {
  color: var(--color-success);
}

.status-badge.stopped {
  color: var(--color-danger);
}

.status-dot {
  width: 6px;
  height: 6px;
  border-radius: 0;
  background: currentColor;
}

.td-port {
  color: var(--text-muted);
  font-size: 11px;
}

.td-actions {
  display: flex;
  gap: 12px;
}

.action-btn {
  font-family: var(--font-mono);
  font-size: 10px;
  font-weight: 700;
  letter-spacing: 1px;
  color: var(--text-primary);
  background: var(--bg-card-header);
  border: var(--border);
  padding: 6px 10px;
  transition: background 0.15s;
}

.action-btn:hover {
  background: var(--bg-body);
}

.action-btn.muted {
  color: var(--text-secondary);
}
</style>
