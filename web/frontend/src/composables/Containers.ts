import { reactive, computed, onMounted } from 'vue'
import {
  getContainerList,
  startContainer as apiStartContainer,
  stopContainer as apiStopContainer,
  restartContainer as apiRestartContainer,
  removeContainer as apiRemoveContainer
} from '@/services/modules/container'
import { ElMessage } from 'element-plus'
import { useActionDialog } from '@/composables/useActionDialog'

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

  const { dialog, runWithDialog } = useActionDialog()

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

  const shortId = (id: string) => id.substring(0, 12)

  const startContainer = async (id: string) => {
    await runWithDialog(
      {
        title: '启动容器',
        pendingText: `正在启动容器 ${shortId(id)}，请稍候...`,
        successText: `✅ 容器已启动: ${shortId(id)}`,
        failureText: (e) => `❌ 启动容器失败: \n${e?.message || '未知错误'}`
      },
      async () => {
        await apiStartContainer(id)
        await loadData()
      }
    )
  }

  const stopContainer = async (id: string) => {
    await runWithDialog(
      {
        title: '停止容器',
        pendingText: `正在停止容器 ${shortId(id)}，请稍候...`,
        successText: `✅ 容器已停止: ${shortId(id)}`,
        failureText: (e) => `❌ 停止容器失败: \n${e?.message || '未知错误'}`
      },
      async () => {
        await apiStopContainer(id)
        await loadData()
      }
    )
  }

  const removeContainer = async (id: string) => {
    await runWithDialog(
      {
        title: '移除容器',
        pendingText: `正在移除容器 ${shortId(id)}，请稍候...`,
        successText: `✅ 容器已移除: ${shortId(id)}`,
        failureText: (e) => `❌ 移除容器失败: \n${e?.message || '未知错误'}`
      },
      async () => {
        await apiRemoveContainer(id, { force: true })
        await loadData()
      }
    )
  }

  const restartContainer = async (id: string) => {
    await runWithDialog(
      {
        title: '重启容器',
        pendingText: `正在重启容器 ${shortId(id)}，请稍候...`,
        successText: `✅ 容器已重启: ${shortId(id)}`,
        failureText: (e) => `❌ 重启容器失败: \n${e?.message || '未知错误'}`
      },
      async () => {
        await apiRestartContainer(id)
        await loadData()
      }
    )
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
    dialog,
    filteredContainers,
    pagedContainers,
    loadData,
    startContainer,
    stopContainer,
    restartContainer,
    removeContainer
  }
}
