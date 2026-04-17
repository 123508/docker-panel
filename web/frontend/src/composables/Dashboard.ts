import { reactive, computed } from 'vue'

export function DashboardState() {
  const state = reactive({
    page: 1,
    pageSize: 5,

    metrics: [
      { label: '容器', value: '24', sub: '18 运行中', subClass: 'success' },
      { label: '镜像', value: '42', sub: '12.4 GB', subClass: 'muted' },
      { label: '卷', value: '8', sub: '3.2 GB 已使用', subClass: 'muted' },
      { label: '网络', value: '5', sub: '3 自定义', subClass: 'muted' }
    ],

    containers: [
      { name: 'nginx-web-01', id: 'a3f8d92b', image: 'nginx:latest', status: '运行中', running: true, port: '80:80, 443:443' },
      { name: 'postgres-db', id: 'e7b2c41a', image: 'postgres:14', status: '运行中', running: true, port: '5432:5432' },
      { name: 'redis-cache', id: 'f4c9e3d2', image: 'redis:alpine', status: '已停止', running: false, port: '6379:6379' }
    ],
  })

  const pagedContainers = computed(() => {
    const start = (state.page - 1) * state.pageSize
    return state.containers.slice(start, start + state.pageSize)
  })


  const startContainer = (id: string) => {
    const c = state.containers.find(c => c.id === id)
    if (c) {
      c.running = true
      c.status = '运行中'
    }
  }

  const stopContainer = (id: string) => {
    const c = state.containers.find(c => c.id === id)
    if (c) {
      c.running = false
      c.status = '已停止'
    }
  }

  return {
    state,
    pagedContainers,
    startContainer,
    stopContainer
  }
}