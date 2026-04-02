import { reactive, computed } from 'vue'

export function ContainerState() {
  const state = reactive({
    query: '',
    page: 1,
    pageSize: 10,

    stats:[
        { label: '总数:', value: '24' },
        { label: '运行中:', value: '18', variant: 'running' },
        { label: '已停止:', value: '6', variant: 'stopped' }
    ],

    containers: [
        { name: 'nginx-web-01', id: 'a3f8d92b', image: 'nginx:latest', running: true, port: '80:80, 443:443', created: '2 小时前' },
        { name: 'postgres-db', id: 'e7b2c41a', image: 'postgres:14', running: true, port: '5432:5432', created: '5 天前' },
        { name: 'redis-cache', id: 'f4c9e3d2', image: 'redis:alpine', running: false, port: '6379:6379', created: '1 天前' }
    ],
})

  const startContainer = (id: string) => {
    const c = state.containers.find(c => c.id === id)
    if (c) c.running = true
  }

  const stopContainer = (id: string) => {
    const c = state.containers.find(c => c.id === id)
    if (c) c.running = false
  }

  const removeContainer = (id: string) => {
    const index = state.containers.findIndex(c => c.id === id)
    if (index !== -1) state.containers.splice(index, 1)
  }

  const filteredContainers = computed(() => {
    const query = state.query.toLowerCase()
    if (!query) return state.containers
    return state.containers.filter(c => c.name.toLowerCase().includes(query) || c.image.toLowerCase().includes(query))
  })

  const pagedContainers = computed(() => {
    const start = (state.page - 1) * state.pageSize
    return filteredContainers.value.slice(start, start + state.pageSize)
  })

  return {
    state,
    filteredContainers,
    pagedContainers,
    startContainer,
    stopContainer,
    removeContainer
  }
}