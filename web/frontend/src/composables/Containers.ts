import { reactive, computed, onMounted } from 'vue'
import { getContainerList, startContainer as apiStartContainer, stopContainer as apiStopContainer, removeContainer as apiRemoveContainer } from '@/services/modules/container'
import { ElMessage } from 'element-plus'

export function ContainerState() {
  const state = reactive({
    query: '',
    page: 1,
    pageSize: 10,
    loading: false,

    stats:[
        { label: '总数:', value: '0' },
        { label: '运行中:', value: '0', variant: 'running' },
        { label: '已停止:', value: '0', variant: 'stopped' }
    ],

    containers: [] as any[],
  })

  const loadData = async () => {
    try {
      state.loading = true
      const res = await getContainerList({ all: true })
      if (Array.isArray(res)) {
        state.containers = res.map((c: any) => {
          const name = (c.names && c.names.length > 0) ? c.names[0].replace(/^\//, '') : (c.id?.substring(0, 12) || 'unknown')
          const createdTime = c.created ? new Date(c.created * 1000).toLocaleString() : 'N/A'
          
          return {
            name,
            id: c.id.substring(0, 12),
            fullId: c.id,
            image: c.image || 'N/A',
            running: c.state === 'running',
            port: c.ports ? c.ports.map((p: any) => `${p.public_port || p.private_port}:${p.private_port}`).join(', ') : '',
            created: createdTime
          }
        })

        state.stats[0].value = state.containers.length.toString()
        state.stats[1].value = state.containers.filter(c => c.running).length.toString()
        state.stats[2].value = state.containers.filter(c => !c.running).length.toString()
      }
    } catch (e: any) {
      ElMessage.error(e.message || '获取容器列表失败')
    } finally {
      state.loading = false
    }
  }

  onMounted(() => {
    loadData()
  })

  const startContainer = async (id: string) => {
    try {
      await apiStartContainer(id)
      ElMessage.success('容器已启动')
      await loadData()
    } catch(e: any) {
      ElMessage.error(e.message || '启动容器失败')
    }
  }

  const stopContainer = async (id: string) => {
    try {
      await apiStopContainer(id)
      ElMessage.success('容器已停止')
      await loadData()
    } catch(e: any) {
      ElMessage.error(e.message || '停止容器失败')
    }
  }

  const removeContainer = async (id: string) => {
    try {
      await apiRemoveContainer(id, { force: true })
      ElMessage.success('容器已删除')
      await loadData()
    } catch(e: any) {
      ElMessage.error(e.message || '删除容器失败')
    }
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