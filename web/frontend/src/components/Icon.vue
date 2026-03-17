<script setup lang="ts">
import SunIcon from './icons/SunIcon.vue'
import MoonIcon from './icons/MoonIcon.vue'
import MonitorIcon from './icons/MonitorIcon.vue'

const props = defineProps<{
  name: 'sun' | 'moon' | 'monitor'
  size?: number | string
}>()

const iconMap = {
  sun: SunIcon,
  moon: MoonIcon,
  monitor: MonitorIcon
}
</script>

<template>
  <component
      :is="iconMap[name]"
      :style="{ width: size + 'px', height: size + 'px'}"
  />
</template>