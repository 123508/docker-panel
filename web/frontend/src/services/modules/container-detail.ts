import {
  getContainerInspect,
  startContainer,
  stopContainer,
  restartContainer
} from './container'
import { http } from '@/services/api'
import { token } from '@/store/auth'

export async function getContainerDetail(id: string): Promise<any> {
  return await getContainerInspect(id)
}

export async function getContainerDetailPage(id: string): Promise<{ detail: any; logs: string }> {
  const [detail, logs] = await Promise.all([
    getContainerDetail(id),
    getContainerDetailLogs(id, '100').catch(() => '')
  ])
  return { detail, logs }
}

export async function getContainerDetailLogs(id: string, tail = '100'): Promise<string> {
  const response = await http.get(`/api/v1/containers/${id}/logs`, {
    params: {
      show_stdout: true,
      show_stderr: true,
      timestamps: true,
      tail
    },
    responseType: 'text',
    transformResponse: [(data) => data]
  })
  return typeof response.data === 'string' ? response.data : String(response.data || '')
}

export function createContainerLogsWebSocket(
  id: string,
  params: { tail?: string; timestamps?: boolean } = {}
): WebSocket {
  const wsProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  const searchParams = new URLSearchParams({
    show_stdout: 'true',
    show_stderr: 'true',
    follow: 'true',
    tail: params.tail ?? '100',
    timestamps: String(params.timestamps ?? true)
  })
  if (token.value) {
    searchParams.set('token', token.value)
  }
  const wsUrl = `${wsProtocol}//${window.location.host}/api/v1/containers/${encodeURIComponent(id)}/logs/ws?${searchParams.toString()}`
  return new WebSocket(wsUrl)
}

export async function startContainerById(id: string): Promise<any> {
  return await startContainer(id)
}

export async function stopContainerById(id: string): Promise<any> {
  return await stopContainer(id)
}

export async function restartContainerById(id: string): Promise<any> {
  return await restartContainer(id)
}
