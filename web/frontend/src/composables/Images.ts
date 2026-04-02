import { reactive, computed } from 'vue'

export function ImageState() {
  const state = reactive({
    query: '',
    page: 1,
    pageSize: 5,

    stats: [
      { label: '总数:', value: '42' },
      { label: '大小:', value: '12.4 GB' },
      { label: '悬空:', value: '3', variant: 'dangling' }
    ],

    images: [
      { repo: 'nginx', tag: 'latest', id: 'a72860cb95fd', size: '187 MB', created: '2 天前', color: '#3B82F6' },
      { repo: 'postgres', tag: '16-alpine', id: 'b4d181a07f80', size: '432 MB', created: '5 天前', color: '#22C55E' },
      { repo: 'redis', tag: '7-alpine', id: '3c41ce05add9', size: '41 MB', created: '1 周前', color: '#EF4444' },
      { repo: 'node', tag: '20-alpine', id: 'c7b5a7e3f2d1', size: '124 MB', created: '2 周前', color: '#FACC15' }
    ]
  })

  const pagedImages = computed(() => {
    const start = (state.page - 1) * state.pageSize
    return state.images.slice(start, start + state.pageSize)
  })

  const filteredImages = computed(() => {
    const q = state.query.toLowerCase()
    if (!q) return state.images
    return state.images.filter(img =>
      img.repo.toLowerCase().includes(q) || img.tag.toLowerCase().includes(q)
    )
  })

  return {
    state,
    pagedImages,
    filteredImages
  }
}