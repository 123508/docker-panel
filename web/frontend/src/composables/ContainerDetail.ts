import { reactive, computed, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  getContainerDetail,
  startContainerById,
  stopContainerById,
  restartContainerById,
  removeContainerById
} from '@/services/modules/container-detail'

export function ContainerDetailState(containerIdRef: () => string) {
  const state = reactive({
    loading: false,
    actionLoading: false,
    detail: null as any | null
  })

  const isRunning = computed(() => !!state.detail?.state?.running)
  const displayName = computed(() => {
    const name = state.detail?.name || ''
    return name.replace(/^\//, '') || 'N/A'
  })

  const loadData = async () => {
    const id = containerIdRef()
    if (!id) return
    try {
      state.loading = true
      state.detail = await getContainerDetail(id)
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

  const remove = async (): Promise<boolean> => {
    const id = containerIdRef()
    if (!id) return false
    try {
      await ElMessageBox.confirm('确认删除该容器？', '删除确认', {
        confirmButtonText: '删除',
        cancelButtonText: '取消',
        type: 'warning'
      })
      state.actionLoading = true
      await removeContainerById(id)
      ElMessage.success('容器已删除')
      return true
    } catch (e: any) {
      if (e !== 'cancel' && e !== 'close') {
        ElMessage.error(e.message || '删除容器失败')
      }
      return false
    } finally {
      state.actionLoading = false
    }
  }

  watch(
    () => containerIdRef(),
    () => {
      loadData()
    },
    { immediate: true }
  )

  return {
    state,
    isRunning,
    displayName,
    loadData,
    start,
    stop,
    restart,
    remove
  }
}
