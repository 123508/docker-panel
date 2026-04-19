import { computed, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import {
  createNewContainer,
  getCreateContainerImageOptions,
  getCreateContainerNetworkOptions
} from '@/services/modules/container-create'

function parseByLine(text: string): string[] {
  return text
    .split('\n')
    .map(s => s.trim())
    .filter(Boolean)
}

function parseCmd(text: string): string[] {
  return text
    .split(' ')
    .map(s => s.trim())
    .filter(Boolean)
}

export function ContainerCreateState() {
  const state = reactive({
    loading: false,
    creating: false,
    imageOptions: [] as string[],
    networkOptions: [] as string[],
    form: {
      name: '',
      image: '',
      cmd: '',
      envText: '',
      networkMode: 'bridge',
      restartPolicy: 'unless-stopped',
      portHost: '',
      portContainer: ''
    }
  })

  const canSubmit = computed(() => !!state.form.name && !!state.form.image)

  const loadOptions = async () => {
    try {
      state.loading = true
      const [images, networks] = await Promise.all([
        getCreateContainerImageOptions(),
        getCreateContainerNetworkOptions()
      ])

      state.imageOptions = images
        .map((img: any) => {
          if (Array.isArray(img.repo_tags) && img.repo_tags.length > 0) return img.repo_tags[0]
          return ''
        })
        .filter(Boolean)

      state.networkOptions = networks.map((n: any) => n.name).filter(Boolean)
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
      const portBindings = (state.form.portHost && state.form.portContainer)
        ? {
            [`${state.form.portContainer}/tcp`]: [{ host_port: state.form.portHost }]
          }
        : undefined

      const payload: any = {
        name: state.form.name,
        image: state.form.image,
        cmd: state.form.cmd ? parseCmd(state.form.cmd) : undefined,
        env: state.form.envText ? parseByLine(state.form.envText) : undefined,
        host_config: {
          network_mode: state.form.networkMode,
          restart_policy: {
            name: state.form.restartPolicy,
            maximum_retry_count: 0
          },
          port_bindings: portBindings
        }
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
    submit
  }
}
