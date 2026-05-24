import { computed, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import {
  createNewContainer,
  getCreateContainerImageOptions,
  getCreateContainerNetworkOptions,
  getCreateContainerVolumeOptions
} from '@/services/modules/container-create'
import type { ContainerCreateRequest, MountRequest } from '@/services/modules/container'

type RestartPolicyName = 'no' | 'always' | 'on-failure' | 'unless-stopped'

interface PortRow {
  containerPort: string
  hostPort: string
}

interface EnvRow {
  key: string
  value: string
}

interface VolumeRow {
  sourceMode: 'volume' | 'custom'
  source: string
  target: string
}

function parseCmd(text: string): string[] {
  return text
    .split(' ')
    .map((s) => s.trim())
    .filter(Boolean)
}

function isLikelyHostPath(source: string): boolean {
  return source.includes('/') || source.includes('\\') || /^[A-Za-z]:/.test(source)
}

export function ContainerCreateState() {
  const state = reactive({
    loading: false,
    creating: false,
    imageOptions: [] as string[],
    networkOptions: [] as string[],
    volumeOptions: [] as string[],
    form: {
      name: '',
      image: '',
      cmd: '',
      networkMode: 'bridge',
      restartPolicy: 'unless-stopped' as RestartPolicyName,
      ports: [{ containerPort: '', hostPort: '' }] as PortRow[],
      envs: [{ key: '', value: '' }] as EnvRow[],
      volumes: [{ sourceMode: 'volume', source: '', target: '' }] as VolumeRow[]
    }
  })

  const canSubmit = computed(() => !!state.form.name && !!state.form.image)

  const addPortRow = () => {
    state.form.ports.push({ containerPort: '', hostPort: '' })
  }

  const removePortRow = (index: number) => {
    state.form.ports.splice(index, 1)
    if (state.form.ports.length === 0) {
      addPortRow()
    }
  }

  const addEnvRow = () => {
    state.form.envs.push({ key: '', value: '' })
  }

  const removeEnvRow = (index: number) => {
    state.form.envs.splice(index, 1)
    if (state.form.envs.length === 0) {
      addEnvRow()
    }
  }

  const addVolumeRow = () => {
    state.form.volumes.push({ sourceMode: 'volume', source: '', target: '' })
  }

  const removeVolumeRow = (index: number) => {
    state.form.volumes.splice(index, 1)
    if (state.form.volumes.length === 0) {
      addVolumeRow()
    }
  }

  const loadOptions = async () => {
    try {
      state.loading = true
      const [images, networks, volumes] = await Promise.all([
        getCreateContainerImageOptions(),
        getCreateContainerNetworkOptions(),
        getCreateContainerVolumeOptions()
      ])

      state.imageOptions = images
        .map((img: any) => {
          if (Array.isArray(img.repo_tags) && img.repo_tags.length > 0) return img.repo_tags[0]
          return ''
        })
        .filter(Boolean)

      state.networkOptions = networks.map((n: any) => n.name).filter(Boolean)
      state.volumeOptions = volumes.map((v: any) => v.name).filter(Boolean)
      if (state.volumeOptions.length > 0) {
        state.form.volumes.forEach((row) => {
          if (row.sourceMode === 'volume' && !row.source) {
            row.source = state.volumeOptions[0]
          }
        })
      }

      if (!state.form.image && state.imageOptions.length > 0) {
        state.form.image = state.imageOptions[0]
      }
      if (state.networkOptions.length > 0 && !state.networkOptions.includes(state.form.networkMode)) {
        state.form.networkMode = state.networkOptions[0]
      }
    } catch (e: any) {
      ElMessage.error(e.message || '加载创建选项失败')
    } finally {
      state.loading = false
    }
  }

  const submit = async (): Promise<string | null> => {
    if (!canSubmit.value) {
      ElMessage.warning('请填写必填项')
      return null
    }

    try {
      state.creating = true

      const envList = state.form.envs
        .filter((row) => row.key.trim())
        .map((row) => `${row.key.trim()}=${row.value.trim()}`)

      const portBindings = state.form.ports.reduce<Record<string, Array<{ host_ip?: string; host_port: string }>>>(
        (acc, row) => {
          const containerPort = row.containerPort.trim()
          const hostPort = row.hostPort.trim()
          if (!containerPort || !hostPort) return acc
          acc[`${containerPort}/tcp`] = [{ host_port: hostPort }]
          return acc
        },
        {}
      )

      const mounts = state.form.volumes.reduce<MountRequest[]>((acc, row) => {
        const source = row.source.trim()
        const target = row.target.trim()
        if (!source || !target) return acc
        acc.push({
          type: isLikelyHostPath(source) ? 'bind' : 'volume',
          source,
          target
        })
        return acc
      }, [])

      const payload: ContainerCreateRequest = {
        name: state.form.name.trim(),
        image: state.form.image.trim(),
        cmd: state.form.cmd.trim() ? parseCmd(state.form.cmd) : undefined,
        env: envList.length ? envList : undefined,
        host_config: {
          network_mode: state.form.networkMode,
          restart_policy: {
            name: state.form.restartPolicy,
            maximum_retry_count: 0
          },
          port_bindings: Object.keys(portBindings).length ? portBindings : undefined,
          mounts: mounts.length ? mounts : undefined
        },
        networking_config: undefined
      }

      const res = await createNewContainer(payload)
      ElMessage.success('容器创建成功')
      return res?.id || state.form.name
    } catch (e: any) {
      ElMessage.error(e.message || '创建容器失败')
      return null
    } finally {
      state.creating = false
    }
  }

  onMounted(() => {
    loadOptions()
  })

  return {
    state,
    canSubmit,
    submit,
    addPortRow,
    removePortRow,
    addEnvRow,
    removeEnvRow,
    addVolumeRow,
    removeVolumeRow
  }
}
