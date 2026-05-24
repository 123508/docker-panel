import { reactive, computed, onMounted } from 'vue'
import {
  getComposePS,
  uploadComposeProject,
  upComposeProject,
  downComposeProject,
  stopComposeProject,
  restartComposeProject
} from '@/services/modules/compose'
import { call } from '@/services/common'
import { ElMessage, ElMessageBox } from 'element-plus'

interface ProjectItem {
  name: string
  file_path: string
  status: string
  services: number
  running: number
  created_at: string
  updated_at: string
}

export function ComposeProjectState() {
  const state = reactive({
    loading: false,
    projects: [] as ProjectItem[],
    stats: [
      { label: '项目:', value: '0' },
      { label: '运行中:', value: '0', variant: 'success' as const },
      { label: '已停止:', value: '0', variant: 'danger' as const }
    ]
  })

  const loadData = async () => {
    try {
      state.loading = true
      const res = await call('GET', '/api/v1/compose/projects', null).unwrap()
      const list = Array.isArray(res) ? res : []
      state.projects = list.map((p: any) => ({
        name: p.name || '-',
        file_path: p.file_path || '-',
        status: p.status || 'unknown',
        services: p.services?.length || 0,
        running: p.services?.filter((s: any) => s.status === 'running').length || 0,
        created_at: p.created_at || '',
        updated_at: p.updated_at || ''
      }))
      updateStats()
    } catch {
      state.projects = []
      updateStats()
    } finally {
      state.loading = false
    }
  }

  const updateStats = () => {
    const total = state.projects.length
    const running = state.projects.filter(p => p.status === 'running').length
    const stopped = state.projects.filter(p => p.status === 'stopped' || p.status === 'exited').length
    state.stats[0].value = String(total)
    state.stats[1].value = String(running)
    state.stats[2].value = String(stopped)
  }

  onMounted(() => { loadData() })

  const uploadNew = async () => {
    try {
      const result = await ElMessageBox.prompt('请输入项目名称', '新建编排', {
        confirmButtonText: '下一步',
        cancelButtonText: '取消',
        inputPattern: /\S+/,
        inputErrorMessage: '名称不能为空'
      })
      const name = result.value.trim()
      const content = await ElMessageBox.prompt('请输入 docker-compose.yml 内容', '编排内容', {
        confirmButtonText: '创建',
        cancelButtonText: '取消',
        inputType: 'textarea',
        inputPattern: /\S+/,
        inputErrorMessage: '内容不能为空'
      })
      await uploadComposeProject({ name, content: content.value })
      ElMessage.success(`编排项目 ${name} 创建成功`)
      await loadData()
    } catch (e: any) {
      if (!e?.toString().includes('cancel')) {
        ElMessage.error(e?.message || '操作失败')
      }
    }
  }

  const startProject = async (name: string) => {
    try {
      await upComposeProject(name)
      ElMessage.success(`项目 ${name} 已启动`)
      await loadData()
    } catch (e: any) {
      ElMessage.error(e?.message || `启动 ${name} 失败`)
    }
  }

  const stopProject = async (name: string) => {
    try {
      await stopComposeProject(name)
      ElMessage.success(`项目 ${name} 已停止`)
      await loadData()
    } catch (e: any) {
      ElMessage.error(e?.message || `停止 ${name} 失败`)
    }
  }

  const restartProject = async (name: string) => {
    try {
      await restartComposeProject(name)
      ElMessage.success(`项目 ${name} 已重启`)
      await loadData()
    } catch (e: any) {
      ElMessage.error(e?.message || `重启 ${name} 失败`)
    }
  }

  const removeProject = async (name: string) => {
    try {
      await ElMessageBox.confirm(`确认删除项目 ${name}？删除后将无法恢复。`, '删除确认', {
        confirmButtonText: '删除',
        cancelButtonText: '取消',
        type: 'warning'
      })
    } catch {
      return
    }
    try {
      await downComposeProject(name, { remove_volumes: false })
      ElMessage.success(`项目 ${name} 已删除`)
      await loadData()
    } catch (e: any) {
      ElMessage.error(e?.message || `删除 ${name} 失败`)
    }
  }

  const viewProject = async (name: string) => {
    try {
      const ps = await getComposePS(name)
      await ElMessageBox.alert(JSON.stringify(ps, null, 2), `项目 ${name} 状态`, {
        confirmButtonText: '关闭',
        type: 'info'
      })
    } catch (e: any) {
      ElMessage.error(e?.message || '获取状态失败')
    }
  }

  return {
    state,
    loadData,
    uploadNew,
    startProject,
    stopProject,
    restartProject,
    removeProject,
    viewProject
  }
}
