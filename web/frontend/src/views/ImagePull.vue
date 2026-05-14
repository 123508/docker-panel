<template>
  <div class="pull-page">
    <div class="pull-header">
      <div class="title-wrap">
        <h1 class="pull-title">拉取镜像</h1>
      </div>
      <div class="pull-actions">
        <dp-button text="← 返回" size="medium" variant="outlined" @click="goBack" />
      </div>
    </div>

    <div class="search-area">
      <div class="search-wrap">
        <input
          class="search-input"
          v-model="state.query"
          placeholder="搜索镜像名称..."
          @keyup.enter="searchOnEnter"
        />
        <dp-button
          text="搜索"
          size="medium"
          type="primary"
          :loading="state.searching"
          @click="searchImages"
        />
      </div>
    </div>

    <div v-if="state.searching" class="placeholder">搜索中...</div>
    <div v-else-if="!state.searched" class="empty-state">
      <div class="empty-icon">09</div>
      <p class="empty-text">输入镜像名称开始搜索</p>
      <p class="empty-sub">支持搜索 Docker Hub 官方镜像，如 nginx、redis、postgres</p>
    </div>
    <div v-else-if="state.results.length === 0" class="empty-state">
      <p class="empty-text">未找到匹配的镜像</p>
      <p class="empty-sub">请尝试其他关键词</p>
    </div>
    <template v-else>
      <div class="meta-bar">
        <span class="meta-col col-name">镜像名称</span>
        <span class="meta-col col-desc">描述</span>
      </div>
      <div class="result-list">
        <div v-for="img in state.results" :key="img.name" class="result-row">
          <div class="result-left">
            <span class="img-name">{{ img.name }}</span>
            <span v-if="img.description" class="img-desc">{{ img.description }}</span>
          </div>
          <div class="result-right">
            <dp-button
              text="拉取"
              size="small"
              type="primary"
              :loading="state.pulling === img.name"
              :disabled="state.pulling !== '' && state.pulling !== img.name"
              @click="pullImage(img.name)"
            />
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import DpButton from '@/components/dp-button.vue'
import { ImagePullState } from '@/composables/ImagePull'

const router = useRouter()
const { state, searchImages, searchOnEnter, pullImage } = ImagePullState()

const goBack = () => router.push('/dashboard/images')
</script>

<style scoped>
.pull-page {
  padding: var(--page-padding-y) var(--page-padding-x);
  height: 100vh;
  display: flex;
  flex-direction: column;
  gap: 24px;
  overflow: auto;
}

.pull-header {
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

.pull-title {
  margin: 0;
  font-family: var(--font-mono);
  font-size: 32px;
  line-height: 1;
  font-weight: 700;
  color: var(--text-primary);
}

.pull-actions {
  display: flex;
  gap: 12px;
}

.search-area {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.search-wrap {
  display: flex;
  gap: 12px;
  align-items: center;
}

.search-input {
  flex: 1;
  height: 44px;
  background: var(--button-default);
  border: none;
  padding: 0 16px;
  font-family: var(--font-mono);
  font-size: 12px;
  color: var(--text-primary);
  outline: none;
}

.search-input::placeholder {
  color: var(--text-tertiary);
}

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

.meta-bar {
  display: flex;
  align-items: center;
  gap: 16px;
  min-height: 44px;
  padding: 0 20px;
  background: var(--bg-card-header);
  border: var(--border);
}

.meta-col {
  font-family: var(--font-mono);
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 1px;
  color: var(--text-secondary);
}

.col-name { width: 220px; }
.col-desc { flex: 1; min-width: 0; }

.result-list {
  display: flex;
  flex-direction: column;
  background: var(--bg-card);
  border: var(--border);
  border-top: none;
}

.result-row {
  display: flex;
  align-items: center;
  gap: 16px;
  min-height: 56px;
  padding: 0 20px;
  border-top: var(--border);
  transition: background 0.15s;
}

.result-row:hover {
  background: var(--table-row-hover);
}

.result-left {
  width: 220px;
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.img-name {
  font-family: var(--font-mono);
  font-size: 12px;
  font-weight: 700;
  color: var(--text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.img-desc {
  font-family: var(--font-mono);
  font-size: 11px;
  color: var(--accent);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.result-right {
  flex: 1;
  display: flex;
  justify-content: flex-end;
}
</style>
