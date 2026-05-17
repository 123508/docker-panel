<template>
  <div class="detail-page">
    <div class="detail-header">
      <div class="title-wrap">
        <h1 class="detail-title">{{ networkName }}</h1>
        <p class="detail-subtitle">网络 ID: {{ shortId }}</p>
      </div>
      <div class="detail-actions">
        <dp-button text="← 返回" size="medium" variant="outlined" @click="goBack" />
        <dp-button text="删除" type="danger" size="medium" :disabled="state.loading" @click="confirmDelete" />
      </div>
    </div>

    <div v-if="state.loading" class="placeholder">加载中...</div>
    <div v-else-if="!state.detail" class="placeholder">未获取到网络详情</div>
    <template v-else>
      <section class="status-card">
        <div class="status-item">
          <span class="item-label">状态</span>
          <span class="status-value active">
            <i class="status-dot"></i>
            活跃
          </span>
        </div>
        <div class="status-item">
          <span class="item-label">创建时间</span>
          <span class="item-value">{{ created }}</span>
        </div>
        <div class="status-item">
          <span class="item-label">驱动</span>
          <span class="item-value">{{ driver }}</span>
        </div>
        <div class="status-item">
          <span class="item-label">作用域</span>
          <span class="item-value">{{ scope }}</span>
        </div>
        <div class="status-item">
          <span class="item-label">连接容器</span>
          <span class="item-value">{{ containers.length }}</span>
        </div>
      </section>

      <section class="section-block">
        <h2 class="section-title">网络信息</h2>
        <div class="info-grid-card">
          <div class="info-grid">
            <div class="info-row">
              <div class="row-label">网络名称</div>
              <div class="row-value strong">{{ networkName }}</div>
            </div>
            <div class="info-row">
              <div class="row-label">网络 ID</div>
              <div class="row-value muted mono">{{ fullId }}</div>
            </div>
            <div class="info-row">
              <div class="row-label">驱动</div>
              <div class="row-value muted">{{ driver }}</div>
            </div>
            <div class="info-row">
              <div class="row-label">子网</div>
              <div class="row-value muted">{{ ipamConfigs[0]?.subnet || '-' }}</div>
            </div>
            <div class="info-row">
              <div class="row-label">网关</div>
              <div class="row-value muted">{{ ipamConfigs[0]?.gateway || '-' }}</div>
            </div>
          </div>
        </div>
      </section>

      <section v-if="ipamConfigs.length" class="section-block">
        <h2 class="section-title">IPAM 配置</h2>
        <div class="table-card">
          <div class="mini-table">
            <div class="mini-th">
              <span>子网</span>
              <span>网关</span>
              <span>IP 范围</span>
            </div>
            <div v-for="(cfg, i) in ipamConfigs" :key="i" class="mini-tr">
              <span class="cell-val">{{ cfg.subnet }}</span>
              <span class="cell-val">{{ cfg.gateway }}</span>
              <span class="cell-val">{{ cfg.ipRange }}</span>
            </div>
          </div>
        </div>
      </section>

      <section v-if="labelsList.length" class="section-block">
        <h2 class="section-title">标签</h2>
        <div class="code-block">
          <div v-for="l in labelsList" :key="l" class="code-line">{{ l }}</div>
        </div>
      </section>

      <section v-if="optionsList.length" class="section-block">
        <h2 class="section-title">选项</h2>
        <div class="code-block">
          <div v-for="o in optionsList" :key="o" class="code-line">{{ o }}</div>
        </div>
      </section>

      <section class="section-block">
        <h2 class="section-title">连接容器</h2>
        <div class="table-card">
          <div v-if="containers.length" class="mini-table">
            <div class="mini-th">
              <span>容器名称</span>
              <span>IPv4</span>
              <span>MAC 地址</span>
            </div>
            <div v-for="c in containers" :key="c.id" class="mini-tr">
              <span class="cell-key">{{ c.name }}</span>
              <span class="cell-val">{{ c.ipv4 }}</span>
              <span class="cell-src muted">{{ c.mac }}</span>
            </div>
          </div>
          <div v-else class="empty-hint">暂无容器连接到此网络</div>
        </div>
      </section>

      <section class="section-block">
        <h2 class="section-title">网络信息</h2>
        <div class="info-grid-card">
          <div class="info-grid">
            <div class="info-row">
              <div class="row-label">IPv6</div>
              <div class="row-value muted">{{ enableIPv6 ? '已启用' : '未启用' }}</div>
            </div>
            <div class="info-row">
              <div class="row-label">内部网络</div>
              <div class="row-value muted">{{ internal ? '是' : '否' }}</div>
            </div>
            <div class="info-row">
              <div class="row-label">可附加</div>
              <div class="row-value muted">{{ attachable ? '是' : '否' }}</div>
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
import { NetworkDetailState } from '@/composables/NetworkDetail'

const route = useRoute()
const router = useRouter()

const networkId = computed(() => String(route.params.id || ''))

const {
  state, networkName, shortId, fullId, driver, scope, created,
  enableIPv6, internal, attachable, ingress,
  ipamDriver, ipamConfigs, containers, labelsList, optionsList,
  loadData, handleRemove
} = NetworkDetailState(() => networkId.value)

onMounted(() => loadData())

const goBack = () => router.push('/dashboard/networks')

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

.status-value {
  display: inline-flex;
  align-items: center;
  gap: 6px;
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

.code-block {
  background: var(--bg-card);
  border: var(--border);
  padding: 16px 20px;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.code-line {
  font-family: var(--font-mono);
  font-size: 12px;
  color: var(--text-muted);
  line-height: 1.6;
}

.empty-hint {
  padding: 20px;
  text-align: center;
  color: var(--text-secondary);
  font-family: var(--font-mono);
  font-size: 12px;
}
</style>
