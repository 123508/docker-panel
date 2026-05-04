<template>
  <dp-dialog
    v-model="state.visible"
    :title="state.title"
    :content="state.content"
    :okText="state.okText"
    :closeOnOverlayClick="!state.isRunning"
    :showCancel="false"
    :loading="state.isRunning"
    @ok="state.visible = false"
    @cancel="state.visible = false"
  >
    <div v-if="state.isRunning" class="running-status">
      <svg class="spinner" viewBox="0 0 50 50">
        <circle class="path" cx="25" cy="25" r="20" fill="none" stroke-width="5"></circle>
      </svg>
      <p>{{ state.content }}</p>
    </div>
    <p v-else class="custom-dialog-content-text">{{ state.content }}</p>
  </dp-dialog>
</template>

<script setup lang="ts">
import DpDialog from '@/components/dp-dialog.vue'
import type { ActionDialogState } from '@/composables/useActionDialog'

defineOptions({ name: 'dp-action-dialog' })

defineProps<{
  state: ActionDialogState
}>()
</script>

<style scoped>
.custom-dialog-content-text {
  margin: 0;
  white-space: pre-wrap;
  word-wrap: break-word;
  word-break: break-all;
  color: var(--text-primary);
}

.running-status {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 16px;
  padding: 24px 0;
  color: var(--text-primary);
}

.spinner {
  animation: rotate 2s linear infinite;
  z-index: 2;
  width: 40px;
  height: 40px;
}

.spinner .path {
  stroke: var(--btn-primary-bg, #1890FF);
  stroke-linecap: round;
  animation: dash 1.5s ease-in-out infinite;
}

@keyframes rotate {
  100% {
    transform: rotate(360deg);
  }
}

@keyframes dash {
  0% {
    stroke-dasharray: 1, 150;
    stroke-dashoffset: 0;
  }
  50% {
    stroke-dasharray: 90, 150;
    stroke-dashoffset: -35;
  }
  100% {
    stroke-dasharray: 90, 150;
    stroke-dashoffset: -124;
  }
}
</style>
