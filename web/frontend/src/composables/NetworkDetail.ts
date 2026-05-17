import { reactive, computed } from 'vue'
import {
  getNetworkInspect,
  removeNetwork
} from '@/services/modules/network'
import { ElMessage, ElMessageBox } from 'element-plus'

export function NetworkDetailState(getId: () => string) {
  const state = reactive({
    loading: false,
    detail: null as any
  })

  const loadData = async () => {
    const id = getId()
    if (!id) return
    try {
      state.loading = true
      state.detail = await getNetworkInspect(id)
    } catch (e: any) {
      ElMessage.error(e.message || '获取网络详情失败')
    } finally {
      state.loading = false
    }
  }

  const networkName = computed(() => state.detail?.name || getId() || '-')
  const shortId = computed(() => {
    const id = state.detail?.id || ''
    return id.substring(0, 12)
  })
  const fullId = computed(() => state.detail?.id || '-')
  const driver = computed(() => state.detail?.driver || '-')
  const scope = computed(() => state.detail?.scope || '-')
  const created = computed(() => {
    if (!state.detail?.created) return '-'
    return new Date(state.detail.created).toLocaleString()
  })
  const enableIPv6 = computed(() => !!state.detail?.enable_ipv6)
  const internal = computed(() => !!state.detail?.internal)
  const attachable = computed(() => !!state.detail?.attachable)
  const ingress = computed(() => !!state.detail?.ingress)

  const ipamDriver = computed(() => state.detail?.ipam?.driver || '-')
  const ipamConfigs = computed(() => {
    const configs = state.detail?.ipam?.config
    if (!Array.isArray(configs)) return []
    return configs.map((c: any) => ({
      subnet: c.subnet || '-',
      gateway: c.gateway || '-',
      ipRange: c.ip_range || '-'
    }))
  })

  const containers = computed(() => {
    const cons = state.detail?.containers
    if (!cons || typeof cons !== 'object') return []
    return Object.entries(cons).map(([id, info]: [string, any]) => ({
      id,
      name: info.name || id.substring(0, 12),
      ipv4: info.ipv4_address || '-',
      ipv6: info.ipv6_address || '-',
      mac: info.mac_address || '-'
    }))
  })

  const labelsList = computed(() => {
    const labels = state.detail?.labels
    if (!labels || typeof labels !== 'object') return []
    return Object.entries(labels).map(([key, value]) => `${key}=${value}`)
  })

  const optionsList = computed(() => {
    const opts = state.detail?.options
    if (!opts || typeof opts !== 'object') return []
    return Object.entries(opts).map(([key, value]) => `${key}=${value}`)
  })

  const handleRemove = async () => {
    const id = getId()
    if (!id) return false
    try {
      await ElMessageBox.confirm('确认删除该网络？', '删除确认', {
        confirmButtonText: '删除',
        cancelButtonText: '取消',
        type: 'warning'
      })
    } catch {
      return false
    }
    try {
      await removeNetwork(id)
      ElMessage.success('网络已删除')
      return true
    } catch (e: any) {
      ElMessage.error(e.message || '删除网络失败')
      return false
    }
  }

  return {
    state,
    networkName,
    shortId,
    fullId,
    driver,
    scope,
    created,
    enableIPv6,
    internal,
    attachable,
    ingress,
    ipamDriver,
    ipamConfigs,
    containers,
    labelsList,
    optionsList,
    loadData,
    handleRemove
  }
}
