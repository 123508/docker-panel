<template>
  <div class="content-page">
    <div class="content-header">
      <div class="title-wrap">
        <h1 class="content-title">编排项目内容</h1>
      </div>
      <div class="content-actions">
        <dp-button text="+ 新建内容" size="medium" type="primary" @click="goCreate" />
      </div>
    </div>

    <div class="content-stats">
      <div class="cstat">
        <span class="cstat-label">服务:</span>
        <span class="cstat-value">{{ state.stats.total }}</span>
      </div>
      <div class="cstat">
        <span class="cstat-label">运行中:</span>
        <span class="cstat-value success">{{ state.stats.running }}</span>
      </div>
      <div class="cstat">
        <span class="cstat-label">已停止:</span>
        <span class="cstat-value danger">{{ state.stats.stopped }}</span>
      </div>
    </div>

    <div v-if="state.loading" class="placeholder">加载中...</div>
    <template v-else-if="state.services.length === 0">
      <div class="empty-state">
        <div class="empty-icon">13</div>
        <p class="empty-text">暂无服务数据</p>
        <p class="empty-sub">该项目下未发现运行中的服务，请检查项目是否已启动</p>
      </div>
    </template>
    <template v-else>
      <div class="content-table-wrap">
        <div class="ct-header">
          <span class="cth-name">服务名称</span>
          <span class="cth-status">状态</span>
          <span class="cth-image">镜像</span>
          <span class="cth-containers">容器</span>
          <span class="cth-replicas">副本</span>
          <span class="cth-actions">操作</span>
        </div>
        <div v-for="svc in state.services" :key="svc.name" class="ct-row">
          <span class="ct-cell ct-name">{{ svc.name }}</span>
          <div class="ct-cell ct-status">
            <span
              class="status-dot"
              :class="{ running: isRunning(svc.status), stopped: !isRunning(svc.status) }"
            ></span>
            <span :class="statusColor(svc.status)">{{ statusText(svc.status) }}</span>
          </div>
          <span class="ct-cell ct-image muted">{{ svc.image }}</span>
          <div class="ct-cell ct-containers-badge">
            <span class="svc-badge">{{ svc.containers?.length || 0 }}</span>
          </div>
          <span class="ct-cell ct-replicas">{{ svc.replicas }}</span>
          <div class="ct-cell ct-actions">
            <dp-button text="详情" size="small" variant="text" type="info" @click="viewService(svc.name)" />
            <dp-button
              v-if="isRunning(svc.status)"
              text="停止"
              size="small"
              variant="text"
              type="info"
              @click="stopService(svc.name)"
            />
            <dp-button
              v-else
              text="启动"
              size="small"
              variant="text"
              type="primary"
              @click="startService(svc.name)"
            />
            <dp-button text="重启" size="small" variant="text" type="info" @click="restartService(svc.name)" />
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { getComposePS } from '@/services/modules/compose'
import { call } from '@/services/common'
import { ElMessage } from 'element-plus'
import DpButton from '@/components/dp-button.vue'

const router = useRouter()
const route = useRoute()
const projectName = route.params.name as string

const state = reactive({
  loading: false,
  project: null as any,
  services: [] as any[],
  stats: { total: '0', running: '0', stopped: '0' }
})

const loadData = async () => {
  state.loading = true
  try {
    const svcList = await getComposePS(projectName)
    state.services = Array.isArray(svcList) ? svcList : []
  } catch {
    state.services = []
  }

  try {
    const res = await call('GET', '/api/v1/compose/projects', null).unwrap()
    const list = Array.isArray(res) ? res : []
    state.project = list.find((p: any) => p.name === projectName) || null
  } catch {
    state.project = null
  }

  updateStats()
  state.loading = false
}

const updateStats = () => {
  const total = state.services.length
  const running = state.services.filter(s => isRunning(s.status)).length
  const stopped = total - running
  state.stats.total = String(total)
  state.stats.running = String(running)
  state.stats.stopped = String(stopped)
}

const isRunning = (s: string) => {
  const lower = (s || '').toLowerCase()
  return lower.includes('running') || lower.includes('up')
}

const statusText = (s: string) => {
  const lower = (s || '').toLowerCase()
  if (lower.includes('running') || lower.includes('up')) return '运行中'
  if (lower.includes('exit') || lower.includes('stopped')) return '已停止'
  return s || '未知'
}

const statusColor = (s: string) => {
  if (isRunning(s)) return 'text-success'
  return 'text-danger'
}

const goCreate = () => router.push('/dashboard/compose/create')

const viewService = (svcName: string) => {
  ElMessage.info(`服务: ${svcName}`)
}

const stopService = async (svcName: string) => {
  try {
    await call('POST', `/api/v1/compose/projects/${projectName}/stop`, { services: [svcName] }).unwrap()
    ElMessage.success(`服务 ${svcName} 已停止`)
    await loadData()
  } catch (e: any) {
    ElMessage.error(e?.message || `停止 ${svcName} 失败`)
  }
}

const startService = async (svcName: string) => {
  try {
    await call('POST', `/api/v1/compose/projects/${projectName}/up`, { services: [svcName] }).unwrap()
    ElMessage.success(`服务 ${svcName} 已启动`)
    await loadData()
  } catch (e: any) {
    ElMessage.error(e?.message || `启动 ${svcName} 失败`)
  }
}

const restartService = async (svcName: string) => {
  try {
    await call('POST', `/api/v1/compose/projects/${projectName}/restart`, { services: [svcName] }).unwrap()
    ElMessage.success(`服务 ${svcName} 已重启`)
    await loadData()
  } catch (e: any) {
    ElMessage.error(e?.message || `重启 ${svcName} 失败`)
  }
}

onMounted(() => { loadData() })
</script>

<style scoped>
.content-page {
  padding: var(--page-padding-y) var(--page-padding-x);
  height: 100vh;
  display: flex;
  flex-direction: column;
  gap: 24px;
  overflow: auto;
}

.content-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.title-wrap {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.content-title {
  margin: 0;
  font-family: var(--font-display);
  font-size: 40px;
  line-height: 1;
  font-weight: 700;
  letter-spacing: -1px;
  color: var(--text-primary);
}

.content-actions {
  display: flex;
  gap: 12px;
}

.content-stats {
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

.content-table-wrap {
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
.cth-image,
.cth-containers,
.cth-replicas,
.cth-actions {
  font-family: var(--font-mono);
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 1px;
  color: var(--text-secondary);
}

.cth-name { width: 140px; }
.cth-status { width: 100px; }
.cth-image { flex: 1; min-width: 0; }
.cth-containers { width: 60px; }
.cth-replicas { width: 60px; }
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

.ct-name { width: 140px; font-weight: 700; }
.ct-status { width: 100px; display: flex; align-items: center; gap: 6px; }
.ct-image { flex: 1; font-size: 11px; }
.ct-containers-badge { width: 60px; }
.ct-replicas { width: 60px; }
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
