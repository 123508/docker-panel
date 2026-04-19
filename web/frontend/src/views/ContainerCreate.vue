<template>
  <dp-page-header title="创建容器" gap="20px">
    <template #actions>
      <container-header-actions>
        <dp-button text="返回列表" variant="outlined" @click="goBack" />
        <dp-button text="创建" type="primary" :disabled="!canSubmit || state.creating" @click="handleSubmit" />
      </container-header-actions>
    </template>

    <container-section-card>
      <div class="form-content">
        <container-form-field label="容器名称" required>
          <input v-model.trim="state.form.name" placeholder="例如: nginx-web-01" />
        </container-form-field>

        <container-form-field label="镜像" required>
          <select v-model="state.form.image">
            <option v-for="img in state.imageOptions" :key="img" :value="img">{{ img }}</option>
          </select>
        </container-form-field>

        <container-form-field label="启动命令">
          <input v-model.trim="state.form.cmd" placeholder="例如: nginx -g daemon off;" />
        </container-form-field>

        <container-form-field label="环境变量 (每行一个 KEY=VALUE)">
          <textarea v-model="state.form.envText" rows="5" placeholder="APP_ENV=prod&#10;TZ=Asia/Shanghai" />
        </container-form-field>

        <div class="row">
          <container-form-field label="网络模式">
            <select v-model="state.form.networkMode">
              <option v-for="net in state.networkOptions" :key="net" :value="net">{{ net }}</option>
            </select>
          </container-form-field>
          <container-form-field label="重启策略">
            <select v-model="state.form.restartPolicy">
              <option value="no">no</option>
              <option value="always">always</option>
              <option value="on-failure">on-failure</option>
              <option value="unless-stopped">unless-stopped</option>
            </select>
          </container-form-field>
        </div>

        <div class="row">
          <container-form-field label="宿主端口">
            <input v-model.trim="state.form.portHost" placeholder="例如: 8080" />
          </container-form-field>
          <container-form-field label="容器端口">
            <input v-model.trim="state.form.portContainer" placeholder="例如: 80" />
          </container-form-field>
        </div>
      </div>
    </container-section-card>
  </dp-page-header>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import DpPageHeader from '@/components/dp-page-header.vue'
import DpButton from '@/components/dp-button.vue'
import ContainerHeaderActions from '@/components/container/ContainerHeaderActions.vue'
import ContainerSectionCard from '@/components/container/ContainerSectionCard.vue'
import ContainerFormField from '@/components/container/ContainerFormField.vue'
import { ContainerCreateState } from '@/composables/ContainerCreate'

const router = useRouter()
const { state, canSubmit, submit } = ContainerCreateState()

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
.form-content {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

input,
select,
textarea {
  width: 100%;
  border: var(--border);
  background: var(--bg-body);
  color: var(--text-primary);
  padding: 10px 12px;
  font-size: 13px;
  outline: none;
}

textarea {
  resize: vertical;
  min-height: 90px;
}

@media (max-width: 900px) {
  .row {
    grid-template-columns: 1fr;
  }
}
</style>
