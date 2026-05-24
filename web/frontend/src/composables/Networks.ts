import { reactive, computed, onMounted } from 'vue'
import {
  getNetworkList,
  getNetworkInspect,
  removeNetwork as apiRemoveNetwork
} from '@/services/modules/network'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useActionDialog } from '@/composables/useActionDialog'

const BUILTIN_NETWORKS = ['bridge', 'host', 'none']

export function NetworkState() {
  const state = reactive({
    page: 1,
    pageSize: 5,

    stats: [
      { label: '总数:', value: '0' },
      { label: '自定义:', value: '0' },
      { label: '内置:', value: '0', variant: 'builtin' }
    ],

    networks: [] as any[]
  })

  const { dialog, runWithDialog } = useActionDialog()

  const loadData = async () => {
    try {
      const res = await getNetworkList()
      if (!Array.isArray(res)) {
        state.networks = []
        state.stats[0].value = '0'
        state.stats[1].value = '0'
        state.stats[2].value = '0'
        return
      }

      const details = await Promise.all(
        res.map(async (n: any) => {
          try {
            const inspect = await getNetworkInspect(n.id)
            return { id: n.id, inspect }
          } catch {
            return { id: n.id, inspect: null }
          }
        })
      )

      const inspectMap = new Map(details.map(item => [item.id, item.inspect]))

      state.networks = res.map((n: any) => {
        const inspect = inspectMap.get(n.id)
        const firstConfig = n.ipam?.config?.[0]
        const subnet = firstConfig?.subnet || '-'
        const gateway = firstConfig?.gateway || '-'
        const containers = inspect?.containers ? Object.keys(inspect.containers).length : 0
        const isBuiltin = BUILTIN_NETWORKS.includes(n.name)

        return {
          id: n.id,
          name: n.name || 'N/A',
          driver: n.driver || 'N/A',
          scope: n.scope || 'N/A',
          subnet,
          gateway,
          containers,
          removable: isBuiltin
        }
      })

      const total = res.length
      const builtin = res.filter((n: any) => BUILTIN_NETWORKS.includes(n.name)).length
      const custom = total - builtin

      state.stats[0].value = String(total)
      state.stats[1].value = String(custom)
      state.stats[2].value = String(builtin)
    } catch (e: any) {
      ElMessage.error(e.message || '获取网络列表失败')
    }
  }

  onMounted(() => {
    loadData()
  })

  const removeNetwork = async (id: string) => {
    try {
      await ElMessageBox.confirm('确认删除该网络？', '删除确认', {
        confirmButtonText: '删除',
        cancelButtonText: '取消',
        type: 'warning'
      })
    } catch {
      return
    }

    const shortId = id.substring(0, 12)
    await runWithDialog(
      {
        title: '移除网络',
        pendingText: `正在移除网络 ${shortId}，请稍候...`,
        successText: `✅ 网络已移除: ${shortId}`,
        failureText: (e) => `❌ 移除网络失败: \n${e?.message || '未知错误'}`
      },
      async () => {
        await apiRemoveNetwork(id)
        await loadData()
      }
    )
  }

  const pagedNetworks = computed(() => {
    const start = (state.page - 1) * state.pageSize
    return state.networks.slice(start, start + state.pageSize)
  })

  return {
    state,
    dialog,
    pagedNetworks,
    removeNetwork
  }
}
