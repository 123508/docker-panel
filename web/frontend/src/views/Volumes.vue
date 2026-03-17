<template>
  <div class="page">
    <div class="page-header">
      <h1 class="page-title">卷</h1>
      <div class="header-actions">
        <button class="create-btn">+ 创建</button>
      </div>
    </div>

    <div class="stats-bar">
      <div class="stat-item">
        <span class="stat-label">总数:</span>
        <span class="stat-value">8</span>
      </div>
      <div class="stat-item">
        <span class="stat-label">已使用:</span>
        <span class="stat-value">3.2 GB</span>
      </div>
      <div class="stat-item">
        <span class="stat-label">未使用:</span>
        <span class="stat-value dangling">2</span>
      </div>
    </div>

    <div class="table-wrap">
      <div class="table-head">
        <span class="th col-name">名称</span>
        <span class="th col-driver">驱动</span>
        <span class="th col-mount">挂载点</span>
        <span class="th col-size">大小</span>
        <span class="th col-containers">容器</span>
        <span class="th" style="flex: 1">操作</span>
      </div>
      <div v-for="v in volumes" :key="v.name" class="table-row">
        <span class="td col-name">{{ v.name }}</span>
        <span class="td col-driver td-muted">{{ v.driver }}</span>
        <span class="td col-mount td-dim">{{ v.mountpoint }}</span>
        <span class="td col-size">{{ v.size }}</span>
        <div class="td col-containers">
          <span class="container-badge" :class="v.containers > 0 ? 'active' : 'empty'">
            {{ v.containers }}
          </span>
        </div>
        <div class="td td-actions" style="flex: 1">
          <button class="link-btn info">查看</button>
          <button class="link-btn danger">移除</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const volumes = [
  { name: 'postgres_data', driver: 'local', mountpoint: '/var/lib/docker/volumes/postgres_data/_data', size: '1.2 GB', containers: 1 },
  { name: 'redis_cache', driver: 'local', mountpoint: '/var/lib/docker/volumes/redis_cache/_data', size: '512 MB', containers: 1 },
  { name: 'app_logs', driver: 'local', mountpoint: '/var/lib/docker/volumes/app_logs/_data', size: '256 MB', containers: 0 }
]
</script>

<style scoped>
.page {
  padding: var(--page-padding-y) var(--page-padding-x);
  display: flex;
  flex-direction: column;
  gap: 24px;
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

.stats-bar {
  display: flex;
  gap: 24px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.stat-label {
  font-size: 11px;
  letter-spacing: 1px;
  color: var(--text-secondary);
}

.stat-value {
  font-size: 14px;
  font-weight: 700;
  color: var(--text-primary);
}

.stat-value.dangling {
  color: var(--accent);
}

.table-wrap {
  background: var(--bg-card);
  flex: 1;
}

.table-head {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 0 20px;
  height: 44px;
  background: var(--bg-card-header);
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
  padding: 0 20px;
  height: 56px;
  transition: background 0.15s;
}

.table-row:hover {
  background: var(--table-row-hover);
}

.td {
  font-size: 12px;
  color: var(--text-primary);
}

.col-name { width: 140px; }
.col-driver { width: 80px; }
.col-mount { width: 320px; }
.col-size { width: 80px; }
.col-containers { width: 100px; }

.td-muted {
  color: var(--text-muted);
}

.td-dim {
  color: var(--text-secondary);
  font-size: 11px;
}

.container-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 20px;
  font-size: 11px;
  font-weight: 700;
}

.container-badge.active {
  background: rgba(34, 197, 94, 0.125);
  color: var(--color-success);
}

.container-badge.empty {
  background: rgba(113, 113, 122, 0.125);
  color: var(--text-secondary);
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
</style>
