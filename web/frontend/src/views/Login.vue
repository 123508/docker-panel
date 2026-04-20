<template>
  <div class="login-container">
    <div class="theme-toggle">
      <button
        class="toggle-btn"
        @click="toggleTheme"
        :title="themePreference === 'dark' ? '切换为浅色模式' : '切换为深色模式'"
      >
        <dp-icon :name="themePreference === 'dark' ? 'moon' : 'sun'" size="16"/>
      </button>
    </div>

    <div class="login-card">
      <div class="header">
        <h1 class="title">Docker Panel</h1>
        <p class="subtitle">登录以管理您的容器</p>
      </div>

      <div class="form">
        <div class="form-group">
          <label class="label">用户名</label>
          <div class="input-wrapper">
            <input type="text" class="input" placeholder="请输入用户名" v-model="username" @keyup.enter="handleLogin" />
          </div>
        </div>

        <div class="form-group">
          <label class="label">密码</label>
          <div class="input-wrapper">
            <input :type="showPassword ? 'text' : 'password'" class="input" placeholder="请输入密码" v-model="password" @keyup.enter="handleLogin" />
            <button class="eye-btn" @click="showPassword = !showPassword" type="button">
              <dp-icon :name="showPassword ? 'eye' : 'eye-off'" size="18"/>
            </button>
          </div>
        </div>

        <div v-if="errorMessage" class="error-msg">{{ errorMessage }}</div>

        <button class="btn" @click="handleLogin">登录</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { themePreference, setThemePreference } from '@/store/ui'
import DpIcon from "@/components/dp-icon.vue"

const router = useRouter()
const username = ref('')
const password = ref('')
const showPassword = ref(false)
const errorMessage = ref('')

const toggleTheme = () => {
  const newTheme = themePreference.value === 'dark' ? 'light' : 'dark'
  setThemePreference(newTheme)
}

const handleLogin = () => {
  errorMessage.value = ''
  if (username.value === 'admin' && password.value === 'admin') {
    router.push('/dashboard')
  } else {
    errorMessage.value = '用户名或密码不正确'
  }
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background-color: var(--bg-body);
  position: relative;
}

.theme-toggle {
  position: absolute;
  top: 24px;
  right: 24px;
}

.toggle-btn {
  width: 36px;
  height: 36px;
  border-radius: 18px;
  border: 1px solid var(--border-color);
  background-color: transparent;
  display: flex;
  justify-content: center;
  align-items: center;
  cursor: pointer;
  color: var(--text-primary);
  transition: all 0.2s;
}

.toggle-btn:hover {
  background-color: var(--bg-card);
}

.login-card {
  width: 400px;
  background-color: var(--bg-body);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 40px;
  display: flex;
  flex-direction: column;
  gap: 32px;
}

.header {
  display: flex;
  flex-direction: column;
  gap: 8px;
  align-items: center;
}

.title {
  margin: 0;
  font-family: Inter, sans-serif;
  font-size: 28px;
  font-weight: bold;
  color: var(--text-primary);
  text-align: center;
}

.subtitle {
  margin: 0;
  font-family: Inter, sans-serif;
  font-size: 14px;
  color: var(--text-secondary);
  text-align: center;
}

.form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.label {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
}

.input-wrapper {
  display: flex;
  align-items: center;
  height: 40px;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  padding: 0 12px;
  background-color: var(--bg-body);
  transition: border-color 0.2s;
}

.input-wrapper:focus-within {
  border-color: var(--accent);
}

.input {
  flex: 1;
  border: none;
  background: transparent;
  outline: none;
  color: var(--text-primary);
  font-size: 14px;
}

.input::placeholder {
  color: var(--text-secondary);
}

.eye-btn {
  background: none;
  border: none;
  padding: 0;
  margin-left: 8px;
  cursor: pointer;
  color: var(--text-secondary);
  display: flex;
  align-items: center;
}

.eye-btn:hover {
  color: var(--text-primary);
}

.btn {
  height: 40px;
  background-color: var(--button-primary);
  color: var(--button-primary-text);
  border: none;
  border-radius: 4px;
  font-size: 14px;
  font-weight: bold;
  cursor: pointer;
  display: flex;
  justify-content: center;
  align-items: center;
  transition: background-color 0.2s;
}

.btn:hover {
  opacity: 0.9;
}

.error-msg {
  color: #f56c6c;
  font-size: 13px;
  text-align: center;
  margin-top: -10px;
  margin-bottom: -10px;
}
</style>
