<template>
  <div class="card">
    <h2>{{ title }}</h2>

    <div v-if="items.length" class="list">
      <div
          v-for="item in items"
          :key="item.id"
          class="list-item"
      >
        <slot name="item" :item="item"/>
      </div>
    </div>

    <div v-else class="empty">
      <slot name="empty"/>
    </div>
  </div>
</template>

<script setup lang="ts" generic="T extends { id: string }">

  defineProps({
    title:{
      type : String,
      default: 'Card Title',
    },
    items:{
      type: Array<T>,
      default:() => [],
    }
  })

  defineSlots<{
    item(props:{ item: T }):any
    empty?:()=>any
  }>()

</script>

<style scoped>
.card {
  background: white;
  border-radius: 12px;
  padding: 1.5rem;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  transition: transform 0.2s;
}

.card:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 12px rgba(0, 0, 0, 0.15);
}

.card h2 {
  margin: 0 0 1rem;
  color: #333;
  font-size: 1.25rem;
  border-bottom: 2px solid #667eea;
  padding-bottom: 0.5rem;
}

.list {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.list-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.75rem;
  background: #f9fafb;
  border-radius: 8px;
  border-left: 3px solid #667eea;
}

.empty {
  text-align: center;
  color: #999;
  padding: 2rem;
  font-style: italic;
}

</style>