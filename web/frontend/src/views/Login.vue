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
        <div class="login-mark">▦</div><div class="login-kicker">CONTAINER WORKSPACE</div><h1 class="title">欢迎回来</h1>
        <p class="subtitle">登录 Docker Panel，继续管理你的工作空间</p>
      </div>

      <div class="form">
        <div class="form-group">
          <label class="label">用户名</label>
          <div class="input-wrapper">
            <input type="text" autocomplete="username" class="input" placeholder="请输入用户名" v-model="username" @keyup.enter="handleLogin" />
          </div>
        </div>

        <div class="form-group">
          <label class="label">密码</label>
          <div class="input-wrapper">
            <input :type="showPassword ? 'text' : 'password'" autocomplete="current-password" class="input" placeholder="请输入密码" v-model="password" @keyup.enter="handleLogin" />
            <button class="eye-btn" :aria-label="showPassword ? '隐藏密码' : '显示密码'" @click="showPassword = !showPassword" type="button">
              {{ showPassword ? '隐藏' : '显示' }}
            </button>
          </div>
        </div>

        <div v-if="errorMessage" class="error-msg">{{ errorMessage }}</div>

        <button class="btn" @click="handleLogin" :disabled="isLoading">
          {{ isLoading ? '登录中...' : '登录' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { themePreference, setThemePreference } from '@/store/ui'
import { setToken } from '@/store/auth'
import { http } from '@/services/api'
import DpIcon from "@/components/dp-icon.vue"

const router = useRouter()
const username = ref('')
const password = ref('')
const showPassword = ref(false)
const errorMessage = ref('')
const isLoading = ref(false)

const toggleTheme = () => {
  const newTheme = themePreference.value === 'dark' ? 'light' : 'dark'
  setThemePreference(newTheme)
}

const handleLogin = async () => {
  if (!username.value || !password.value) {
    errorMessage.value = '请输入用户名和密码'
    return
  }

  errorMessage.value = ''
  isLoading.value = true

  try {
    const response = await http.post('/api/login', {
      username: username.value,
      password: password.value
    })

    const res = response.data
    if (res.code === 0 && res.data && res.data.token) {
      setToken(res.data.token)
      router.push('/dashboard')
    } else {
      errorMessage.value = res.message || '登录失败'
    }
  } catch (error: any) {
    if (error.response && error.response.data && error.response.data.message) {
      errorMessage.value = error.response.data.message
    } else {
      errorMessage.value = '登录请求失败，请检查网络'
    }
  } finally {
    isLoading.value = false
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
  transition: background-color 0.2s, opacity 0.2s;
}

.btn:hover:not(:disabled) {
  opacity: 0.9;
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.error-msg {
  color: #f56c6c;
  font-size: 13px;
  text-align: center;
  margin-top: -10px;
  margin-bottom: -10px;
}
</style>

<style scoped>
.login-container { min-height: 100vh; height: auto; padding: 72px 20px; background: radial-gradient(ellipse at 20% 10%, var(--login-glow-one), transparent 42%), radial-gradient(ellipse at 85% 85%, var(--login-glow-two), transparent 36%), var(--bg-body); }
.login-card { width: min(440px, 100%); padding: 44px; border-radius: 20px; background: var(--bg-card); box-shadow: 0 24px 70px #05172b24; }
.login-mark { width: 48px; height: 48px; display: grid; place-items: center; border-radius: 14px; background: var(--brand-gradient); box-shadow: 0 8px 20px var(--brand-shadow); color: var(--brand-icon-text); font-size: 30px; margin-bottom: 12px; }
.login-kicker { font: 700 10px var(--font-mono); letter-spacing: 1.4px; color: var(--accent); }
.title { font-family: var(--font-display), sans-serif; font-size: 27px; }
.subtitle { font-family: var(--font-display), sans-serif; font-size: 13px; }
.input-wrapper { height: 46px; border-radius: 9px; background: var(--bg-body); }
.btn { height: 46px; border-radius: 9px; }
@media (max-width: 480px) { .login-card { padding: 30px 24px; } }
</style>
