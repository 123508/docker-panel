<template>
  <div class="detail-page">
    <div class="detail-header">
      <div class="title-wrap">
        <h1 class="detail-title">{{ displayName }}</h1>
        <p class="detail-subtitle">容器 ID: {{ shortId }}</p>
      </div>
      <dp-container-header-actions class="detail-actions">
        <dp-button
          text="返回列表"
          size="medium"
          variant="outlined"
          class="btn-back"
          :disabled="state.actionLoading"
          @click="goBack"
        />
        <dp-button
          v-if="isRunning"
          text="停止"
          type="default"
          variant="text"
          size="medium"
          class="btn-stop"
          :disabled="state.actionLoading"
          @click="stop"
        />
        <dp-button
          v-else
          text="启动"
          type="primary"
          size="medium"
          class="btn-start"
          :disabled="state.actionLoading"
          @click="start"
        />
        <dp-button
          text="重启"
          type="primary"
          size="medium"
          class="btn-restart"
          :disabled="state.actionLoading"
          @click="restart"
        />
        <dp-button
          text="删除"
          type="danger"
          variant="filled"
          size="medium"
          class="btn-delete"
          :disabled="state.actionLoading"
          @click="handleRemove"
        />
      </dp-container-header-actions>
    </div>

    <div v-if="state.loading" class="placeholder">加载中...</div>
    <div v-else-if="!state.detail" class="placeholder">未获取到容器详情</div>
    <template v-else>
      <section class="status-card">
        <div class="status-item">
          <span class="item-label">状态</span>
          <dp-status-badge :status="isRunning ? 'running' : 'stopped'">
            {{ isRunning ? '运行中' : '已停止' }}
          </dp-status-badge>
        </div>
        <div class="status-item">
          <span class="item-label">运行时长</span>
          <span class="item-value">{{ uptimeText }}</span>
        </div>
        <div class="status-item">
          <span class="item-label">CPU</span>
          <span class="item-value">--</span>
        </div>
        <div class="status-item">
          <span class="item-label">内存</span>
          <span class="item-value">{{ memoryLimitText }}</span>
        </div>
        <div class="status-item">
          <span class="item-label">网络</span>
          <span class="item-value">{{ networkText }}</span>
        </div>
      </section>

      <section class="section-block">
        <h2 class="section-title">容器信息</h2>
        <dp-container-section-card class="info-grid-card">
          <div class="info-grid">
            <div v-for="row in infoRows" :key="row.key" class="info-row">
              <span class="row-label">{{ row.label }}</span>
              <span class="row-value">{{ row.value }}</span>
            </div>
          </div>
        </dp-container-section-card>
      </section>

      <section class="section-block">
        <h2 class="section-title">容器日志</h2>
        <div class="logs-terminal">
          <p v-for="(line, idx) in logLines" :key="`${idx}-${line.text}`" :class="line.levelClass">
            {{ line.text }}
          </p>
        </div>
      </section>
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import DpButton from '@/components/dp-button.vue'
import DpStatusBadge from '@/components/dp-status-badge.vue'
import DpContainerHeaderActions from '@/components/container/dp-container-header-actions.vue'
import DpContainerSectionCard from '@/components/container/dp-container-section-card.vue'
import { ContainerDetailState } from '@/composables/ContainerDetail'

const route = useRoute()
const router = useRouter()

const containerId = computed(() => {
  const pathId = route.params.id
  if (pathId != null && String(pathId).trim() !== '') return String(pathId)
  const queryId = route.query.id
  if (queryId != null && String(queryId).trim() !== '') return String(queryId)
  return ''
})
const { state, isRunning, displayName, start, stop, restart, remove } = ContainerDetailState(() => containerId.value)

const shortId = computed(() => {
  const id = state.detail?.id || containerId.value
  return String(id || '').slice(0, 12) || 'N/A'
})

const uptimeText = computed(() => {
  const startedAt = state.detail?.state?.started_at
  if (!startedAt || !isRunning.value) return 'N/A'
  const started = new Date(startedAt).getTime()
  if (!Number.isFinite(started) || started <= 0) return 'N/A'
  const diffMinutes = Math.max(0, Math.floor((Date.now() - started) / 60000))
  const hours = Math.floor(diffMinutes / 60)
  const minutes = diffMinutes % 60
  if (hours <= 0) return `${minutes} 分钟`
  return `${hours} 小时 ${minutes} 分钟`
})

const memoryLimitText = computed(() => {
  const bytes = Number(state.detail?.host_config?.memory || 0)
  if (!bytes) return '不限'
  if (bytes < 1024 * 1024 * 1024) return `${Math.round(bytes / (1024 * 1024))}MB`
  return `${(bytes / (1024 * 1024 * 1024)).toFixed(1)}GB`
})

const networkText = computed(() => {
  const networks = state.detail?.network_settings?.networks || {}
  const names = Object.keys(networks)
  if (!names.length) return 'N/A'
  return names.join(', ')
})

