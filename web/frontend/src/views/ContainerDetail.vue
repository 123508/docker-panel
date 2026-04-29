<template>
  <div class="detail-page">
    <div class="detail-header">
      <div class="title-wrap">
        <h1 class="detail-title">{{ displayName }}</h1>
        <p class="detail-subtitle">容器 ID: {{ shortId }}</p>
      </div>
      <dp-container-header-actions class="detail-actions">
        <dp-button
          text="← 返回"
          size="medium"
          variant="outlined"
          class="btn-back"
          :disabled="state.actionLoading"
          @click="goBack"
        />
        <dp-button
          v-if="isRunning"
          text="停止"
          type="danger"
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
      </dp-container-header-actions>
    </div>

    <div v-if="state.loading" class="placeholder">加载中...</div>
    <div v-else-if="!state.detail" class="placeholder">未获取到容器详情</div>
    <template v-else>
      <section class="status-card">
        <div class="status-item">
          <span class="item-label">状态</span>
          <span class="status-value" :class="{ stopped: !isRunning }">
            <i class="status-dot"></i>
            {{ isRunning ? '运行中' : '已停止' }}
          </span>
        </div>
        <div class="status-item">
          <span class="item-label">运行时长</span>
          <span class="item-value">{{ uptimeText }}</span>
        </div>
        <div class="status-item">
          <span class="item-label">CPU</span>
          <span class="item-value">{{ cpuText }}</span>
        </div>
        <div class="status-item">
          <span class="item-label">内存</span>
          <span class="item-value">{{ memoryText }}</span>
        </div>
        <div class="status-item">
          <span class="item-label">网络</span>
          <span class="item-value">{{ networkSummary }}</span>
        </div>
      </section>

      <section class="section-block">
        <h2 class="section-title">容器信息</h2>
        <dp-container-section-card class="info-grid-card">
          <div class="info-grid">
            <div class="info-row">
              <div class="row-label">镜像</div>
              <div class="row-value strong">{{ imageText }}</div>
            </div>
            <div class="info-row">
              <div class="row-label">命令</div>
              <div class="row-value muted">{{ commandText }}</div>
            </div>
            <div class="info-row">
              <div class="row-label">端口</div>
              <div class="row-value">
                <div v-if="portRows.length" class="ports-table">
                  <div class="ports-row ports-head">
                    <span>容器端口</span>
                    <span>主机端口</span>
                    <span>协议</span>
                  </div>
                  <div v-for="port in portRows" :key="port.key" class="ports-row">
                    <span>{{ port.private }}</span>
                    <span>{{ port.public }}</span>
                    <span>{{ port.type }}</span>
                  </div>
                </div>
                <span v-else class="empty-value">N/A</span>
              </div>
            </div>
            <div class="info-row">
              <div class="row-label">环境变量</div>
              <div class="row-value list-value">
                <div v-for="item in envRows" :key="item.raw" class="kv-line">
                  <strong>{{ item.key }}</strong><span>{{ item.value }}</span>
                </div>
                <span v-if="!envRows.length" class="empty-value">N/A</span>
              </div>
            </div>
            <div class="info-row">
              <div class="row-label">卷</div>
              <div class="row-value list-value">
                <div v-for="mount in mountRows" :key="mount.key" class="mount-line">
                  <span class="path-text">{{ mount.source }}</span>
                  <span class="arrow">→</span>
                  <span class="muted">{{ mount.destination }}</span>
                  <span class="mode">{{ mount.mode }}</span>
                </div>
                <span v-if="!mountRows.length" class="empty-value">N/A</span>
              </div>
            </div>
            <div class="info-row">
              <div class="row-label">网络</div>
              <div class="row-value list-value">
                <div v-for="network in networkRows" :key="network.name" class="network-line">
                  <strong>{{ network.name }}</strong>
                  <span class="muted">{{ network.ip }}</span>
                  <span class="secondary">Gateway: {{ network.gateway }}</span>
                </div>
                <span v-if="!networkRows.length" class="empty-value">N/A</span>
              </div>
            </div>
            <div class="info-row">
              <div class="row-label">重启策略</div>
              <div class="row-value strong">{{ restartPolicyText }}</div>
            </div>
          </div>
        </dp-container-section-card>
      </section>

      <section class="section-block">
        <h2 class="section-title">容器日志</h2>
        <dp-container-section-card class="logs-card">
          <div class="logs-terminal">
            <p v-for="(line, idx) in logLines" :key="`${idx}-${line.text}`" :class="line.levelClass">
              {{ line.text }}
            </p>
          </div>
        </dp-container-section-card>
      </section>
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import DpButton from '@/components/dp-button.vue'
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

