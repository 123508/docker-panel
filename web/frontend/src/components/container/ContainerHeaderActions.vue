<template>
  <div class="header-actions">
    <slot />
  </div>
</template>

<style scoped>
.header-actions {
  display: flex;
  gap: 10px;
  align-items: center;
  flex-wrap: wrap;
}
</style>
