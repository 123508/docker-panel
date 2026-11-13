<template>
  <nav class="nav">
    <div class="logo">
      <div class="logo-icon"><span>▦</span></div>
      <div class="logo-copy"><span class="logo-text">docker<span class="logo-accent">.</span>panel</span><span class="logo-caption">CONTAINER WORKSPACE</span></div>
    </div>

    <div class="menu-area"><div class="menu-heading">WORKSPACE</div><ul class="menu">
      <li v-for="item in menus" :key="item.path">
        <router-link :to="item.path" class="menu-item" exact-active-class="active">
          <span class="menu-number">{{ item.number }}</span>
          <span class="menu-label">{{ item.name }}</span>
        </router-link>
      </li>
    </ul></div>

    <div class="footer"><div class="footer-label">外观设置</div>
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
          <dp-icon :name="btn.icon" size="18"/>
        </button>
      </div>
    </div>
  </nav>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { themePreference, setThemePreference } from '@/store/ui'
import DpIcon from "@/components/dp-icon.vue";

const themeMap = ['system', 'light', 'dark'] as const

const menus = [
  { name: '仪表盘', path: '/dashboard', number: '01' },
  { name: '容器', path: '/dashboard/containers', number: '02' },
  { name: '镜像', path: '/dashboard/images', number: '03' },
  { name: '卷', path: '/dashboard/volumes', number: '04' },
  { name: '网络', path: '/dashboard/networks', number: '05' },
  { name: '编排', path: '/dashboard/compose', number: '06' }
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
  left: 8px;
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

<style scoped>
.nav { border-right: var(--border); justify-content: flex-start; }
.logo { padding: 32px 22px 42px; gap: 12px; }
.logo-icon { width: 38px; height: 38px; border-radius: 12px; display: grid; place-items: center; color: var(--brand-icon-text); font-size: 25px; line-height: 1; background: var(--brand-gradient); box-shadow: 0 6px 18px var(--brand-shadow); }
.logo-copy { display: flex; flex-direction: column; gap: 2px; }
.logo-text { font-family: var(--font-display); font-size: 17px; letter-spacing: -.6px; }
.logo-accent { color: var(--accent); }
.logo-caption { font: 700 8px var(--font-mono); letter-spacing: 1px; color: var(--text-secondary); }
.menu-area { flex: 1; }
.menu-heading { padding: 0 30px 12px; font: 700 10px var(--font-mono); letter-spacing: 1.5px; color: var(--text-dim); }
.menu { gap: 5px; }
.menu-item { border-radius: 9px; padding: 12px 14px; font-family: var(--font-display), sans-serif; }
.menu-item:hover { background: var(--table-row-hover); color: var(--text-primary); }
.menu-item.active { font-weight: 700; }
.menu-number { font-family: var(--font-mono); opacity: .7; }
.footer { flex-direction: column; align-items: flex-start; gap: 12px; border-top: var(--border); }
.footer-label { font-size: 11px; color: var(--text-secondary); }
@media (max-width: 700px) { .nav { width: 100%; height: auto; border-right: 0; border-bottom: var(--border); } .logo { padding: 16px 18px; } .menu-area { overflow-x: auto; } .menu-heading, .menu-number, .footer-label { display: none; } .menu { flex-direction: row; width: max-content; padding: 0 12px 10px; } .menu-item { padding: 8px 12px; } .footer { position: absolute; right: 16px; top: 14px; padding: 0; border: 0; } .toggle-container { transform: scale(.8); transform-origin: right top; } }
</style>
