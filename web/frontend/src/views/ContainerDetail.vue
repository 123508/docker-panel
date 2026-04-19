<template>
  <dp-page-header :title="`容器详情 - ${displayName}`" gap="20px">
    <template #actions>
      <container-header-actions>
        <dp-button text="返回列表" variant="outlined" @click="goBack" />
        <dp-button
          v-if="isRunning"
          text="停止"
          type="danger"
          :disabled="state.actionLoading"
          @click="stop"
        />
        <dp-button
          v-else
          text="启动"
          type="primary"
          :disabled="state.actionLoading"
          @click="start"
        />
        <dp-button text="重启" variant="outlined" :disabled="state.actionLoading" @click="restart" />
        <dp-button text="删除" type="danger" variant="text" :disabled="state.actionLoading" @click="handleRemove" />
      </container-header-actions>
    </template>

    <div v-if="state.loading" class="placeholder">加载中...</div>
    <div v-else-if="!state.detail" class="placeholder">未获取到容器详情</div>
    <div v-else class="detail-wrap">
      <container-section-card title="基础信息">
        <div class="grid">
          <div class="item"><span>ID</span><b>{{ state.detail.id }}</b></div>
          <div class="item"><span>名称</span><b>{{ displayName }}</b></div>
          <div class="item"><span>镜像</span><b>{{ state.detail.image || 'N/A' }}</b></div>
          <div class="item"><span>状态</span><b>{{ state.detail.state?.running ? '运行中' : '已停止' }}</b></div>
          <div class="item"><span>平台</span><b>{{ state.detail.platform || 'N/A' }}</b></div>
          <div class="item"><span>重启次数</span><b>{{ state.detail.restart_count ?? 0 }}</b></div>
        </div>
      </container-section-card>

      <container-section-card title="网络信息">
        <pre>{{ JSON.stringify(state.detail.network_settings?.networks || {}, null, 2) }}</pre>
      </container-section-card>

      <container-section-card title="挂载信息">
        <pre>{{ JSON.stringify(state.detail.mounts || [], null, 2) }}</pre>
      </container-section-card>
    </div>
  </dp-page-header>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import DpPageHeader from '@/components/dp-page-header.vue'
import DpButton from '@/components/dp-button.vue'
import ContainerHeaderActions from '@/components/container/ContainerHeaderActions.vue'
import ContainerSectionCard from '@/components/container/ContainerSectionCard.vue'
import { ContainerDetailState } from '@/composables/ContainerDetail'

const route = useRoute()
const router = useRouter()

const containerId = computed(() => String(route.params.id || ''))
const { state, isRunning, displayName, start, stop, restart, remove } = ContainerDetailState(() => containerId.value)

const goBack = () => {
  router.push('/containers')
}

const handleRemove = async () => {
  const ok = await remove()
  if (ok) {
    router.push('/containers')
  }
}
</script>

<style scoped>
.placeholder {
  padding: 24px;
  color: var(--text-muted);
  border: var(--border);
  background: var(--bg-card);
}

.detail-wrap {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.item span {
  font-size: 12px;
  color: var(--text-muted);
}

.item b {
  font-size: 13px;
  color: var(--text-primary);
  word-break: break-all;
}

pre {
  margin: 0;
  white-space: pre-wrap;
  word-break: break-all;
  font-size: 12px;
  color: var(--text-muted);
  background: var(--bg-body);
  padding: 12px;
  border: var(--border);
}
</style>
