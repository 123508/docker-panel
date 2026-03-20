<template>
  <div class="dp-pagination">

    <!-- 上一页 -->
    <dp-button
        text="上一页"
        size="medium"
        variant="outlined"
        :disabled="page === 1"
        @click="change(page - 1)"
    />

    <!-- 页码 -->
    <span
        v-for="p in pages"
        :key="p"
        class="page-item"
        :class="{ active: p === page }"
        @click="change(p)"
    >
      {{ p }}
    </span>

    <!-- 下一页 -->
    <dp-button
        text="下一页"
        size="medium"
        variant="outlined"
        :disabled="page === totalPages"
        @click="change(page + 1)"
    />

  </div>
</template>

<script setup lang="ts">
import {computed} from "vue";
import DpButton from "@/components/dp-button.vue";

defineOptions({
  name:'dp-pagination'
})

const props = defineProps<{
  page: number
  total: number
  pageSize: number
}>()

const emit = defineEmits<{
  (e:'change',page:number):void
}>()

const totalPages = computed(() => Math.ceil(props.total / props.pageSize))

const pages = computed(() => {
  return Array.from({length: totalPages.value}, (_, i) => i + 1)
})

function change(p:number){
  if(p<1 || p>totalPages.value) return
  emit('change',p)
}
</script>

<style scoped>
.dp-pagination {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 16px;
  justify-content: flex-end;
}

/* 页码基础样式 */
.page-item {
  min-width: 24px;
  height: 24px;
  padding: 0 6px;
  display: flex;
  align-items: center;
  justify-content: center;

  cursor: pointer;
  font-size: 12px;
}

/* hover */
.page-item:hover {
  background: var(--table-row-hover);
}

/* 当前页 */
.page-item.active {
  color: var(--accent);
}
</style>