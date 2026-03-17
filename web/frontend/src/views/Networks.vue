<template>
  <div class="page">
    <div class="page-header">
      <h1 class="page-title">网络</h1>
      <div class="header-actions">
        <button class="create-btn">+ 创建</button>
      </div>
    </div>

    <div class="stats-bar">
      <div class="stat-item">
        <span class="stat-label">总数:</span>
        <span class="stat-value">5</span>
      </div>
      <div class="stat-item">
        <span class="stat-label">自定义:</span>
        <span class="stat-value">3</span>
      </div>
      <div class="stat-item">
        <span class="stat-label">内置:</span>
        <span class="stat-value builtin">2</span>
      </div>
    </div>

    <div class="table-wrap">
      <div class="table-head">
        <span class="th col-name">名称</span>
        <span class="th col-driver">驱动</span>
        <span class="th col-scope">范围</span>
        <span class="th col-subnet">子网</span>
        <span class="th col-gateway">网关</span>
        <span class="th col-containers">容器</span>
        <span class="th" style="flex: 1">操作</span>
      </div>
      <div v-for="n in networks" :key="n.name" class="table-row">
        <span class="td col-name">{{ n.name }}</span>
        <span class="td col-driver td-muted">{{ n.driver }}</span>
        <span class="td col-scope td-muted">{{ n.scope }}</span>
        <span class="td col-subnet td-dim">{{ n.subnet }}</span>
        <span class="td col-gateway td-dim">{{ n.gateway }}</span>
        <div class="td col-containers">
          <span class="container-badge" :class="n.containers > 0 ? 'active' : 'empty'">
            {{ n.containers }}
          </span>
        </div>
        <div class="td td-actions" style="flex: 1">
          <button class="link-btn info">查看</button>
          <button class="link-btn" :class="n.removable ? 'danger' : 'disabled'">移除</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const networks = [
  { name: 'bridge', driver: 'bridge', scope: 'local', subnet: '172.17.0.0/16', gateway: '172.17.0.1', containers: 3, removable: false },
  { name: 'app-network', driver: 'bridge', scope: 'local', subnet: '192.168.1.0/24', gateway: '192.168.1.1', containers: 5, removable: true },
  { name: 'host', driver: 'host', scope: 'local', subnet: '-', gateway: '-', containers: 0, removable: false }
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

.stat-value.builtin {
  color: var(--text-muted);
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

.link-btn.disabled {
  color: var(--text-secondary);
  cursor: default;
}
</style>