const {
  state,
  isRunning,
  displayName,
  shortId,
  uptimeText,
  imageText,
  commandText,
  cpuText,
  memoryText,
  portRows,
  envRows,
  mountRows,
  networkRows,
  networkSummary,
  restartPolicyText,
  logLines,
  start,
  stop,
  restart
} = ContainerDetailState(() => containerId.value)

const goBack = () => {
  router.push('/dashboard/containers')
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
  margin: 0;
  font-family: var(--font-display);
  font-size: 40px;
  line-height: 1;
  font-weight: 700;
  letter-spacing: -1px;
  color: var(--text-primary);
}

.detail-subtitle {
  margin: 0;
  font-family: var(--font-mono);
  font-size: 11px;
  color: var(--text-secondary);
}

.detail-actions {
  justify-content: flex-end;
  gap: 12px;
}

.detail-actions :deep(.button) {
  font-family: var(--font-mono);
}

.detail-actions :deep(.btn-back.button-outlined) {
  padding: 0 14px;
  background: var(--bg-card);
}

.detail-actions :deep(.btn-back.button-outlined:hover:not([data-disabled="true"])) {
  background: var(--bg-card-header);
}

.detail-actions :deep(.btn-stop.button-filled.button-danger),
.detail-actions :deep(.btn-start.button-filled.button-primary),
.detail-actions :deep(.btn-restart.button-filled.button-primary) {
  padding: 0 18px;
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
  min-width: 0;
}

.item-label,
.row-label {
  font-family: var(--font-mono);
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 1px;
  color: var(--text-secondary);
}

.item-value,
.status-value {
  min-width: 0;
  overflow-wrap: anywhere;
  color: var(--text-primary);
  font-family: var(--font-mono);
  font-size: 14px;
}

.status-value {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.status-value.stopped .status-dot {
  background: var(--color-danger);
}

.status-dot {
  width: 8px;
  height: 8px;
  background: var(--color-success);
}

.section-block {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.section-title {
  margin: 0;
  font-family: var(--font-display);
  font-size: 18px;
  line-height: 1.2;
  font-weight: 700;
  color: var(--text-primary);
}

.info-grid-card,
.logs-card {
  padding: 0;
}

.info-grid {
  display: flex;
  flex-direction: column;
}

.info-row {
  display: grid;
  grid-template-columns: 240px minmax(0, 1fr);
  gap: 0;
  padding: 20px;
  border-bottom: var(--border);
}

.info-row:last-child {
  border-bottom: none;
}

.row-value {
  min-width: 0;
  overflow-wrap: anywhere;
  color: var(--text-primary);
  font-family: var(--font-mono);
  font-size: 14px;
}

.strong,
.kv-line strong,
.network-line strong {
  font-weight: 700;
  color: var(--text-primary);
}

.muted,
.kv-line span {
  color: var(--text-muted);
}

.secondary,
.arrow,
.empty-value {
  color: var(--text-secondary);
}

.list-value {
  display: flex;
  flex-direction: column;
  gap: 8px;
  font-size: 12px;
}

.kv-line,
.mount-line,
.network-line {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
}

.mount-line,
.network-line {
  flex-wrap: wrap;
}

.path-text {
  color: var(--text-primary);
}

.mode {
  color: var(--color-info);
  font-size: 11px;
}

.ports-table {
  display: flex;
  flex-direction: column;
  background: var(--button-default);
  width: 100%;
}

.ports-row {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 12px;
  min-height: 36px;
  align-items: center;
  padding: 0 12px;
  font-size: 12px;
  color: var(--text-muted);
}

.ports-head {
  min-height: 32px;
  color: var(--text-secondary);
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 1px;
}

.logs-terminal {
  background: var(--bg-card);
  padding: 16px;
  height: 280px;
  display: flex;
  flex-direction: column;
  gap: 4px;
  overflow: auto;
}

.logs-terminal p {
  margin: 0;
  font-family: var(--font-mono);
  font-size: 11px;
  line-height: 1.5;
  white-space: pre-wrap;
  overflow-wrap: anywhere;
}

.log-info {
  color: var(--text-muted);
}

.log-warn {
  color: var(--accent);
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

  .info-row {
    grid-template-columns: 1fr;
    gap: 10px;
  }
}
</style>
