import { reactive, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import {
  getContainerList,
  startContainer as apiStartContainer,
  stopContainer as apiStopContainer,
  restartContainer as apiRestartContainer,
  removeContainer as apiRemoveContainer
} from '@/services/modules/container'
import { getImageList } from '@/services/modules/image'
import { getVolumeList, getVolumeInspect } from '@/services/modules/volume'
import { getNetworkList } from '@/services/modules/network'

function formatBytesToGB(bytes: number): string {
  if (!Number.isFinite(bytes) || bytes <= 0) return '0 GB'
  return `${(bytes / 1024 / 1024 / 1024).toFixed(2)} GB`
}

export function DashboardState() {
  const state = reactive({
    page: 1,
    pageSize: 5,
    loading: false,
    actionLoading: false,
    metrics: [
      { label: '容器', value: '0', sub: '0 运行中', subClass: 'success' },
      { label: '镜像', value: '0', sub: '0 GB', subClass: 'muted' },
      { label: '卷', value: '0', sub: '0 GB 已使用', subClass: 'muted' },
      { label: '网络', value: '0', sub: '0 自定义', subClass: 'muted' }
    ],
    containers: [] as any[]
  })

  const loadData = async () => {
    try {
      state.loading = true
      const [containersRes, imagesRes, volumesRes, networksRes] = await Promise.all([
        getContainerList({ all: true }),
        getImageList({ all: true }),
        getVolumeList(),
        getNetworkList()
      ])

      const containers = Array.isArray(containersRes) ? containersRes : []
      const images = Array.isArray(imagesRes) ? imagesRes : []
      const volumes = Array.isArray(volumesRes) ? volumesRes : []
      const networks = Array.isArray(networksRes) ? networksRes : []

      state.containers = containers.map((c: any) => {
        const name = Array.isArray(c.names) && c.names.length > 0
          ? c.names[0].replace(/^\//, '')
          : (c.id?.substring(0, 12) || 'unknown')
        const running = c.state === 'running'
        return {
          name,
          id: c.id?.substring(0, 12) || 'N/A',
          fullId: c.id,
          image: c.image || 'N/A',
          status: running ? '运行中' : '已停止',
          running,
          port: Array.isArray(c.ports)
            ? c.ports.map((p: any) => `${p.public_port || p.private_port}:${p.private_port}`).join(', ')
            : ''
        }
      })

      const runningCount = containers.filter((c: any) => c.state === 'running').length
      const imagesTotalSize = images.reduce((sum: number, img: any) => sum + Number(img.size || 0), 0)

      const volumeInspects = await Promise.all(
        volumes.map(async (v: any) => {
          try {
            return await getVolumeInspect(v.name)
          } catch {
            return null
          }
        })
      )
      const usedVolumeSize = volumeInspects.reduce(
        (sum: number, detail: any) => sum + Number(detail?.usage_data?.size || 0),
        0
      )
      const customNetworks = networks.filter((n: any) => !['bridge', 'host', 'none'].includes(n.name)).length

      state.metrics = [
        { label: '容器', value: String(containers.length), sub: `${runningCount} 运行中`, subClass: 'success' },
        { label: '镜像', value: String(images.length), sub: formatBytesToGB(imagesTotalSize), subClass: 'muted' },
        { label: '卷', value: String(volumes.length), sub: `${formatBytesToGB(usedVolumeSize)} 已使用`, subClass: 'muted' },
        { label: '网络', value: String(networks.length), sub: `${customNetworks} 自定义`, subClass: 'muted' }
      ]
    } catch (e: any) {
      ElMessage.error(e.message || '获取 Dashboard 数据失败')
    } finally {
      state.loading = false
    }
  }

  onMounted(() => {
    loadData()
  })

  const startContainer = async (id: string) => {
    try {
      state.actionLoading = true
      await apiStartContainer(id)
      ElMessage.success('容器已启动')
      await loadData()
    } catch (e: any) {
      ElMessage.error(e.message || '启动容器失败')
    } finally {
      state.actionLoading = false
    }
  }

  const stopContainer = async (id: string) => {
    try {
      state.actionLoading = true
      await apiStopContainer(id)
      ElMessage.success('容器已停止')
      await loadData()
    } catch (e: any) {
      ElMessage.error(e.message || '停止容器失败')
    } finally {
      state.actionLoading = false
    }
  }

  const restartContainer = async (id: string) => {
    try {
      state.actionLoading = true
      await apiRestartContainer(id)
      ElMessage.success('容器已重启')
      await loadData()
    } catch (e: any) {
      ElMessage.error(e.message || '重启容器失败')
    } finally {
      state.actionLoading = false
    }
  }

  const removeContainer = async (id: string) => {
    try {
      state.actionLoading = true
      await apiRemoveContainer(id, { force: true })
      ElMessage.success('容器已删除')
      await loadData()
    } catch (e: any) {
      ElMessage.error(e.message || '删除容器失败')
    } finally {
      state.actionLoading = false
    }
  }

  const pagedContainers = computed(() => {
    const start = (state.page - 1) * state.pageSize
    return state.containers.slice(start, start + state.pageSize)
  })

  return {
    state,
    pagedContainers,
    loadData,
    startContainer,
    stopContainer,
    restartContainer,
    removeContainer
  }
}
