import { reactive, computed } from 'vue'
import {
  getImageInspect,
  removeImage
} from '@/services/modules/image'
import { ElMessage, ElMessageBox } from 'element-plus'

export function ImageDetailState(getId: () => string) {
  const state = reactive({
    loading: false,
    detail: null as any
  })

  const loadData = async () => {
    const id = getId()
    if (!id) return
    try {
      state.loading = true
      state.detail = await getImageInspect(id)
    } catch (e: any) {
      ElMessage.error(e.message || '获取镜像详情失败')
    } finally {
      state.loading = false
    }
  }

  const imageName = computed(() => {
    const tags = state.detail?.repo_tags
    if (tags && tags.length > 0 && tags[0] !== '<none>:<none>') return tags[0]
    return state.detail?.id?.replace('sha256:', '').substring(0, 12) || '-'
  })

  const shortId = computed(() => {
    const id = state.detail?.id || ''
    return id.replace('sha256:', '').substring(0, 12)
  })

  const fullId = computed(() => state.detail?.id || '-')

  const repo = computed(() => {
    const tags = state.detail?.repo_tags
    if (tags && tags.length > 0 && tags[0] !== '<none>:<none>') {
      const parts = tags[0].split(':')
      return parts[0]
    }
    const digests = state.detail?.repo_digests
    if (digests && digests.length > 0) return digests[0].split('@')[0]
    return '-'
  })

  const tag = computed(() => {
    const tags = state.detail?.repo_tags
    if (tags && tags.length > 0 && tags[0] !== '<none>:<none>') {
      return tags[0].split(':').slice(1).join(':') || 'latest'
    }
    return '-'
  })

  const created = computed(() => {
    if (!state.detail?.created) return '-'
    return new Date(state.detail.created).toLocaleString()
  })

  const size = computed(() => {
    const raw = typeof state.detail?.size === 'number' ? state.detail.size : 0
    if (raw <= 0) return '-'
    const mb = raw / 1024 / 1024
    if (mb < 1024) return mb.toFixed(1) + ' MB'
    return (mb / 1024).toFixed(2) + ' GB'
  })

  const os = computed(() => state.detail?.os || '-')
  const architecture = computed(() => state.detail?.architecture || '-')
  const dockerVersion = computed(() => state.detail?.docker_version || '-')

  const layerCount = computed(() => {
    const rootFs = state.detail?.root_fs
    return rootFs?.layers?.length ?? 0
  })

  const containers = computed(() => {
    const c = state.detail?.containers
    return typeof c === 'number' ? c : 0
  })

  const envVars = computed(() => {
    const config = state.detail?.config
    if (!config?.env || !Array.isArray(config.env)) return []
    return config.env.map((e: string) => {
      const idx = e.indexOf('=')
      if (idx === -1) return { key: e, value: '' }
      return { key: e.substring(0, idx + 1), value: e.substring(idx + 1) }
    })
  })

  const exposedPorts = computed(() => {
    const config = state.detail?.config
    const ports = config?.exposed_ports || config?.ExposedPorts
    if (!ports || typeof ports !== 'object') return []
    return Object.keys(ports).map(p => {
      const [port, proto] = p.split('/')
      return { port, protocol: proto || 'tcp' }
    })
  })

  const digest = computed(() => {
    const digests = state.detail?.repo_digests
    if (digests && digests.length > 0) {
      const d = digests[0]
      const idx = d.indexOf('@')
      return idx >= 0 ? d.substring(idx + 1) : d
    }
    return '-'
  })

  const handleRemove = async () => {
    const id = getId()
    if (!id) return false
    try {
      await ElMessageBox.confirm('确认删除该镜像？', '删除确认', {
        confirmButtonText: '删除',
        cancelButtonText: '取消',
        type: 'warning'
      })
    } catch {
      return false
    }
    try {
      await removeImage(id, { force: true })
      ElMessage.success('镜像已删除')
      return true
    } catch (e: any) {
      ElMessage.error(e.message || '删除镜像失败')
      return false
    }
  }

  return {
    state,
    imageName,
    shortId,
    fullId,
    repo,
    tag,
    created,
    size,
    os,
    architecture,
    dockerVersion,
    layerCount,
    containers,
    envVars,
    exposedPorts,
    digest,
    loadData,
    handleRemove
  }
}
