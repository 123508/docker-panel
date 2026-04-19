<template>
  <dp-page-header title="镜像" gap="24px">
    <template #actions>
      <div class="header-actions">
        <dp-search-input v-model="state.query" placeholder="搜索镜像..." />
        <dp-button text="搜索" size="medium" type="primary" class="search-btn" variant="outlined" @click="searchImages" />
        <dp-button text="+ 拉取" size="medium" type="primary" class="create-btn" @click="pullImage" />
      </div>
    </template>

    <dp-stat-bar :items="state.stats" />

    <dp-data-table :bordered="false" header-bg row-gap fixed-row-height :columns="columns">
      <div v-for="img in pagedImages" :key="img.fullId" class="table-row">
        <div class="td" style="width: 40px; display: flex; align-items: center">
          <div class="row-icon" :style="{ background: img.color }"></div>
        </div>
        <span class="td col-repo">{{ img.repo }}</span>
        <span class="td col-tag td-muted">{{ img.tag }}</span>
        <span class="td col-id td-dim">{{ img.id }}</span>
        <span class="td col-size">{{ img.size }}</span>
        <span class="td col-created td-muted-sm">{{ img.created }}</span>
        <div class="td td-actions" style="flex: 1">
          <dp-button text="查看" size="small" variant="text" type="info" @click="inspectImage(img.fullId)" />
          <dp-button text="运行" size="small" variant="text" type="primary" @click="runImage(img.fullId, `${img.repo}:${img.tag}`)" />
          <dp-button text="移除" size="small" variant="text" type="danger" @click="removeImage(img.fullId)" />
        </div>
      </div>
    </dp-data-table>

    <dp-pagination
      :page="state.page"
      :total="state.images.length"
      :page-size="state.pageSize"
      @change="state.page = $event"
    />

    <dp-dialog
      v-model="state.runDialogVisible"
      :title="state.runDialogTitle"
      :content="state.runDialogContent"
      :okText="state.runDialogOkText"
      :cancelText="state.runDialogCancelText"
      :closeOnOverlayClick="!state.runDialogIsRunning"
      :showCancel="!state.runDialogIsRunning"
      :loading="state.runDialogIsRunning"
      @ok="state.runDialogVisible = false"
      @cancel="state.runDialogVisible = false"
    >
      <div v-if="state.runDialogIsRunning" class="running-status">
        <svg class="spinner" viewBox="0 0 50 50">
          <circle class="path" cx="25" cy="25" r="20" fill="none" stroke-width="5"></circle>
        </svg>
        <p>{{ state.runDialogContent }}</p>
      </div>
      <p v-else class="custom-dialog-content-text">{{ state.runDialogContent }}</p>
    </dp-dialog>
  </dp-page-header>
</template>

<script setup lang="ts">
import DpPageHeader from '@/components/dp-page-header.vue'
import DpStatBar from '@/components/dp-stat-bar.vue'
import DpDataTable from '@/components/dp-data-table.vue'
import DpButton from '@/components/dp-button.vue'
import DpSearchInput from '@/components/dp-search-input.vue'
import DpPagination from '@/components/dp-pagination.vue'
import DpDialog from '@/components/dp-dialog.vue'
import { ImageState } from '@/composables/Images'

const { state, pagedImages, searchImages, pullImage, inspectImage, runImage, removeImage } = ImageState()

const columns = [
  { key: 'icon', width: 40 },
  { key: 'repo', label: '仓库', width: 120 },
  { key: 'tag', label: '标签', width: 100 },
  { key: 'id', label: '镜像 ID', width: 120 },
  { key: 'size', label: '大小', width: 80 },
  { key: 'created', label: '创建时间', width: 80 },
  { key: 'actions', label: '操作', flex: 1 }
]
</script>

<style scoped>
.header-actions {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
}

.search-btn {
  color: var(--text-primary);
  background: var(--bg-card);
  border: var(--border);
  padding: 10px 14px;
}

.create-btn {
  width: 120px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.table-row {
  display: flex;
  align-items: center;
  gap: 16px;
  transition: background 0.15s;
}

.table-row:hover {
  background: var(--table-row-hover);
}

.td {
  font-size: 12px;
  color: var(--text-primary);
}

.col-repo { width: 120px; }
.col-tag { width: 100px; }
.col-id { width: 120px; }
.col-size { width: 80px; }
.col-created { width: 80px; }

.row-icon {
  width: 32px;
  height: 32px;
  flex-shrink: 0;
}

.td-muted {
  color: var(--text-muted);
}

.td-dim {
  color: var(--text-secondary);
  font-size: 11px;
}

.td-muted-sm {
  color: var(--text-muted);
  font-size: 11px;
}

.td-actions {
  display: flex;
  align-items: center;
  gap: 16px;
}

.custom-dialog-content-text {
  margin: 0;
  white-space: pre-wrap;
  word-wrap: break-word;
  word-break: break-all;
  color: var(--text-primary);
}

.running-status {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 16px;
  padding: 24px 0;
  color: var(--text-primary);
}

.spinner {
  animation: rotate 2s linear infinite;
  z-index: 2;
  width: 40px;
  height: 40px;
}

.spinner .path {
  stroke: var(--btn-primary-bg, #1890FF);
  stroke-linecap: round;
  animation: dash 1.5s ease-in-out infinite;
}

@keyframes rotate {
  100% {
    transform: rotate(360deg);
  }
}

@keyframes dash {
  0% {
    stroke-dasharray: 1, 150;
    stroke-dashoffset: 0;
  }
  50% {
    stroke-dasharray: 90, 150;
    stroke-dashoffset: -35;
  }
  100% {
    stroke-dasharray: 90, 150;
    stroke-dashoffset: -124;
  }
}
</style>
