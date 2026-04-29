import { computed, onBeforeUnmount, reactive, watch } from 'vue'
import { ElMessage } from 'element-plus'
import {
  createContainerLogsWebSocket,
  getContainerDetailPage,
  restartContainerById,
  startContainerById,
  stopContainerById
} from '@/services/modules/container-detail'

export function ContainerDetailState(containerIdRef: () => string) {
  const state = reactive({
    loading: false,
    actionLoading: false,
    detail: null as any | null,
    logs: ''
  })

  let logsWS: WebSocket | null = null
  let reconnectTimer: number | null = null
  let shouldReconnectLogs = false

  const isRunning = computed(() => !!state.detail?.state?.running)

  const displayName = computed(() => {
    const name = state.detail?.name || ''
    return name.replace(/^\//, '') || 'N/A'
  })

  const shortId = computed(() => {
    const id = state.detail?.id || containerIdRef()
    return String(id || '').slice(0, 12) || 'N/A'
  })

  const uptimeText = computed(() => {
    const startedAt = state.detail?.state?.started_at
    if (!startedAt || !isRunning.value) return 'N/A'
    const started = new Date(startedAt).getTime()
    if (!Number.isFinite(started) || started <= 0) return 'N/A'
    const diffMinutes = Math.max(0, Math.floor((Date.now() - started) / 60000))
    const hours = Math.floor(diffMinutes / 60)
    const minutes = diffMinutes % 60
    if (hours <= 0) return `${minutes} 分钟`
    return `${hours} 小时 ${minutes} 分钟`
  })

  const imageText = computed(() => state.detail?.config?.image || state.detail?.image || 'N/A')

  const commandText = computed(() => {
    const detail = state.detail || {}
    const parts = [...(detail.config?.entrypoint || []), ...(detail.config?.cmd || detail.args || [])]
    return parts.length ? parts.join(' ') : detail.path || 'N/A'
  })

  const cpuText = computed(() => {
    const nanoCPUs = Number(state.detail?.host_config?.nano_cpus || 0)
    if (!nanoCPUs) return '--'
    return `${(nanoCPUs / 1_000_000_000).toFixed(2)} CPU`
  })

  const formatBytes = (bytes: number) => {
    if (!bytes) return 'N/A'
    if (bytes < 1024 * 1024 * 1024) return `${Math.round(bytes / (1024 * 1024))}MB`
    return `${(bytes / (1024 * 1024 * 1024)).toFixed(1)}GB`
  }

  const memoryText = computed(() => {
    const memory = Number(state.detail?.host_config?.memory || 0)
    const swap = Number(state.detail?.host_config?.memory_swap || 0)
    if (!memory) return '不限'
    return swap > 0 && swap !== memory ? `${formatBytes(memory)} / ${formatBytes(swap)}` : formatBytes(memory)
  })

  const portRows = computed(() => {
    const ports = state.detail?.network_settings?.ports || state.detail?.host_config?.port_bindings || {}
    return Object.entries(ports).flatMap(([portKey, bindings]: [string, any]) => {
      const [privatePort, type = 'tcp'] = portKey.split('/')
      const rows = Array.isArray(bindings) && bindings.length ? bindings : [{ private_port: privatePort, type }]
      return rows.map((binding: any, index: number) => ({
        key: `${portKey}-${index}-${binding.ip || ''}-${binding.public_port || ''}`,
        private: `${binding.private_port || privatePort}/${binding.type || type}`,
        public: binding.public_port ? `${binding.ip || '0.0.0.0'}:${binding.public_port}` : '未发布',
        type: binding.type || type
      }))
    })
  })

  const envRows = computed(() => {
    const env = state.detail?.config?.env || []
    return env.slice(0, 8).map((raw: string) => {
      const index = raw.indexOf('=')
      if (index < 0) return { raw, key: `${raw}=`, value: '' }
      return { raw, key: raw.slice(0, index + 1), value: raw.slice(index + 1) }
    })
  })

  const mountRows = computed(() => {
    const mounts = state.detail?.mounts || []
    return mounts.map((mount: any, index: number) => ({
      key: `${mount.source || mount.name || index}-${mount.destination || index}`,
      source: mount.source || mount.name || mount.type || 'N/A',
      destination: mount.destination || 'N/A',
      mode: mount.rw ? 'rw' : 'ro'
    }))
  })

  const networkRows = computed(() => {
    const networks = state.detail?.network_settings?.networks || {}
    return Object.entries(networks).map(([name, network]: [string, any]) => ({
      name,
      ip: network?.ip_address || 'N/A',
      gateway: network?.gateway || 'N/A'
    }))
  })

  const networkSummary = computed(() => {
    if (!networkRows.value.length) return 'N/A'
    return networkRows.value.map((network) => network.name).join(', ')
  })

  const restartPolicyText = computed(() => {
    const policy = state.detail?.host_config?.restart_policy
    if (!policy?.name || policy.name === 'no') return 'no'
    if (policy.maximum_retry_count) return `${policy.name} (${policy.maximum_retry_count})`
    return policy.name
  })

  const logLines = computed(() => {
    const rawLines = String(state.logs || '')
      .split(/\r?\n/)
      .map((line) => line.trim())
      .filter(Boolean)
      .slice(-100)
    const lines = rawLines.length ? rawLines : ['[info] 暂无容器日志']
    return lines.map((text) => ({
      text,
      levelClass: /\b(warn|warning)\b/i.test(text) ? 'log-warn' : 'log-info'
    }))
  })

  const closeLogsWS = () => {
    if (reconnectTimer != null) {
      window.clearTimeout(reconnectTimer)
      reconnectTimer = null
    }
    if (logsWS) {
      logsWS.onopen = null
      logsWS.onmessage = null
      logsWS.onclose = null
      logsWS.onerror = null
      logsWS.close()
      logsWS = null
    }
  }

  const connectLogsWS = () => {
    const id = containerIdRef()
    if (!id) return

    closeLogsWS()
    shouldReconnectLogs = true

    const ws = createContainerLogsWebSocket(id, { tail: '100', timestamps: true })
    logsWS = ws

    ws.onmessage = (event) => {
      const text = typeof event.data === 'string' ? event.data : ''
      if (!text) return

      state.logs = `${state.logs}${state.logs.endsWith('\n') || state.logs === '' ? '' : '\n'}${text}`
      const lines = state.logs.split(/\r?\n/)
      if (lines.length > 1200) {
        state.logs = lines.slice(-1000).join('\n')
      }
    }

    ws.onclose = () => {
      if (!shouldReconnectLogs) return
      reconnectTimer = window.setTimeout(() => {
        connectLogsWS()
      }, 3000)
    }

    ws.onerror = () => {
      ws.close()
    }
  }

  const loadData = async () => {
    const id = containerIdRef()
    if (!id) return
    try {
      state.loading = true
      const page = await getContainerDetailPage(id)
      state.detail = page.detail
      state.logs = page.logs
      connectLogsWS()
    } catch (e: any) {
      ElMessage.error(e.message || '获取容器详情失败')
    } finally {
      state.loading = false
    }
  }

  const start = async () => {
    const id = containerIdRef()
    if (!id) return
    try {
      state.actionLoading = true
      await startContainerById(id)
      ElMessage.success('容器已启动')
      await loadData()
    } catch (e: any) {
      ElMessage.error(e.message || '启动容器失败')
    } finally {
      state.actionLoading = false
    }
  }

  const stop = async () => {
    const id = containerIdRef()
    if (!id) return
    try {
      state.actionLoading = true
      await stopContainerById(id)
      ElMessage.success('容器已停止')
      await loadData()
    } catch (e: any) {
      ElMessage.error(e.message || '停止容器失败')
    } finally {
      state.actionLoading = false
    }
  }

  const restart = async () => {
    const id = containerIdRef()
    if (!id) return
    try {
      state.actionLoading = true
      await restartContainerById(id)
      ElMessage.success('容器已重启')
      await loadData()
    } catch (e: any) {
      ElMessage.error(e.message || '重启容器失败')
    } finally {
      state.actionLoading = false
    }
  }

  watch(
    () => containerIdRef(),
    (newId, oldId) => {
      if (newId !== oldId) {
        shouldReconnectLogs = false
        closeLogsWS()
      }
      loadData()
    },
    { immediate: true }
  )

  onBeforeUnmount(() => {
    shouldReconnectLogs = false
    closeLogsWS()
  })

  return {
    state,
    isRunning,
    displayName,
    shortId,
    uptimeText,
    imageText,
    commandText,
    cpuText,
    memoryText,
    portRows,
    envRows,
    mountRows,
    networkRows,
    networkSummary,
    restartPolicyText,
    logLines,
    loadData,
    start,
    stop,
    restart
  }
}
