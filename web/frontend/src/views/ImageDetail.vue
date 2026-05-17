<template>
  <div class="detail-page">
    <div class="detail-header">
      <div class="title-wrap">
        <h1 class="detail-title">{{ imageName }}</h1>
        <p class="detail-subtitle">镜像 ID: {{ shortId }}</p>
      </div>
      <div class="detail-actions">
        <dp-button text="← 返回" size="medium" variant="outlined" @click="goBack" />
        <dp-button text="删除" type="danger" size="medium" :disabled="state.loading" @click="confirmDelete" />
      </div>
    </div>

    <div v-if="state.loading" class="placeholder">加载中...</div>
    <div v-else-if="!state.detail" class="placeholder">未获取到镜像详情</div>
    <template v-else>
      <section class="status-card">
        <div class="status-item">
          <span class="item-label">大小</span>
          <span class="item-value">{{ size }}</span>
        </div>
        <div class="status-item">
          <span class="item-label">层数</span>
          <span class="item-value">{{ layerCount }}</span>
        </div>
        <div class="status-item">
          <span class="item-label">架构</span>
          <span class="item-value">{{ architecture }}</span>
        </div>
        <div class="status-item">
          <span class="item-label">操作系统</span>
          <span class="item-value">{{ os }}</span>
        </div>
        <div class="status-item">
          <span class="item-label">容器数</span>
          <span class="item-value">{{ containers }}</span>
        </div>
      </section>

      <section class="section-block">
        <h2 class="section-title">镜像信息</h2>
        <div class="info-grid-card">
          <div class="info-grid">
            <div class="info-row">
              <div class="row-label">仓库</div>
              <div class="row-value strong">{{ repo }}</div>
            </div>
            <div class="info-row">
              <div class="row-label">标签</div>
              <div class="row-value muted">{{ tag }}</div>
            </div>
            <div class="info-row">
              <div class="row-label">创建时间</div>
              <div class="row-value muted">{{ created }}</div>
            </div>
            <div class="info-row">
              <div class="row-label">Docker 版本</div>
              <div class="row-value muted">{{ dockerVersion }}</div>
            </div>
          </div>
        </div>
      </section>

      <section v-if="exposedPorts.length" class="section-block">
        <h2 class="section-title">暴露端口</h2>
        <div class="table-card">
          <div class="mini-table">
            <div class="mini-th">
              <span>端口</span>
              <span>协议</span>
              <span>用途</span>
            </div>
            <div v-for="p in exposedPorts" :key="p.port + p.protocol" class="mini-tr">
              <span class="cell-key">{{ p.port }}</span>
              <span class="cell-val">{{ p.protocol }}</span>
              <span class="cell-src muted">-</span>
            </div>
          </div>
        </div>
      </section>

      <section v-if="envVars.length" class="section-block">
        <h2 class="section-title">环境变量</h2>
        <div class="table-card">
          <div class="info-grid info-grid-compact">
            <div v-for="env in envVars" :key="env.key" class="info-row-info">
              <div class="row-label-small">{{ env.key }}</div>
              <div class="row-value muted">{{ env.value }}</div>
            </div>
          </div>
        </div>
      </section>

      <section class="section-block">
        <h2 class="section-title">摘要</h2>
        <div class="table-card">
          <div class="info-grid info-grid-compact">
            <div class="info-row-info">
              <div class="row-label-small">镜像 ID</div>
              <div class="row-value muted mono">{{ fullId }}</div>
            </div>
            <div class="info-row-info">
              <div class="row-label-small">摘要 (Digest)</div>
              <div class="row-value muted mono">{{ digest }}</div>
            </div>
          </div>
        </div>
      </section>
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import DpButton from '@/components/dp-button.vue'
import { ImageDetailState } from '@/composables/ImageDetail'

const route = useRoute()
const router = useRouter()

const imageId = computed(() => String(route.params.id || ''))

const { state, imageName, shortId, fullId, repo, tag, created, size, os, architecture, dockerVersion, layerCount, containers, envVars, exposedPorts, digest, loadData, handleRemove } = ImageDetailState(() => imageId.value)

onMounted(() => loadData())

const goBack = () => router.push('/dashboard/images')

const confirmDelete = async () => {
  const ok = await handleRemove()
  if (ok) goBack()
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
  display: flex;
  gap: 12px;
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

.item-value {
  min-width: 0;
  overflow-wrap: anywhere;
  color: var(--text-primary);
  font-family: var(--font-mono);
  font-size: 14px;
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
.table-card {
  background: var(--bg-card);
  border: var(--border);
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

.info-grid-compact .info-row-info {
  display: grid;
  grid-template-columns: 240px minmax(0, 1fr);
  padding: 12px 20px;
  border-bottom: var(--border);
}

.info-grid-compact .info-row-info:last-child {
  border-bottom: none;
}

.row-label-small {
  font-family: var(--font-mono);
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 1px;
  color: var(--text-secondary);
}

.row-value {
  min-width: 0;
  overflow-wrap: anywhere;
  color: var(--text-primary);
  font-family: var(--font-mono);
  font-size: 14px;
}

.mono { font-size: 11px; }

.strong { font-weight: 700; color: var(--text-primary); }
.muted { color: var(--text-muted); }

.mini-table {
  display: flex;
  flex-direction: column;
}

.mini-th,
.mini-tr {
  display: grid;
  grid-template-columns: 2fr 1fr 2fr;
  gap: 12px;
  padding: 0 20px;
  min-height: 36px;
  align-items: center;
  font-size: 12px;
  border-bottom: var(--border);
}

.mini-th {
  min-height: 32px;
  color: var(--text-secondary);
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 1px;
}

.mini-tr:last-child {
  border-bottom: none;
}

.cell-key { color: var(--text-primary); font-weight: 700; }
.cell-val { color: var(--text-muted); }
.cell-src { font-size: 11px; }
</style>
