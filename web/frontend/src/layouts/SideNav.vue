<template>
  <nav class="nav">
    <div class="logo">
      <div class="logo-icon"></div>
      <span class="logo-text">DOCKER</span>
    </div>

    <ul class="menu">
      <li v-for="item in menus" :key="item.path">
        <router-link :to="item.path" class="menu-item" exact-active-class="active">
          <span class="menu-number">{{ item.number }}</span>
          <span class="menu-label">{{ item.name }}</span>
        </router-link>
      </li>
    </ul>

    <div class="footer">
      <div class="toggle-container">
        <div class="toggle-slider" :style="sliderStyle"></div>
        <button
          v-for="(btn, i) in themeButtons"
          :key="i"
          class="toggle-btn"
          :class="{ active: themeIndex === i }"
          @click="setTheme(i)"
          :title="btn.title"
        >
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path v-if="btn.icon === 'monitor'" d="M2 3h20v14H2zM8 21h8M12 17v4" />
            <circle v-if="btn.icon === 'sun'" cx="12" cy="12" r="5" />
            <g v-if="btn.icon === 'sun'">
              <line x1="12" y1="1" x2="12" y2="3" />
              <line x1="12" y1="21" x2="12" y2="23" />
              <line x1="4.22" y1="4.22" x2="5.64" y2="5.64" />
              <line x1="18.36" y1="18.36" x2="19.78" y2="19.78" />
              <line x1="1" y1="12" x2="3" y2="12" />
              <line x1="21" y1="12" x2="23" y2="12" />
              <line x1="4.22" y1="19.78" x2="5.64" y2="18.36" />
              <line x1="18.36" y1="5.64" x2="19.78" y2="4.22" />
            </g>
            <path v-if="btn.icon === 'moon'" d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z" />
          </svg>
        </button>
      </div>
    </div>
  </nav>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { themePreference, setThemePreference } from '@/store/ui'

const themeMap = ['system', 'light', 'dark'] as const

const menus = [
  { name: '仪表盘', path: '/', number: '01' },
  { name: '容器', path: '/containers', number: '02' },
  { name: '镜像', path: '/images', number: '03' },
  { name: '卷', path: '/volumes', number: '04' },
  { name: '网络', path: '/networks', number: '05' },
  { name: '编排', path: '/compose', number: '06' }
]

const themeButtons = [
  { icon: 'monitor', title: '跟随系统' },
  { icon: 'sun', title: '浅色模式' },
  { icon: 'moon', title: '深色模式' }
]

const themeIndex = computed(() => themeMap.indexOf(themePreference.value))

const sliderStyle = computed(() => ({
  transform: `translateX(${themeIndex.value * 40}px)`
}))

function setTheme(index: number) {
  setThemePreference(themeMap[index])
}
</script>

<style scoped>
.nav {
  width: var(--sidebar-width);
  height: 100vh;
  background: var(--bg-sidebar);
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.logo {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 24px 24px 16px;
}

.logo-icon {
  width: 18px;
  height: 18px;
  background: var(--accent);
  border-radius: 2px;
}

.logo-text {
  font-family: var(--font-mono);
  font-size: 14px;
  font-weight: 700;
  letter-spacing: 2px;
  color: var(--text-primary);
}

.menu {
  list-style: none;
  padding: 0 16px;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.menu-item {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 14px 16px;
  border-radius: 0;
  font-family: var(--font-mono);
  font-size: 13px;
  color: var(--text-secondary);
  transition: color 0.2s;
}

.menu-item:hover {
  color: var(--text-muted);
}

.menu-item.active {
  color: var(--accent);
  background: var(--nav-active-bg);
  border-radius: 8px;
}

.menu-number {
  font-size: 11px;
  font-weight: 500;
  min-width: 16px;
}

.menu-label {
  font-weight: 500;
}

.footer {
  padding: 24px;
  display: flex;
  justify-content: center;
}

.toggle-container {
  position: relative;
  display: flex;
  background: var(--toggle-bg);
  border-radius: 100px;
  padding: 4px;
}

.toggle-slider {
  position: absolute;
  width: 32px;
  height: 32px;
  background: var(--accent);
  border-radius: 50%;
  left: 4px;
  top: 4px;
  transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  pointer-events: none;
}

.toggle-btn {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  color: var(--text-secondary);
  position: relative;
  z-index: 1;
  padding: 0;
  margin: 0 4px;
  transition: color 0.2s;
}

.toggle-btn:first-child {
  margin-left: 0;
}

.toggle-btn:last-child {
  margin-right: 0;
}

.toggle-btn.active {
  color: var(--toggle-active-icon);
}

.toggle-btn svg {
  width: 14px;
  height: 14px;
}
</style>
