<template>
  <div class="create-page">
    <div class="create-header">
      <div class="title-wrap">
        <h1 class="create-title">创建编排</h1>
      </div>
      <div class="create-actions">
        <dp-button text="← 返回" size="medium" variant="outlined" @click="goBack" />
        <dp-button text="创建" size="medium" type="primary" :loading="submitting" @click="handleCreate" />
      </div>
    </div>

    <div class="form-card">
      <div class="field">
        <label class="field-label">编排名称</label>
        <input class="field-input" v-model="form.name" placeholder="my-compose" />
      </div>

      <div class="field">
        <label class="field-label">编排内容 (YAML)</label>
        <textarea
          class="field-textarea"
          v-model="form.content"
          placeholder="version: '3.8'&#10;services:&#10;  web:&#10;    image: nginx:latest&#10;    ports:&#10;      - &quot;8080:80&quot;"
        ></textarea>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { uploadComposeProject } from '@/services/modules/compose'
import { ElMessage } from 'element-plus'
import DpButton from '@/components/dp-button.vue'

const router = useRouter()
const submitting = ref(false)

const form = reactive({
  name: '',
  content: `version: '3.8'
services:
  web:
    image: nginx:latest
    ports:
      - "8080:80"
`
})

const goBack = () => router.push('/dashboard/compose')

const handleCreate = async () => {
  if (!form.name.trim()) {
    ElMessage.warning('请输入编排名称')
    return
  }
  if (!form.content.trim()) {
    ElMessage.warning('请输入编排内容')
    return
  }
  submitting.value = true
  try {
    await uploadComposeProject({ name: form.name.trim(), content: form.content })
    ElMessage.success(`编排项目 ${form.name} 创建成功`)
    router.push('/dashboard/compose')
  } catch (e: any) {
    ElMessage.error(e?.message || '创建编排失败')
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

.field-textarea {
  min-height: 420px;
  background: var(--button-default);
  border: none;
  padding: 16px;
  font-family: var(--font-mono);
  font-size: 12px;
  color: var(--text-primary);
  outline: none;
  resize: vertical;
  line-height: 1.6;
}

.field-textarea::placeholder {
  color: var(--text-tertiary);
}
</style>
