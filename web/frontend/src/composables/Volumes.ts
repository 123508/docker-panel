import { reactive, computed, onMounted } from 'vue'
import { getVolumeList, getVolumeInspect, getVolumeContainers } from '@/services/modules/volume'
import { ElMessage } from 'element-plus'

function formatBytes(bytes: number): string {
  if (!Number.isFinite(bytes) || bytes <= 0) return '0 B'

  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  let value = bytes
  let unitIndex = 0

  while (value >= 1024 && unitIndex < units.length - 1) {
    value /= 1024
    unitIndex += 1
  }

  if (unitIndex === 0) {
    return `${Math.round(value)} ${units[unitIndex]}`
  }

  return `${value.toFixed(2)} ${units[unitIndex]}`
}

export function VolumeState() {
  const state = reactive({
    page: 1,
    pageSize: 5,

    stats: [
      { label: '总数:', value: '0' },
      { label: '已使用:', value: '0 B' },
      { label: '未使用:', value: '0', variant: 'dangling' }
    ],

    volumes: [] as any[]
  })

  const loadData = async () => {
    try {
      const res = await getVolumeList()
      if (!Array.isArray(res)) {
        state.volumes = []
        state.stats[0].value = '0'
        state.stats[1].value = '0 B'
        state.stats[2].value = '0'
        return
      }

      const extraData = await Promise.all(
        res.map(async (v: any) => {
          try {
            const [inspect, containersRes] = await Promise.all([
              getVolumeInspect(v.name),
              getVolumeContainers(v.name)
            ])

            const rawSize = typeof inspect?.usage_data?.size === 'number' ? inspect.usage_data.size : 0
            const containersCount = Array.isArray(containersRes?.containers)
              ? containersRes.containers.length
              : 0

            return {
              name: v.name,
              rawSize,
              size: rawSize > 0 ? formatBytes(rawSize) : '-',
              containers: containersCount
            }
          } catch {
            return {
              name: v.name,
              rawSize: 0,
              size: '-',
              containers: 0
            }
          }
        })
      )

      const extraMap = new Map(extraData.map(item => [item.name, item]))

      state.volumes = res.map((v: any) => {
        const extra = extraMap.get(v.name)
        return {
          name: v.name || 'N/A',
          driver: v.driver || 'N/A',
          mountpoint: v.mountpoint || '-',
          size: extra?.size || '-',
          containers: extra?.containers ?? 0
        }
      })

      const total = res.length
      const usedSize = extraData.reduce((sum, item) => sum + item.rawSize, 0)
      const dangling = extraData.filter(item => item.containers === 0).length

      state.stats[0].value = String(total)
      state.stats[1].value = formatBytes(usedSize)
      state.stats[2].value = String(dangling)
    } catch (e: any) {
      ElMessage.error(e.message || '获取数据卷列表失败')
    }
  }

  onMounted(() => {
    loadData()
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
