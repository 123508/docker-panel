<template>
  <div class="create-page">
    <div class="create-header">
      <div class="title-wrap">
        <h1 class="create-title">创建网络</h1>
      </div>
      <div class="create-actions">
        <dp-button text="← 返回" size="medium" variant="outlined" @click="goBack" />
        <dp-button text="创建" size="medium" type="primary" :loading="submitting" @click="handleCreate" />
      </div>
    </div>

    <div class="form-card">
      <div class="field">
        <label class="field-label">网络名称</label>
        <input class="field-input" v-model="form.name" placeholder="my-network" />
      </div>

      <div class="field">
        <label class="field-label">驱动</label>
        <div class="select-group">
          <button
            v-for="d in drivers"
            :key="d"
            class="select-opt"
            :class="{ active: form.driver === d }"
            @click="form.driver = d"
          >{{ d }}</button>
        </div>
      </div>

      <div class="field">
        <label class="field-label">子网</label>
        <input class="field-input" v-model="form.subnet" placeholder="172.18.0.0/16" />
      </div>

      <div class="field">
        <label class="field-label">网关</label>
        <input class="field-input" v-model="form.gateway" placeholder="172.18.0.1" />
      </div>

      <div class="field">
        <label class="field-label">IP 范围</label>
        <input class="field-input" v-model="form.ipRange" placeholder="172.18.0.0/24" />
      </div>

      <div class="field toggles">
        <label class="field-label">选项</label>
        <label class="toggle-row">
          <input type="checkbox" v-model="form.internal" />
          <span>内部网络（禁止外部访问）</span>
        </label>
        <label class="toggle-row">
          <input type="checkbox" v-model="form.attachable" />
          <span>允许独立容器附加</span>
        </label>
        <label class="toggle-row">
          <input type="checkbox" v-model="form.enableIPv6" />
          <span>启用 IPv6</span>
        </label>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { createNetwork as apiCreateNetwork } from '@/services/modules/network'
import { ElMessage } from 'element-plus'
import DpButton from '@/components/dp-button.vue'

const router = useRouter()
const submitting = ref(false)

const drivers = ['bridge', 'overlay', 'macvlan', 'ipvlan']

const form = reactive({
  name: '',
  driver: 'bridge',
  subnet: '',
  gateway: '',
  ipRange: '',
  internal: false,
  attachable: false,
  enableIPv6: false
})

const goBack = () => router.push('/dashboard/networks')

const handleCreate = async () => {
  if (!form.name.trim()) {
    ElMessage.warning('请输入网络名称')
    return
  }
  submitting.value = true
  try {
    const ipam: any = { driver: 'default', config: [] }
    const config: any = {}
    if (form.subnet) config.subnet = form.subnet
    if (form.gateway) config.gateway = form.gateway
    if (form.ipRange) config.ip_range = form.ipRange
    if (Object.keys(config).length) ipam.config.push(config)

    await apiCreateNetwork({
      name: form.name.trim(),
      driver: form.driver,
      internal: form.internal || undefined,
      attachable: form.attachable || undefined,
      enable_ipv6: form.enableIPv6 || undefined,
      ipam: ipam.config.length > 0 ? ipam : undefined
    })
    ElMessage.success(`网络创建成功: ${form.name}`)
    router.push('/dashboard/networks')
  } catch (e: any) {
    ElMessage.error(e.message || '创建网络失败')
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.create-page {
  padding: var(--page-padding-y) var(--page-padding-x);
  height: 100vh;
  display: flex;
  flex-direction: column;
  gap: 32px;
  overflow: auto;
}

.create-header {
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

.create-title {
  margin: 0;
  font-family: var(--font-display);
  font-size: 40px;
  line-height: 1;
  font-weight: 700;
  letter-spacing: -1px;
  color: var(--text-primary);
}

.create-actions {
  display: flex;
  gap: 12px;
}

.form-card {
  background: var(--bg-card);
  border: var(--border);
  padding: 32px;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.field-label {
  font-family: var(--font-mono);
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 1px;
  color: var(--text-secondary);
}

.field-input {
  height: 44px;
  background: var(--button-default);
  border: none;
  padding: 0 16px;
  font-family: var(--font-mono);
  font-size: 12px;
  color: var(--text-primary);
  outline: none;
}

.field-input::placeholder {
  color: var(--text-tertiary);
}

.select-group {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.select-opt {
  height: 44px;
  padding: 0 20px;
  background: var(--button-default);
  border: none;
  font-family: var(--font-mono);
  font-size: 12px;
  color: var(--text-muted);
  cursor: pointer;
}

.select-opt.active {
  background: var(--accent);
  color: var(--toggle-active-icon);
  font-weight: 700;
}

.toggles {
  gap: 12px;
}

.toggle-row {
  display: flex;
  align-items: center;
  gap: 10px;
  font-family: var(--font-mono);
  font-size: 12px;
  color: var(--text-primary);
  cursor: pointer;
}

.toggle-row input[type="checkbox"] {
  accent-color: var(--accent);
}
</style>
