<template>
  <div class="header-actions">
    <slot />
  </div>
</template>

<script setup lang="ts">
defineOptions({
  name: 'dp-container-header-actions'
})
</script>

<style scoped>
.header-actions {
  display: flex;
  gap: 10px;
  align-items: center;
  flex-wrap: wrap;
}
</style>
