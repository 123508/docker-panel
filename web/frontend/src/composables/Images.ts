import { reactive, computed, onMounted } from 'vue'
import { getImageList, removeImage as apiRemoveImage } from '@/services/modules/image'
import { ElMessage } from 'element-plus'

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
    try {
      await apiRemoveImage(id, { force: true })
      ElMessage.success('镜像已删除')
      await loadData()
    } catch(e: any) {
      ElMessage.error(e.message || '删除镜像失败')
    }
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
    pagedImages,
    filteredImages,
    removeImage
  }
}