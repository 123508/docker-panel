<template>
  <div class="table-wrap" :class="{ 'bordered': bordered, 'flex': flex }">
    <div class="table-head" :class="{ 'has-bg': headerBg, 'has-gap': rowGap }">
      <slot name="head"></slot>
    </div>
    <div class="table-body" :class="{ 'has-gap': rowGap, 'fixed-row-height': fixedRowHeight, 'bordered-rows': bordered }">
      <slot></slot>
    </div>
  </div>
</template>

<script setup lang="ts">
withDefaults(defineProps<{
  headerBg?: boolean
  rowGap?: boolean
  fixedRowHeight?: boolean
  bordered?: boolean
  flex?: boolean
}>(), {
  bordered: true,
  flex: true
})
</script>

<style scoped>
.table-wrap {
  background: var(--bg-card);
  overflow: hidden;
}

.table-wrap.bordered {
  border: var(--border);
}

.table-wrap.flex {
  flex: 1;
}

.table-head {
  display: flex;
  align-items: center;
  padding: 0 20px;
  height: 44px;
  border-bottom: var(--border);
}

.table-head.has-bg {
  background: var(--bg-card-header);
  border-bottom: none;
}

.table-head.has-gap {
  gap: 16px;
}

/* Row Styles */
/* These styles are applied to slotted content using :deep() */
.table-body :deep(.table-row) {
  display: flex;
  align-items: center;
  padding: 16px 20px;
  transition: background 0.15s;
}

.table-body :deep(.table-row:hover) {
  background: var(--table-row-hover);
}

.table-body.has-gap :deep(.table-row) {
  gap: 16px;
}

.table-body.fixed-row-height :deep(.table-row) {
  height: 56px;
  padding: 0 20px;
}

.table-body.bordered-rows :deep(.table-row:not(:last-child)) {
  border-bottom: var(--border);
}
</style>
