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
}

.sidebar {
  width: 240px;
  border-right: 1px solid #eee;
  background: #fafafa;
}

.content {
  flex: 1;
  padding: 24px;
  overflow: auto;
}
</style>