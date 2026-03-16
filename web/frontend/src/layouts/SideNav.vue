<template>
  <nav class="nav">
    <div class="logo">
      <span class="whale">🐳</span>
      <span class="docker-text">DOCKER</span>
    </div>

    <ul class="menu">
      <li v-for="menu in menus" :key="menu.path">
        <router-link :to="menu.path" class="menu-item" @click="clearFooterSelection">
          <span class="icon">{{ menu.icon }}</span>
          <span class="label">{{ menu.name }}</span>
        </router-link>
      </li>
    </ul>

    <div class="footer">
      <div class="toggle-container">
        <div class="toggle-slider" :style="{ transform: `translateX(${activeIndex * 44}px)` }"></div>
        <button
            v-for="(btn, index) in footerButtons"
            :key="index"
            class="footer-btn"
            :class="{ active: activeIndex === index }"
            @click="activeIndex = index"
        >
          <span :class="['btn-icon', btn.icon]"></span>
        </button>
      </div>
    </div>
  </nav>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const activeIndex = ref(1)

const menus = [
  { name: '仪表盘', path: '/', icon: '01' },
  { name: '容器', path: '/containers', icon: '02' },
  { name: '镜像', path: '/images', icon: '03' },
  { name: '卷', path: '/volumes', icon: '04' },
  { name: '网络', path: '/networks', icon: '05' },
  { name: '编排', path: '/compose', icon: '06' }
]

const footerButtons = [
  { icon: 'icon-window' },
  { icon: 'icon-circle' },
  { icon: 'icon-refresh' }
]

const clearFooterSelection = () => {
  activeIndex.value = -1
}
</script>

<style scoped>
.nav {
  width: 230px;
  height: 100vh;
  background: #f0f4f8;
  color: #333;
  display: flex;
  flex-direction: column;
  border-right: 1px solid #e0e6ed;
}

.logo {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 16px 20px;
  border-bottom: 1px solid #e0e6ed;
  height: 60px;
}

.whale {
  width: 20px;
  height: 20px;
  background: #0099ff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  border-radius: 2px;
}

.docker-text {
  font-size: 12px;
  font-weight: 600;
  letter-spacing: 1px;
  color: #333;
}

.menu {
  list-style: none;
  padding: 0;
  margin: 0;
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.menu li {
  margin: 0;
}

.menu-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 20px;
  text-decoration: none;
  color: #666;
  transition: all 0.3s;
  font-size: 13px;
  border-left: 3px solid transparent;
  cursor: pointer;
}

.menu-item:hover {
  background: #e6f0f8;
  color: #333;
}

.router-link-active {
  background: #d4e8f5;
  color: #0099ff;
  border-left-color: #0099ff;
  font-weight: 500;
}

.icon {
  width: 20px;
  text-align: center;
  font-size: 11px;
  color: #999;
  font-weight: 500;
}

.router-link-active .icon {
  color: #0099ff;
  font-weight: 600;
}

.label {
  flex: 1;
}

.footer {
  display: flex;
  gap: 8px;
  padding: 16px 20px;
  border-top: 1px solid #e0e6ed;
  justify-content: center;
}

.toggle-container {
  position: relative;
  display: flex;
  gap: 6px;
  background: #e6eef5;
  border-radius: 50px;
  padding: 6px;
  width: fit-content;
}

.toggle-slider {
  position: absolute;
  width: 36px;
  height: 36px;
  background: #0099ff;
  border-radius: 50%;
  left: 6px;
  top: 6px;
  transition: transform 0.3s cubic-bezier(0.4, 0.0, 0.2, 1);
  pointer-events: none;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  opacity: 1;
}

.toggle-container:has(.footer-btn:nth-child(n+4):not(.active)) .toggle-slider {
  opacity: 0;
  pointer-events: none;
}

.footer-btn {
  width: 36px;
  height: 36px;
  border: none;
  background: transparent;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: color 0.3s;
  position: relative;
  z-index: 1;
  color: #666;
  border-radius: 50%;
}

.footer-btn:hover {
  color: #333;
}

.footer-btn.active {
  color: white;
  font-weight: 600;
}

.btn-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
}

/* 窗口/方块图标 */
.icon-window {
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.icon-window::before {
  content: '';
  width: 14px;
  height: 14px;
  border: 2px solid currentColor;
  border-radius: 1px;
}

/* 圆形图标 */
.icon-circle {
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.icon-circle::before {
  content: '';
  width: 14px;
  height: 14px;
  border: 2px solid currentColor;
  border-radius: 50%;
}

/* 刷新图标 */
.icon-refresh {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  transform: rotate(0deg);
}

.icon-refresh::before {
  content: '';
  width: 12px;
  height: 12px;
  border: 2px solid currentColor;
  border-right: none;
  border-bottom: none;
  border-radius: 50%;
  position: relative;
}

.icon-refresh::after {
  content: '';
  position: absolute;
  width: 4px;
  height: 4px;
  background: currentColor;
  top: 1px;
  right: 1px;
  clip-path: polygon(0 0, 100% 0, 0 100%);
}
</style>