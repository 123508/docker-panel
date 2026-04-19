<template>
  <div v-if="modelValue" class="custom-dialog-overlay" @click="handleOverlayClick">
    <div class="custom-dialog-container" @click.stop>
      <!-- Header -->
      <div class="custom-dialog-header">
        <span class="custom-dialog-title">{{ title }}</span>
        <button class="custom-dialog-close" @click="close" aria-label="关闭">
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <line x1="18" y1="6" x2="6" y2="18"></line>
            <line x1="6" y1="6" x2="18" y2="18"></line>
          </svg>
        </button>
      </div>

      <!-- Body -->
      <div class="custom-dialog-body">
        <slot>
          <p class="custom-dialog-content-text">{{ content }}</p>
        </slot>
      </div>

      <!-- Footer -->
      <div class="custom-dialog-footer" v-if="showFooter">
        <button v-if="showCancel" class="custom-dialog-btn custom-dialog-btn-cancel" @click="handleCancel" :disabled="loading">
          {{ cancelText }}
        </button>
        <button v-if="showOk" class="custom-dialog-btn custom-dialog-btn-ok" @click="handleOk" :disabled="loading">
          {{ okText }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { defineProps, defineEmits } from 'vue'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  title: {
    type: String,
    default: '弹窗标题'
  },
  content: {
    type: String,
    default: '这里是一些弹窗内容...'
  },
  cancelText: {
    type: String,
    default: '取消'
  },
  okText: {
    type: String,
    default: '确定'
  },
  closeOnOverlayClick: {
    type: Boolean,
    default: true
  },
  showFooter: {
    type: Boolean,
    default: true
  },
  showCancel: {
    type: Boolean,
    default: true
  },
  showOk: {
    type: Boolean,
    default: true
  },
  loading: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:modelValue', 'cancel', 'ok'])

const close = () => {
  emit('update:modelValue', false)
}

const handleOverlayClick = () => {
  if (props.closeOnOverlayClick) {
    close()
  }
}

const handleCancel = () => {
  emit('cancel')
  close()
}

const handleOk = () => {
  emit('ok')
  close()
}
</script>

<style scoped>
/* 默认暗色主题变量（跟随项目默认的无 .light 状态） */
.custom-dialog-overlay {
  --overlay-bg: rgba(0, 0, 0, 0.6);
  --dialog-bg: #1E1E1E;
  --border-color: #333333;
  --text-primary: #E0E0E0;
  --text-secondary: #A0A0A0;
  --btn-default-bg: #2A2A2A;
  --btn-default-text: #E0E0E0;
  --btn-primary-bg: #facc15; /* 截图暗色主题的黄色按钮 */
  --btn-primary-text: #010000;
  --btn-primary-hover: #fce72d;
  --btn-default-hover: #3a3a3a;
}

/* 适配全局亮色模式 (当 html 添加了 .light 类时) */

/* 遮罩层 */
.custom-dialog-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: var(--overlay-bg);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2000;
  backdrop-filter: blur(2px);
}

/* 弹窗主体 */
.custom-dialog-container {
  width: 480px;
  max-width: 90vw;
  background-color: var(--dialog-bg);
  border-radius: 8px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
  display: flex;
  flex-direction: column;
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
}

/* 标题栏 */
.custom-dialog-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 24px;
  border-bottom: 1px solid var(--border-color);
}

.custom-dialog-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
}

.custom-dialog-close {
  background: transparent;
  border: none;
  cursor: pointer;
  color: var(--text-secondary);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 4px;
  border-radius: 4px;
  transition: background-color 0.2s;
}

.custom-dialog-close:hover {
  background-color: var(--btn-default-bg);
}

/* 内容区 */
.custom-dialog-body {
  padding: 24px;
  color: var(--text-secondary);
  font-size: 14px;
  line-height: 1.5;
}

.custom-dialog-content-text {
  margin: 0;
  white-space: pre-wrap; /* 保留换行 */
  word-wrap: break-word; /* 允许长单词换行 */
  word-break: break-all; /* 在非常长的字符串内换行 */
}

/* 底栏 */
.custom-dialog-footer {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  padding: 12px 24px;
  gap: 12px;
  border-top: 1px solid var(--border-color);
}

/* 按钮通用 */
.custom-dialog-btn {
  padding: 6px 16px;
  border-radius: 4px;
  font-size: 14px;
  cursor: pointer;
  border: 1px solid transparent;
  transition: all 0.2s;
  font-weight: 500;
}

.custom-dialog-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

/* 取消按钮 */
.custom-dialog-btn-cancel {
  background-color: var(--btn-default-bg);
  color: var(--btn-default-text);
  border-color: var(--border-color);
}

.custom-dialog-btn-cancel:hover {
  background-color: var(--btn-default-hover);
}

/* 确认按钮 */
.custom-dialog-btn-ok {
  background-color: var(--btn-primary-bg);
  color: var(--btn-primary-text);
}

.custom-dialog-btn-ok:hover {
  background-color: var(--btn-primary-hover);
}
</style>

<style>
html.light .custom-dialog-overlay {
  --overlay-bg: rgba(0, 0, 0, 0.2) !important;
  --dialog-bg: #FFFFFF !important;
  --border-color: #EAEAEA !important;
  --text-primary: #333333 !important;
  --text-secondary: #666666 !important;
  --btn-default-bg: #F5F5F5 !important;
  --btn-default-text: #333333 !important;
  --btn-primary-bg: #1890FF !important;
  --btn-primary-text: #FFFFFF !important;
  --btn-primary-hover: #40a9ff !important;
  --btn-default-hover: #e8e8e8 !important;
}
</style>
