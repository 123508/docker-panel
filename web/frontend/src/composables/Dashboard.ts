import { reactive, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getContainerList } from '@/services/modules/container'
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
        const name = (c.names && c.names.length > 0)
          ? c.names[0].replace(/^\//, '')
          : (c.id?.substring(0, 12) || 'unknown')
        const running = c.state === 'running'

        return {
          name,
          id: c.id?.substring(0, 12) || 'N/A',
          image: c.image || 'N/A',
          status: running ? '运行中' : '已停止',
          running,
          port: c.ports
            ? c.ports.map((p: any) => `${p.public_port || p.private_port}:${p.private_port}`).join(', ')
            : ''
        }
      })

      const runningCount = containers.filter((c: any) => c.state === 'running').length
      const imagesTotalSize = images.reduce((sum: number, img: any) => sum + (img.size || 0), 0)

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
        (sum: number, detail: any) => sum + (detail?.usage_data?.size || 0),
        0
      )

      const customNetworks = networks.filter(
        (n: any) => !['bridge', 'host', 'none'].includes(n.name)
      ).length

      state.metrics = [
        { label: '容器', value: String(containers.length), sub: `${runningCount} 运行中`, subClass: 'success' },
        { label: '镜像', value: String(images.length), sub: formatBytesToGB(imagesTotalSize), subClass: 'muted' },
        { label: '卷', value: String(volumes.length), sub: `${formatBytesToGB(usedVolumeSize)} 已使用`, subClass: 'muted' },
        { label: '网络', value: String(networks.length), sub: `${customNetworks} 自定义`, subClass: 'muted' }
      ]
    } catch (e: any) {
      ElMessage.error(e.message || '获取 Dashboard 数据失败')
    }
  }

  onMounted(() => {
    loadData()
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
