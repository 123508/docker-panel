<template>
  <div class="create-page">
    <div class="create-header">
      <div class="title-wrap">
        <h1 class="create-title">创建卷</h1>
      </div>
      <div class="create-actions">
        <dp-button text="← 返回" size="medium" variant="outlined" @click="goBack" />
        <dp-button text="创建" size="medium" type="primary" :loading="submitting" @click="handleCreate" />
      </div>
    </div>

    <div class="form-card">
      <div class="field">
        <label class="field-label">卷名称</label>
        <input class="field-input" v-model="form.name" placeholder="my-volume" />
      </div>

      <div class="field">
        <label class="field-label">驱动</label>
        <input class="field-input" v-model="form.driver" placeholder="local" />
      </div>

      <div class="section-block">
        <div class="section-header">
          <label class="field-label">驱动选项</label>
          <dp-button text="+ 添加" size="small" variant="outlined" @click="addDriverOpt" />
        </div>
        <div v-for="(opt, i) in form.driverOpts" :key="i" class="row-group">
          <input class="field-input flex-1" v-model="opt.key" placeholder="键" />
          <input class="field-input flex-1" v-model="opt.value" placeholder="值" />
          <dp-button text="×" size="small" type="danger" variant="text" @click="form.driverOpts.splice(i, 1)" />
        </div>
      </div>

      <div class="section-block">
        <div class="section-header">
          <label class="field-label">标签</label>
          <dp-button text="+ 添加" size="small" variant="outlined" @click="addLabel" />
        </div>
        <div v-for="(label, i) in form.labels" :key="i" class="row-group">
          <input class="field-input flex-1" v-model="label.key" placeholder="键" />
          <input class="field-input flex-1" v-model="label.value" placeholder="值" />
          <dp-button text="×" size="small" type="danger" variant="text" @click="form.labels.splice(i, 1)" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { createVolume as apiCreateVolume } from '@/services/modules/volume'
import { ElMessage } from 'element-plus'
import DpButton from '@/components/dp-button.vue'

const router = useRouter()
const submitting = ref(false)

const form = reactive({
  name: '',
  driver: 'local',
  driverOpts: [] as { key: string; value: string }[],
  labels: [] as { key: string; value: string }[]
})

const addDriverOpt = () => form.driverOpts.push({ key: '', value: '' })
const addLabel = () => form.labels.push({ key: '', value: '' })

const goBack = () => router.push('/dashboard/volumes')

const handleCreate = async () => {
  if (!form.name.trim()) {
    ElMessage.warning('请输入卷名称')
    return
  }
  submitting.value = true
  try {
    const driver_opts: Record<string, string> = {}
    form.driverOpts.forEach(opt => {
      if (opt.key) driver_opts[opt.key] = opt.value
    })
    const labels: Record<string, string> = {}
    form.labels.forEach(l => {
      if (l.key) labels[l.key] = l.value
    })
    await apiCreateVolume({
      name: form.name.trim(),
      driver: form.driver || 'local',
      driver_opts: Object.keys(driver_opts).length ? driver_opts : undefined,
      labels: Object.keys(labels).length ? labels : undefined
    })
    ElMessage.success(`卷创建成功: ${form.name}`)
    router.push('/dashboard/volumes')
  } catch (e: any) {
    ElMessage.error(e.message || '创建卷失败')
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

.section-block {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.row-group {
  display: flex;
  gap: 12px;
  align-items: center;
}

.flex-1 {
  flex: 1;
}
</style>
