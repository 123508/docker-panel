<template>
  <div class="page">
    <div class="page-header">
      <h1 class="page-title">容器</h1>
      <div class="header-actions">
        <div class="search-box">
          <span class="search-text">搜索...</span>
        </div>
        <button class="filter-btn">筛选</button>
        <button class="create-btn">+ 创建</button>
      </div>
    </div>

    <div class="stats-bar">
      <div class="stat-item">
        <span class="stat-label">总数:</span>
        <span class="stat-value">24</span>
      </div>
      <div class="stat-item">
        <span class="stat-label">运行中:</span>
        <span class="stat-value running">18</span>
      </div>
      <div class="stat-item">
        <span class="stat-label">已停止:</span>
        <span class="stat-value stopped">6</span>
      </div>
    </div>

    <div class="table-wrap">
      <div class="table-head">
        <span class="th" style="width: 60px"></span>
        <span class="th" style="width: 220px">名称</span>
        <span class="th" style="width: 180px">镜像</span>
        <span class="th" style="width: 120px">状态</span>
        <span class="th" style="width: 140px">端口</span>
        <span class="th" style="width: 120px">创建时间</span>
        <span class="th" style="flex: 1">操作</span>
      </div>
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
          <span class="status-badge" :class="c.running ? 'running' : 'stopped'">
            <span class="status-dot"></span>
            {{ c.running ? '运行中' : '已停止' }}
          </span>
        </div>
        <span class="td td-muted-sm" style="width: 140px">{{ c.port }}</span>
        <span class="td td-muted-sm" style="width: 120px">{{ c.created }}</span>
        <div class="td td-actions" style="flex: 1">
          <template v-if="c.running">
            <button class="action-btn">停止</button>
            <button class="action-btn muted">重启</button>
            <button class="action-btn muted">日志</button>
          </template>
          <template v-else>
            <button class="action-btn accent">启动</button>
            <button class="action-btn danger">移除</button>
            <button class="action-btn muted">日志</button>
          </template>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const containers = [
  { name: 'nginx-web-01', id: 'a3f8d92b', image: 'nginx:latest', running: true, port: '80:80, 443:443', created: '2 小时前' },
  { name: 'postgres-db', id: 'e7b2c41a', image: 'postgres:14', running: true, port: '5432:5432', created: '5 天前' },
  { name: 'redis-cache', id: 'f4c9e3d2', image: 'redis:alpine', running: false, port: '6379:6379', created: '1 天前' }
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

.header-actions {
  display: flex;
  gap: 12px;
}

.search-box {
  display: flex;
  align-items: center;
  background: var(--bg-card);
  border: var(--border);
  padding: 10px 14px;
  width: 240px;
}

.search-text {
  font-size: 11px;
  color: var(--text-dim);
}

.filter-btn {
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 1px;
  color: var(--text-primary);
  background: var(--bg-card);
  border: var(--border);
  padding: 10px 14px;
}

.create-btn {
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 1px;
  color: var(--accent-text);
  background: var(--accent);
  padding: 10px 18px;
}

.create-btn:hover {
  background: var(--accent-hover);
}

.stats-bar {
  display: flex;
  gap: 12px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 8px;
  background: var(--bg-card);
  border: var(--border);
  padding: 8px 12px;
}

.stat-label {
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 1px;
  color: var(--text-secondary);
}

.stat-value {
  font-size: 11px;
  font-weight: 700;
  color: var(--text-primary);
}

.stat-value.running {
  color: var(--color-success);
}

.stat-value.stopped {
  color: var(--color-danger);
}

.table-wrap {
  background: var(--bg-card);
  border: var(--border);
  flex: 1;
}

.table-head {
  display: flex;
  align-items: center;
  padding: 0 20px;
  height: 44px;
  border-bottom: var(--border);
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
  background: currentColor;
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
