import { reactive, computed } from 'vue'

export function VolumeState() {
  const state = reactive({
    page: 1,
    pageSize: 5,

    stats: [
      { label: '总数:', value: '8' },
      { label: '已使用:', value: '3.2 GB' },
      { label: '未使用:', value: '2', variant: 'dangling' }
    ],

    volumes: [
      { name: 'postgres_data', driver: 'local', mountpoint: '/var/lib/docker/volumes/postgres_data/_data', size: '1.2 GB', containers: 1 },
      { name: 'redis_cache', driver: 'local', mountpoint: '/var/lib/docker/volumes/redis_cache/_data', size: '512 MB', containers: 1 },
      { name: 'app_logs', driver: 'local', mountpoint: '/var/lib/docker/volumes/app_logs/_data', size: '256 MB', containers: 0 }
    ]
  })

  const pagedVolumes = computed(() => {
    const start = (state.page - 1) * state.pageSize
    return state.volumes.slice(start, start + state.pageSize)
  })

  return {
    state,
    pagedVolumes
  }
}