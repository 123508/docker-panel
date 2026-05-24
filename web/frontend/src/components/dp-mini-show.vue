<template>
  <Teleport to="body">
    <div v-if="modelValue" class="mini-overlay" @click="onOverlayClick">
      <div class="mini-panel" @click.stop>
        <slot />
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
defineOptions({ name: 'dp-mini-show' })

const props = withDefaults(
  defineProps<{
    modelValue: boolean
    closeOnOverlay?: boolean
  }>(),
  {
    closeOnOverlay: true
  }
)

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
}>()

const onOverlayClick = () => {
  if (!props.closeOnOverlay) return
  emit('update:modelValue', false)
}
</script>

<style scoped>
.mini-overlay {
  position: fixed;
  inset: 0;
  z-index: 2000;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--mini-overlay-bg);
}

.mini-panel {
  width: min(420px, calc(100vw - 32px));
  border: var(--border);
  background: var(--mini-bg);
  box-shadow: 0 20px 40px var(--mini-shadow);
  padding: 20px;
}
</style>
