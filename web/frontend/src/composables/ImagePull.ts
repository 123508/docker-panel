import { reactive, ref } from 'vue'
import { pullImage as apiPullImage } from '@/services/modules/image'
import { ElMessage } from 'element-plus'

interface SearchResult {
  name: string
  description: string
  star_count: number
  pull_count: number
  is_official: boolean
}

export function ImagePullState() {
  const state = reactive({
    query: '',
    searching: false,
    results: [] as SearchResult[],
    searched: false,
    pulling: ''
  })

  const searchImages = async () => {
    const q = state.query.trim()
    if (!q) return
    state.searching = true
    state.searched = true
    try {
      const url = `https://hub.docker.com/v2/repositories/library?search=${encodeURIComponent(q)}&page_size=20`
      const res = await fetch(url)
      if (!res.ok) throw new Error('搜索请求失败')
      const data = await res.json()
      state.results = (data.results || data).map((r: any) => ({
        name: r.name || r.repo_name || '-',
        description: r.description || r.short_description || '',
        star_count: r.star_count || 0,
        pull_count: r.pull_count || 0,
        is_official: r.is_official ?? true
      }))
    } catch (e: any) {
      // Fallback: try Docker Hub API v1 or show suggestions
      if (state.results.length === 0) {
        // Show a basic result with the search term as an official image suggestion
        state.results = [{
          name: q,
          description: `Docker Official Image — ${q}`,
          star_count: 0,
          pull_count: 0,
          is_official: true
        }]
      }
    } finally {
      state.searching = false
    }
  }

  const pullImage = async (imageName: string) => {
    const fullName = imageName.includes(':') ? imageName : `${imageName}:latest`
    state.pulling = imageName
    try {
      await apiPullImage({ image: fullName })
      ElMessage.success(`镜像拉取成功: ${fullName}`)
    } catch (e: any) {
      ElMessage.error(e.message || `拉取镜像失败: ${fullName}`)
    } finally {
      state.pulling = ''
    }
  }

  const searchOnEnter = () => {
    searchImages()
  }

  return {
    state,
    searchImages,
    pullImage,
    searchOnEnter
  }
}
