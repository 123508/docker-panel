<template>
  <div class="app">
    <header class="header">
      <h1>🐳 Docker 可视化面板</h1>
      <p class="subtitle">Docker Container Management Dashboard</p>
    </header>

    <main class="main-content">
      <div class="status-card">
        <h2>系统状态</h2>
        <div class="status-item">
          <span class="label">后端状态:</span>
          <span class="value" :class="{ 'status-ok': apiStatus === 'ok' }">
            {{ apiStatus === 'ok' ? '✓ 运行中' : '✗ 未连接' }}
          </span>
        </div>
        <div class="status-item">
          <span class="label">API消息:</span>
          <span class="value">{{ apiMessage }}</span>
        </div>
      </div>

      <div class="content-grid">
        <BaseCard
            title="容器列表"
            :items="containers">
          <template #item="{ item }">
            <div class="item-name">{{ item.name }}</div>
            <div class="item-status" :class="'status-' + item.status">{{ item.status }}</div>
          </template>

          <template #empty>
            <div>这里没有任何内容哦~</div>
          </template>
        </BaseCard>

        <BaseCard
            title="镜像列表"
            :items="images">
          <template #item="{ item }">
            <div class="item-name">{{ item.name }}</div>
            <div class="item-meta">{{ item.size }}</div>
          </template>

          <template #empty>
            <div>暂无数据</div>
          </template>
        </BaseCard>

      </div>

      <div class="info-card">
        <h3>📝 说明</h3>
        <p>这是 Docker 可视化面板的基础框架。</p>
        <ul>
          <li>前端：Vue 3 + Vite</li>
          <li>后端：Go + embed 静态文件</li>
          <li>API 端点：/api/health, /api/containers, /api/images</li>
          <li>当前显示的是模拟数据，实际功能待实现</li>
        </ul>
      </div>
    </main>

    <BaseButton text="test" size="small" type="primary"></BaseButton>

    <footer class="footer">
      <p>Docker Panel v0.1.0 | 基础框架</p>
    </footer>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import BaseCard from "@/components/BaseCard.vue";
import BaseButton from "@/components/BaseButton.vue";

const apiStatus = ref('unknown')
const apiMessage = ref('正在连接...')
const containers = ref([  ])
const images = ref([])

async function fetchData() {
  try {
    const healthRes = await fetch('/api/health')
    const healthData = await healthRes.json()
    apiStatus.value = healthData.status
    apiMessage.value = healthData.message

    const containersRes = await fetch('/api/containers')
    const containersData = await containersRes.json()
    containers.value = containersData.containers || []

    const imagesRes = await fetch('/api/images')
    const imagesData = await imagesRes.json()
    images.value = imagesData.images || []
  } catch (error) {
    apiStatus.value = 'error'
    apiMessage.value = '连接失败: ' + error.message
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.app {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
}

.header {
  background: rgba(255, 255, 255, 0.95);
  padding: 2rem;
  text-align: center;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.header h1 {
  margin: 0;
  font-size: 2.5rem;
  color: #333;
  font-weight: 700;
}

.subtitle {
  margin: 0.5rem 0 0;
  color: #666;
  font-size: 1rem;
}

.main-content {
  flex: 1;
  padding: 2rem;
  max-width: 1200px;
  width: 100%;
  margin: 0 auto;
}

.status-card {
  background: white;
  border-radius: 12px;
  padding: 1.5rem;
  margin-bottom: 2rem;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.status-card h2 {
  margin: 0 0 1rem;
  color: #333;
  font-size: 1.5rem;
}

.status-item {
  display: flex;
  gap: 1rem;
  margin-bottom: 0.75rem;
}

.status-item .label {
  font-weight: 600;
  color: #555;
}

.status-item .value {
  color: #888;
}

.status-ok {
  color: #10b981 !important;
  font-weight: 600;
}

.content-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 2rem;
  margin-bottom: 2rem;
}

.card {
  background: white;
  border-radius: 12px;
  padding: 1.5rem;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  transition: transform 0.2s;
}

.card:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 12px rgba(0, 0, 0, 0.15);
}

.card h2 {
  margin: 0 0 1rem;
  color: #333;
  font-size: 1.25rem;
  border-bottom: 2px solid #667eea;
  padding-bottom: 0.5rem;
}

.list {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.list-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.75rem;
  background: #f9fafb;
  border-radius: 8px;
  border-left: 3px solid #667eea;
}

.item-name {
  font-weight: 500;
  color: #333;
}

.item-status {
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  font-size: 0.875rem;
  font-weight: 600;
  text-transform: capitalize;
}

.status-running {
  background: #d1fae5;
  color: #065f46;
}

.status-stopped {
  background: #fee2e2;
  color: #991b1b;
}

.item-meta {
  color: #666;
  font-size: 0.875rem;
}

.empty {
  text-align: center;
  color: #999;
  padding: 2rem;
  font-style: italic;
}

.info-card {
  background: white;
  border-radius: 12px;
  padding: 1.5rem;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.info-card h3 {
  margin: 0 0 1rem;
  color: #333;
}

.info-card p {
  color: #555;
  line-height: 1.6;
  margin: 0 0 0.5rem;
}

.info-card ul {
  margin: 0.5rem 0 0 1.5rem;
  color: #666;
}

.info-card li {
  margin-bottom: 0.25rem;
}

.footer {
  background: rgba(255, 255, 255, 0.95);
  padding: 1rem;
  text-align: center;
  color: #666;
  margin-top: auto;
}

.footer p {
  margin: 0;
}
</style>
