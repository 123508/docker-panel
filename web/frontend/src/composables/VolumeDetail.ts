import { reactive, computed, onMounted, ref, watch } from 'vue'
import {
  getVolumeInspect,
  getVolumeContainers,
  removeVolume
} from '@/services/modules/volume'
import { ElMessage, ElMessageBox } from 'element-plus'

export function VolumeDetailState(getName: () => string) {
  const state = reactive({
    loading: false,
    detail: null as any,
    containers: [] as string[]
  })

  const loadData = async () => {
    const name = getName()
    if (!name) return
    try {
      state.loading = true
      const [detail, containersRes] = await Promise.all([
        getVolumeInspect(name),
        getVolumeContainers(name).catch(() => ({ containers: [] }))
      ])
      state.detail = detail
      state.containers = containersRes?.containers || []
    } catch (e: any) {
      ElMessage.error(e.message || '获取卷详情失败')
    } finally {
      state.loading = false
    }
  }

  const volumeName = computed(() => state.detail?.name || getName() || '')
  const driver = computed(() => state.detail?.driver || '-')
  const mountpoint = computed(() => state.detail?.mountpoint || '-')
  const createdAt = computed(() => {
    if (!state.detail?.created_at) return '-'
    return new Date(state.detail.created_at).toLocaleString()
  })
  const scope = computed(() => state.detail?.scope || '-')
  const size = computed(() => {
    const raw = state.detail?.usage_data?.size
    if (typeof raw !== 'number' || raw <= 0) return '-'
    const mb = raw / 1024 / 1024
    if (mb < 1024) return mb.toFixed(1) + ' MB'
    return (mb / 1024).toFixed(2) + ' GB'
  })
  const refCount = computed(() => state.detail?.usage_data?.ref_count ?? 0)

  const labelsList = computed(() => {
    const labels = state.detail?.labels
    if (!labels || typeof labels !== 'object') return []
    return Object.entries(labels).map(([key, value]) => ({ key, value }))
  })

  const optionsList = computed(() => {
    const opts = state.detail?.options
    if (!opts || typeof opts !== 'object') return []
    return Object.entries(opts).map(([key, value]) => ({ key, value }))
  })

  const handleRemove = async () => {
    const name = volumeName.value
    if (!name) return
    try {
      await ElMessageBox.confirm('确认删除该卷？删除后将无法恢复。', '删除确认', {
        confirmButtonText: '删除',
        cancelButtonText: '取消',
        type: 'warning'
      })
    } catch {
      return false
    }
    try {
      await removeVolume(name, { force: true })
      ElMessage.success(`卷已删除: ${name}`)
      return true
    } catch (e: any) {
      ElMessage.error(e.message || '删除卷失败')
      return false
    }
  }

  return {
    state,
    volumeName,
    driver,
    mountpoint,
    createdAt,
    scope,
    size,
    refCount,
    containers: computed(() => state.containers),
    labelsList,
    optionsList,
    loadData,
    handleRemove
  }
}