const infoRows = computed(() => {
  const detail = state.detail || {}
  const labels = detail?.config?.labels || {}
  return [
    { key: 'id', label: '容器 ID', value: detail.id || 'N/A' },
    { key: 'name', label: '容器名称', value: displayName.value },
    { key: 'image', label: '镜像', value: detail.image || detail.config?.image || 'N/A' },
    { key: 'status', label: '状态', value: detail.state?.status || 'N/A' },
    { key: 'created', label: '创建时间', value: detail.created || 'N/A' },
    { key: 'restartCount', label: '重启次数', value: String(detail.restart_count ?? 0) },
    { key: 'network', label: '网络模式', value: detail.host_config?.network_mode || 'N/A' },
    { key: 'platform', label: '平台', value: detail.platform || 'N/A' },
    { key: 'entrypoint', label: '入口命令', value: (detail.config?.entrypoint || []).join(' ') || 'N/A' },
    { key: 'cmd', label: '启动命令', value: (detail.config?.cmd || []).join(' ') || 'N/A' },
    { key: 'workdir', label: '工作目录', value: detail.config?.working_dir || 'N/A' },
    { key: 'labels', label: '标签', value: Object.keys(labels).length ? JSON.stringify(labels) : 'N/A' }
  ]
})

const logLines = computed(() => {
  const detail = state.detail || {}
  const lines = [
    `[info] container=${displayName.value} status=${detail.state?.status || 'unknown'}`,
    `[info] image=${detail.image || detail.config?.image || 'N/A'}`,
    `[info] command=${(detail.config?.cmd || []).join(' ') || 'N/A'}`,
    `[info] started_at=${detail.state?.started_at || 'N/A'}`,
    `[info] network=${networkText.value}`,
    '[warn] 该页面尚未接入实时日志 API，当前为容器元数据摘要',
    `[info] mounts=${Array.isArray(detail.mounts) ? detail.mounts.length : 0}`
  ]
  return lines.map((text) => ({ text, levelClass: text.startsWith('[warn]') ? 'log-warn' : 'log-info' }))
})

const goBack = () => {
  router.push('/dashboard/containers')
}

const handleRemove = async () => {
  const ok = await remove()
  if (ok) {
    router.push('/dashboard/containers')
  }
}
</script>

<style scoped>
.detail-page {
  padding: var(--page-padding-y) var(--page-padding-x);
  height: 100vh;
  display: flex;
  flex-direction: column;
  gap: var(--section-gap);
  overflow: auto;
}

.detail-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.title-wrap {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.detail-title {
  font-family: var(--font-display);
  font-size: 40px;
  line-height: 1;
  font-weight: 700;
  letter-spacing: -1px;
  color: var(--text-primary);
}

.detail-subtitle {
  font-family: var(--font-mono);
  font-size: 11px;
  color: var(--text-secondary);
}

.detail-actions {
  justify-content: flex-end;
  gap: 20px;
}

.detail-actions :deep(.btn-back.button-outlined) {
  min-width: 116px;
  padding: 0 18px;
  background: var(--bg-card);
}

.detail-actions :deep(.btn-stop.button-text) {
  color: var(--text-primary);
  padding: 0 2px;
}

.detail-actions :deep(.btn-restart.button-filled.button-primary) {
  min-width: 76px;
  padding: 0 20px;
}

.detail-actions :deep(.btn-delete.button-filled.button-danger) {
  min-width: 76px;
  padding: 0 20px;
}

.placeholder {
  padding: 20px;
  border: var(--border);
  background: var(--bg-card);
  color: var(--text-muted);
}

.status-card {
  display: grid;
  grid-template-columns: repeat(5, minmax(0, 1fr));
  gap: 20px;
  background: var(--bg-card);
  border: var(--border);
  padding: 20px;
}

.status-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.item-label {
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 1px;
  color: var(--text-secondary);
}

.item-value {
  color: var(--text-primary);
  font-size: 14px;
}

.section-block {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.section-title {
  font-family: var(--font-display);
  font-size: 18px;
  font-weight: 700;
  color: var(--text-primary);
}

.info-grid-card {
  padding: 0;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.info-row {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 20px;
  border-bottom: var(--border);
}

.info-row:nth-last-child(-n + 2) {
  border-bottom: none;
}

.row-label {
  font-size: 11px;
  font-weight: 700;
  color: var(--text-secondary);
  letter-spacing: 1px;
}

.row-value {
  font-size: 13px;
  color: var(--text-primary);
  overflow-wrap: anywhere;
}

.logs-terminal {
  background: var(--bg-card);
  padding: 16px;
  min-height: 280px;
  display: flex;
  flex-direction: column;
  gap: 4px;
  border: var(--border);
}

.logs-terminal p {
  font-family: var(--font-mono);
  font-size: 11px;
  line-height: 1.5;
}

.log-info {
  color: var(--text-muted);
}

.log-warn {
  color: var(--color-warning);
}

@media (max-width: 1200px) {
  .status-card {
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }
}

@media (max-width: 900px) {
  .detail-header {
    flex-direction: column;
    align-items: flex-start;
  }

  .status-card {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .info-grid {
    grid-template-columns: 1fr;
  }

  .info-row:nth-last-child(-n + 2) {
    border-bottom: var(--border);
  }

  .info-row:last-child {
    border-bottom: none;
  }
}
</style>
