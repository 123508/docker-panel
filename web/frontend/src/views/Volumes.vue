<template>
  <dp-page-header title="卷" gap="24px">
    <template #actions>
      <div class="header-actions">
        <dp-button text="+ 创建" size="large" type="primary" class="btn" />
      </div>
    </template>

    <dp-stat-bar :items="form.stats" />

    <dp-data-table :bordered="false" header-bg row-gap fixed-row-height :columns="columns">

      <div v-for="v in pagedVolumes" :key="v.name" class="table-row">
        <span class="td col-name">{{ v.name }}</span>
        <span class="td col-driver td-muted">{{ v.driver }}</span>
        <span class="td col-mount td-dim">{{ v.mountpoint }}</span>
        <span class="td col-size">{{ v.size }}</span>
        <div class="td col-containers">
          <dp-count-badge :count="v.containers" />
        </div>
        <div class="td td-actions" style="flex: 1">
          <dp-button text="查看" size="small" variant="text" type="info"/>
          <dp-button text="移除" size="small" variant="text" type="danger"/>
        </div>
      </div>
    </dp-data-table>
    <dp-pagination
        :page="form.page"
        :total="form.volumes.length"
        :page-size="form.pageSize"
        @change="form.page = $event"
    />
  </dp-page-header>
</template>

<script setup lang="ts">
import DpPageHeader from '@/components/dp-page-header.vue'
import DpStatBar from '@/components/dp-stat-bar.vue'
import DpDataTable from '@/components/dp-data-table.vue'
import DpCountBadge from '@/components/dp-count-badge.vue'
import DpButton from "@/components/dp-button.vue";
import DpPagination from "@/components/dp-pagination.vue";
import {computed, reactive} from "vue";

const columns = [
  { key: 'name', label: '名称', width: 140 },
  { key: 'driver', label: '驱动', width: 80 },
  { key: 'mount', label: '挂载点', width: 320 },
  { key: 'size', label: '大小', width: 80 },
  { key: 'containers', label: '容器', width: 100 },
  { key: 'actions', label: '操作', flex: 1 }
]

const form = reactive({
  page:1,
  pageSize:5,

  stats:[
    { label: '总数:', value: '8' },
    { label: '已使用:', value: '3.2 GB' },
    { label: '未使用:', value: '2', variant: 'dangling' }
  ],

  volumes:[
    { name: 'postgres_data', driver: 'local', mountpoint: '/var/lib/docker/volumes/postgres_data/_data', size: '1.2 GB', containers: 1 },
    { name: 'redis_cache', driver: 'local', mountpoint: '/var/lib/docker/volumes/redis_cache/_data', size: '512 MB', containers: 1 },
    { name: 'app_logs', driver: 'local', mountpoint: '/var/lib/docker/volumes/app_logs/_data', size: '256 MB', containers: 0 }
  ],

})

const pagedVolumes = computed(() => {
  const start = (form.page - 1) * form.pageSize
  return form.volumes.slice(start, start + form.pageSize)
})


</script>

<style scoped>

.btn{
  width: 140px;
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

.td-actions {
  display: flex;
  align-items: center;
  gap: 16px;
}

</style>
