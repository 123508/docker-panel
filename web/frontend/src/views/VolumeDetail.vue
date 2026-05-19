<template>
  <div class="detail-page">
    <div class="detail-header">
      <div class="title-wrap">
        <h1 class="detail-title">{{ volumeName || '卷详情' }}</h1>
        <p class="detail-subtitle">卷名: {{ volumeName }}</p>
      </div>
      <div class="detail-actions">
        <dp-button text="← 返回" size="medium" variant="outlined" class="btn-back" @click="goBack" />
        <dp-button text="删除" type="danger" size="medium" class="btn-delete" :disabled="state.loading" @click="confirmDelete" />
      </div>
    </div>

    <div v-if="state.loading" class="placeholder">加载中...</div>
    <div v-else-if="!state.detail" class="placeholder">未获取到卷详情</div>
    <template v-else>
      <section class="status-card">
        <div class="status-item">
          <span class="item-label">驱动</span>
          <span class="item-value">{{ driver }}</span>
        </div>
        <div class="status-item">
          <span class="item-label">创建时间</span>
          <span class="item-value">{{ createdAt }}</span>
        </div>
        <div class="status-item">
          <span class="item-label">大小</span>
          <span class="item-value">{{ size }}</span>
        </div>
        <div class="status-item">
          <span class="item-label">引用数</span>
          <span class="item-value">{{ refCount }}</span>
        </div>
        <div class="status-item">
          <span class="item-label">挂载容器</span>
          <span class="item-value">{{ containers.length }}</span>
        </div>
      </section>

      <section class="section-block">
        <h2 class="section-title">基本信息</h2>
        <div class="info-grid-card">
          <div class="info-grid">
            <div class="info-row">
              <div class="row-label">名称</div>
              <div class="row-value strong">{{ volumeName }}</div>
            </div>
            <div class="info-row">
              <div class="row-label">驱动</div>
              <div class="row-value muted">{{ driver }}</div>
            </div>
            <div class="info-row">
              <div class="row-label">挂载点</div>
              <div class="row-value muted">{{ mountpoint }}</div>
            </div>
            <div class="info-row">
              <div class="row-label">作用域</div>
              <div class="row-value muted">{{ scope }}</div>
            </div>
          </div>
        </div>
      </section>

      <section v-if="labelsList.length" class="section-block">
        <h2 class="section-title">标签</h2>
        <div class="table-card">
          <div class="mini-table">
            <div class="mini-th">
              <span>键</span>
              <span>值</span>
              <span>来源</span>
            </div>
            <div v-for="item in labelsList" :key="item.key" class="mini-tr">
              <span class="cell-key">{{ item.key }}</span>
              <span class="cell-val">{{ item.value }}</span>
              <span class="cell-src muted">user</span>
            </div>
          </div>
        </div>
      </section>

      <section v-if="optionsList.length" class="section-block">
        <h2 class="section-title">选项</h2>
        <div class="table-card">
          <div class="mini-table">
            <div class="mini-th">
              <span>键</span>
              <span>值</span>
            </div>
            <div v-for="item in optionsList" :key="item.key" class="mini-tr">
              <span class="cell-key">{{ item.key }}</span>
              <span class="cell-val">{{ item.value }}</span>
            </div>
          </div>
        </div>
      </section>

      <section class="section-block">
        <h2 class="section-title">挂载容器</h2>
        <div class="table-card">
          <div v-if="containers.length" class="mini-table">
            <div class="mini-th">
              <span>容器名称</span>
            </div>
            <div v-for="c in containers" :key="c" class="mini-tr">
              <span class="cell-key">{{ c }}</span>
            </div>
          </div>
          <div v-else class="empty-hint">暂无容器使用此卷</div>
        </div>
      </section>
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import DpButton from '@/components/dp-button.vue'
import { VolumeDetailState } from '@/composables/VolumeDetail'

const route = useRoute()
const router = useRouter()

const volumeNameParam = computed(() => {
  const n = route.params.name
  return n ? String(n) : ''
})

const { state, volumeName, driver, mountpoint, createdAt, scope, size, refCount, containers, labelsList, optionsList, loadData, handleRemove } = VolumeDetailState(() => volumeNameParam.value)

onMounted(() => loadData())

const goBack = () => router.push('/dashboard/volumes')

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

.item-value,
.status-value {
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

.row-value {
  min-width: 0;
  overflow-wrap: anywhere;
  color: var(--text-primary);
  font-family: var(--font-mono);
  font-size: 14px;
}

.strong { font-weight: 700; color: var(--text-primary); }
.muted { color: var(--text-muted); }

.mini-table {
  display: flex;
  flex-direction: column;
}

.mini-th,
.mini-tr {
  display: grid;
  grid-template-columns: 2fr 2fr 1fr;
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

.empty-hint {
  padding: 20px;
  text-align: center;
  color: var(--text-secondary);
  font-family: var(--font-mono);
  font-size: 12px;
}
</style>
