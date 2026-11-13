<template>
  <div class="table-wrap"
       :class="{ 'bordered': bordered, 'flex': flex, 'has-columns': columns.length }"
       :style="{
          height: typeof height === 'number' ? height + 'px' : height,
          '--column-layout': columnLayout
       }">
    <div class="table-head"
         :class="{
            'has-bg': headerBg, 'has-gap': rowGap
         }">
      <!-- 新columns 渲染 -->
      <template v-if="columns.length">
    <span
        v-for="col in columns"
        :key="col.key"
        class="th"
    >
      {{ col.label }}
    </span>
      </template>

      <!-- 兼容旧 slot -->
      <slot v-else name="head"></slot>
    </div>
    <div class="table-body" :class="{ 'has-gap': rowGap, 'fixed-row-height': fixedRowHeight, 'bordered-rows': bordered }">
      <slot></slot>
    </div>
  </div>
</template>

<script setup lang="ts">

  import { computed } from 'vue'

  defineOptions({
    name:'dp-data-table'
  })

  type Column = {
    key: string
    label?: string
    width?: number
    flex?: number
  }

  const props = withDefaults(defineProps<{
    headerBg?: boolean
    rowGap?: boolean
    fixedRowHeight?: boolean
    bordered?: boolean
    flex?: boolean
    columns?: Column[]
    height?: string | number
  }>(), {
    bordered: true,
    flex: true,
    columns: () => []
  })

  const columnLayout = computed(() => props.columns.map((col) =>
    col.key === 'icon' ? `${col.width || 40}px` : `minmax(0, ${col.width || 120}fr)`
  ).join(' '))
</script>

<style scoped>
.table-wrap {
  background: var(--bg-card);
  overflow: auto;
  border-radius: 14px;

  display: flex;
  flex-direction: column;
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
  min-width: 0;
  height: 52px;
  border-bottom: var(--border);
  flex-shrink: 0;
  position: sticky;
  top: 0;
  z-index: 1;
}

.table-head.has-bg {
  background: var(--bg-card-header);
  border-bottom: none;
}

.table-head.has-gap {
  gap: 16px;
}

.table-body {
  flex: 1;
  min-width: 0;
  overflow: visible;
  min-height: 0;
}

.has-columns .table-head,
.has-columns .table-body :deep(.table-row) {
  display: grid;
  grid-template-columns: var(--column-layout);
}
.has-columns .table-head > .th,
.has-columns .table-body :deep(.table-row > .td) {
  width: auto !important;
  min-width: 0;
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

.th {
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 1px;
  color: var(--text-secondary);
}
</style>

<style scoped>
@media (max-width: 1000px) {
  .table-wrap { height: auto !important; min-height: 160px; overflow: visible; }
  .table-wrap.flex { flex: initial; }
  .has-columns .table-head { display: none; }
  .table-body { min-width: 0; overflow: visible; }
  .has-columns .table-body :deep(.table-row), .table-body.fixed-row-height :deep(.table-row) {
    height: auto; display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 16px; padding: 20px;
  }
  .table-body :deep(.table-row:not(:last-child)) { border-bottom: var(--border); }
  .table-body :deep(.table-row > .td) { width: auto !important; min-width: 0; display: flex; flex-direction: column; align-items: flex-start; gap: 5px; overflow: visible; white-space: normal; overflow-wrap: anywhere; }
  .table-body :deep(.table-row > .td::before) { content: attr(data-label); color: var(--text-secondary); font-size: 10px; font-weight: 700; letter-spacing: .4px; }
  .table-body :deep(.table-row > [data-table-icon]) { display: none !important; }
  .table-body :deep(.table-row > .td-actions) { grid-column: 1 / -1; display: flex; flex-direction: row; align-items: center; flex-wrap: wrap; gap: 18px; padding-top: 2px; }
  .table-body :deep(.table-row > .td-actions::before) { display: none; }
  .table-body :deep(.name-col) { width: 100%; }
  .table-body :deep(.name-text), .table-body :deep(.name-id) { max-width: 100%; white-space: normal; overflow-wrap: anywhere; }
}
</style>
