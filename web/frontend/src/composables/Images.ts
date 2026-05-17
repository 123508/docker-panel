import { reactive, computed, onMounted } from 'vue'
import { getImageList, pullImage as apiPullImage, removeImage as apiRemoveImage } from '@/services/modules/image'
import { createContainer, startContainer } from '@/services/modules/container'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useActionDialog } from '@/composables/useActionDialog'

export function ImageState() {
  const state = reactive({
    query: '',
    page: 1,
    pageSize: 5,
    loading: false,

    stats: [
      { label: '总数:', value: '0' },
      { label: '大小:', value: '0 B' },
      { label: '悬空:', value: '0', variant: 'dangling' }
    ],

    images: [] as any[]
  })

  const { dialog, runWithDialog } = useActionDialog()

  const loadData = async () => {
    try {
      state.loading = true
      const res = await getImageList({ all: true })
      if (res) {
        state.images = res.map((img: any) => {
          let repo = '<none>'
          let tag = '<none>'
          if (img.repo_tags && img.repo_tags.length > 0) {
             const parts = img.repo_tags[0].split(':')
             repo = parts[0]
             tag = parts[1] || 'latest'
          } else if (img.repo_digests && img.repo_digests.length > 0) {
             repo = img.repo_digests[0].split('@')[0]
             tag = '<none>'
          }

          return {
            repo,
            tag,
            fullId: img.id,
            id: img.id.replace('sha256:', '').substring(0, 12),
            size: (img.size / 1024 / 1024).toFixed(2) + ' MB',
            rawSize: img.size,
            created: new Date(img.created * 1000).toLocaleString(),
            color: '#3B82F6',
            dangling: repo === '<none>' && tag === '<none>'
          }
        })

        state.stats[0].value = state.images.length.toString()
        const totalSize = state.images.reduce((acc, img) => acc + (img.rawSize || 0), 0)
        state.stats[1].value = (totalSize / 1024 / 1024 / 1024).toFixed(2) + ' GB'
        state.stats[2].value = state.images.filter(img => img.dangling).length.toString()
      }
    } catch (e: any) {
      ElMessage.error(e.message || '获取镜像列表失败')
    } finally {
      state.loading = false
    }
  }

  onMounted(() => {
    loadData()
  })

  const removeImage = async (id: string) => {
    const shortId = id.replace('sha256:', '').substring(0, 12)
    await runWithDialog(
      {
        title: '移除镜像',
        pendingText: `正在移除镜像 ${shortId}，请稍候...`,
        successText: `✅ 镜像已移除: ${shortId}`,
        failureText: (e) => `❌ 移除镜像失败: \n${e?.message || '未知错误'}`
      },
      async () => {
        await apiRemoveImage(id, { force: true })
        await loadData()
      }
    )
  }

  const searchImages = async () => {
    await loadData()
  }

  const pullImage = async () => {
    let value: string
    try {
      const result = await ElMessageBox.prompt('请输入镜像名（如 nginx:latest）', '拉取镜像', {
        confirmButtonText: '拉取',
        cancelButtonText: '取消',
        inputPattern: /\S+/,
        inputErrorMessage: '镜像名不能为空'
      })
      value = result.value
    } catch (e: any) {
      // 用户取消
      return
    }

    const image = value.trim()
    await runWithDialog(
      {
        title: '拉取镜像',
        pendingText: `正在拉取镜像 ${image}，请稍候...`,
        successText: `✅ 镜像拉取成功: ${image}`,
        failureText: (e) => `❌ 拉取镜像失败: \n${e?.message || '未知错误'}`
      },
      async () => {
        await apiPullImage({ image })
        await loadData()
      }
    )
  }

  const runImage = async (imageRef: string, imageName?: string) => {
    const display = imageName || imageRef
    await runWithDialog(
      {
        title: '运行镜像',
        pendingText: `正在启动镜像 ${display}，请稍候...`,
        failureText: (e) => `❌ 运行镜像失败: \n${e?.message || '未知错误'}`
      },
      async () => {
        const suffix = Date.now().toString(36).slice(-6)
        const cleanName = (imageName || 'image').replace(/[^a-zA-Z0-9_.-]/g, '-').toLowerCase()
        const name = `run-${cleanName}-${suffix}`.slice(0, 63)
        const created = await createContainer({
          name,
          image: imageName || imageRef,
          host_config: {},
          networking_config: {}
        })
        const containerId = created?.id || created?.container_id
        if (!containerId) {
          throw new Error('创建容器后未返回容器 ID')
        }
        await startContainer(containerId)
        // 成功文案需要拿到 containerId 才能拼出来，这里直接覆盖
        dialog.content = `✅ 成功启动容器: \n名称: ${name}\nID: ${containerId}`
      }
    )
  }

  const filteredImages = computed(() => {
    const q = state.query.toLowerCase()
    if (!q) return state.images
    return state.images.filter(img =>
      img.repo.toLowerCase().includes(q) || img.tag.toLowerCase().includes(q)
    )
  })

  const pagedImages = computed(() => {
    const start = (state.page - 1) * state.pageSize
    return filteredImages.value.slice(start, start + state.pageSize)
  })

  return {
    state,
    dialog,
    pagedImages,
    filteredImages,
    searchImages,
    pullImage,
    runImage,
    removeImage
  }
}

export type ImageState = ReturnType<typeof ImageState>
