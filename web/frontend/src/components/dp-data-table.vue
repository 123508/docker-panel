<template>
  <div class="table-wrap"
       :class="{ 'bordered': bordered, 'flex': flex }"
       :style="{
          height: typeof height === 'number' ? height + 'px' : height
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
        :style="{
        width: col.width ? col.width + 'px' : undefined,
        flex: col.flex
      }"
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

  defineOptions({
    name:'dp-data-table'
  })

  type Column = {
    key: string
    label?: string
    width?: number
    flex?: number
  }

  withDefaults(defineProps<{
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
</script>

<style scoped>
.table-wrap {
  background: var(--bg-card);
  overflow: hidden;

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
  height: 44px;
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
  overflow-y: auto;
  min-height: 0;
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
