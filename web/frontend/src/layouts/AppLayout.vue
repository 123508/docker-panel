<template>
  <div class="layout">
    <SideNav class="sidebar" />
    <main class="content">
      <router-view />
    </main>
  </div>
</template>

<script setup lang="ts">
import SideNav from './SideNav.vue'
</script>

<style scoped>
.layout {
  display: flex;
  height: 100vh;
  background: var(--bg-body);
}

.sidebar {
  width: var(--sidebar-width);
  flex-shrink: 0;
}

.content {
  flex: 1;
  overflow-y: auto;
  min-width: 0;
}
</style>
