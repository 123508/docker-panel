<template>
  <dp-page-header title="创建容器" gap="32px">
    <template #actions>
      <dp-container-header-actions class="create-actions">
        <dp-button
          text="← 返回"
          type="default"
          variant="filled"
          size="medium"
          class="create-back-btn"
          @click="goBack"
        />
        <dp-button
          text="创建"
          type="primary"
          variant="filled"
          size="medium"
          class="create-submit-btn"
          :disabled="!canSubmit || state.creating"
          :loading="state.creating"
          @click="handleSubmit"
        />
      </dp-container-header-actions>
    </template>

    <dp-container-section-card class="create-card">
      <div class="form-content">
        <dp-container-form-field label="容器名称" required>
          <input v-model.trim="state.form.name" placeholder="my-container" />
        </dp-container-form-field>

        <dp-container-form-field label="镜像" required>
          <select v-model="state.form.image">
            <option v-for="img in state.imageOptions" :key="img" :value="img">{{ img }}</option>
          </select>
        </dp-container-form-field>

        <dp-container-form-field label="端口映射">
          <div class="section-head">
            <span></span>
            <dp-button text="+ 添加端口" type="default" variant="filled" size="small" @click="addPortRow" />
          </div>
          <div class="row-list">
            <div v-for="(row, idx) in state.form.ports" :key="`port-${idx}`" class="inline-row">
              <input v-model.trim="row.containerPort" placeholder="Container Port (e.g., 80)" />
              <input v-model.trim="row.hostPort" placeholder="Host Port (e.g., 8080)" />
              <button class="remove-btn" type="button" @click="removePortRow(idx)">×</button>
            </div>
          </div>
        </dp-container-form-field>

        <dp-container-form-field label="环境变量">
          <div class="section-head">
            <span></span>
            <dp-button text="+ 添加变量" type="default" variant="filled" size="small" @click="addEnvRow" />
          </div>
          <div class="row-list">
            <div v-for="(row, idx) in state.form.envs" :key="`env-${idx}`" class="inline-row">
              <input v-model.trim="row.key" placeholder="键" />
              <input v-model.trim="row.value" placeholder="值" />
              <button class="remove-btn" type="button" @click="removeEnvRow(idx)">×</button>
            </div>
          </div>
        </dp-container-form-field>

        <dp-container-form-field label="卷挂载">
          <div class="section-head">
            <span></span>
            <dp-button text="+ 添加卷" type="default" variant="filled" size="small" @click="addVolumeRow" />
          </div>
          <div class="row-list">
            <div v-for="(row, idx) in state.form.volumes" :key="`volume-${idx}`" class="inline-row">
              <input v-model.trim="row.source" placeholder="主机路径或卷名称" />
              <input v-model.trim="row.target" placeholder="容器路径" />
              <button class="remove-btn" type="button" @click="removeVolumeRow(idx)">×</button>
            </div>
          </div>
        </dp-container-form-field>

        <dp-container-form-field label="网络">
          <select v-model="state.form.networkMode">
            <option v-for="net in state.networkOptions" :key="net" :value="net">{{ net }}</option>
          </select>
        </dp-container-form-field>

        <dp-container-form-field label="重启策略">
          <div class="restart-options">
            <button
              v-for="item in restartOptions"
              :key="item"
              type="button"
              class="restart-btn"
              :class="{ active: state.form.restartPolicy === item }"
              @click="state.form.restartPolicy = item"
            >
              {{ item }}
            </button>
          </div>
        </dp-container-form-field>
      </div>
    </dp-container-section-card>
  </dp-page-header>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import DpPageHeader from '@/components/dp-page-header.vue'
import DpButton from '@/components/dp-button.vue'
import DpContainerHeaderActions from '@/components/container/dp-container-header-actions.vue'
import DpContainerSectionCard from '@/components/container/dp-container-section-card.vue'
import DpContainerFormField from '@/components/container/dp-container-form-field.vue'
import { ContainerCreateState } from '@/composables/ContainerCreate'

const restartOptions = ['no', 'always', 'on-failure', 'unless-stopped'] as const

const router = useRouter()
const {
  state,
  canSubmit,
  submit,
  addPortRow,
  removePortRow,
  addEnvRow,
  removeEnvRow,
  addVolumeRow,
  removeVolumeRow
} = ContainerCreateState()

const goBack = () => {
  router.push('/containers')
}

const handleSubmit = async () => {
  const id = await submit()
  if (id) {
    router.push(`/containers/${id}`)
  }
}
</script>

<style scoped>
.create-actions {
  gap: 12px;
  flex-wrap: nowrap;
}

.create-actions :deep(.create-back-btn.button-medium) {
  width: 92px;
}

.create-actions :deep(.create-submit-btn.button-medium) {
  width: 108px;
}

.create-card :deep(.card) {
  padding: 32px;
}

.form-content {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.form-content :deep(.field) {
  gap: 8px;
}

.form-content :deep(label) {
  color: var(--text-secondary);
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 1px;
}

.form-content :deep(.required) {
  margin-left: 6px;
}

input,
select,
textarea {
  width: 100%;
  border: var(--border);
  background: var(--create-input-bg);
  color: var(--text-primary);
  padding: 0 14px;
  font-size: 13px;
  outline: none;
  height: 40px;
}

.section-head {
  display: flex;
  justify-content: space-between;
  margin-bottom: 12px;
}

.section-head :deep(.button-small) {
  height: 30px;
  padding: 0 14px;
}

.row-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.inline-row {
  display: grid;
  grid-template-columns: 1fr 1fr 40px;
  gap: 10px;
}

.remove-btn {
  height: 40px;
  border: var(--border);
  background: var(--create-input-bg);
  color: var(--color-danger);
  font-size: 20px;
  cursor: pointer;
}

.restart-options {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 10px;
}

.restart-btn {
  height: 40px;
  border: none;
  background: var(--create-input-bg);
  color: var(--text-muted);
  text-transform: lowercase;
  font-family: var(--font-mono);
  font-size: 13px;
  cursor: pointer;
}

.restart-btn.active {
  background: var(--accent-soft);
  color: var(--accent);
  font-weight: 700;
}

.create-card {
  --create-input-bg: #141415;
  --accent-soft: #facc1520;
}

:global(html.light) .create-card {
  --create-input-bg: #e8f1ff;
  --accent-soft: rgba(14, 165, 233, 0.15);
}

@media (max-width: 1024px) {
  .inline-row {
    grid-template-columns: 1fr;
  }

  .restart-options {
    grid-template-columns: 1fr 1fr;
  }
}
</style>
