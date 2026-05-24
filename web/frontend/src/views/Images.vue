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
        <span class="td col-repo td-ellipsis" :title="img.repo">{{ img.repo }}</span>
        <span class="td col-tag td-muted td-ellipsis" :title="img.tag">{{ img.tag }}</span>
        <span class="td col-id td-dim td-ellipsis" :title="img.id">{{ img.id }}</span>
        <span class="td col-size td-ellipsis" :title="img.size">{{ img.size }}</span>
        <span class="td col-created td-muted-sm td-ellipsis" :title="img.created">{{ img.created }}</span>
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

    <dp-action-dialog :state="dialog" />
  </dp-page-header>
</template>

<script setup lang="ts">
import DpPageHeader from '@/components/dp-page-header.vue'
import DpStatBar from '@/components/dp-stat-bar.vue'
import DpDataTable from '@/components/dp-data-table.vue'
import DpButton from '@/components/dp-button.vue'
import DpSearchInput from '@/components/dp-search-input.vue'
import DpPagination from '@/components/dp-pagination.vue'
import DpActionDialog from '@/components/dp-action-dialog.vue'
import { useRouter } from 'vue-router'
import { ImageState } from '@/composables/Images'

const router = useRouter()
const { state, dialog, pagedImages, searchImages, runImage, removeImage } = ImageState()

const inspectImage = (id: string) => router.push(`/dashboard/images/${id}`)
const pullImage = () => router.push('/dashboard/images/pull')

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

/* 内容超长时截断显示省略号，鼠标悬停可见完整文本 */
.td-ellipsis {
  display: inline-block;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  vertical-align: middle;
}

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
</style>
