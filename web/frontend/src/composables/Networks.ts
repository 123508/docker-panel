import { reactive, computed, onMounted } from 'vue'
import {
  getNetworkList,
  getNetworkInspect,
  createNetwork as apiCreateNetwork,
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

  const createNetwork = async () => {
    let value: string
    try {
      const result = await ElMessageBox.prompt('请输入网络名称', '创建网络', {
        confirmButtonText: '创建',
        cancelButtonText: '取消',
        inputPattern: /\S+/,
        inputErrorMessage: '网络名称不能为空'
      })
      value = result.value
    } catch {
      return
    }

    const name = value.trim()
    await runWithDialog(
      {
        title: '创建网络',
        pendingText: `正在创建网络 ${name}，请稍候...`,
        successText: `✅ 网络创建成功: ${name}`,
        failureText: (e) => `❌ 创建网络失败: \n${e?.message || '未知错误'}`
      },
      async () => {
        await apiCreateNetwork({ name })
        await loadData()
      }
    )
  }

  const inspectNetwork = async (id: string) => {
    try {
      const detail = await getNetworkInspect(id)
      await ElMessageBox.alert(JSON.stringify(detail, null, 2), '网络详情', {
        confirmButtonText: '关闭'
      })
    } catch (e: any) {
      ElMessage.error(e.message || '获取网络详情失败')
    }
  }

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
    createNetwork,
    inspectNetwork,
    removeNetwork
  }
}
