import {
  getContainerInspect,
  startContainer,
  stopContainer,
  restartContainer,
  removeContainer
} from './container'

export async function getContainerDetail(id: string): Promise<any> {
  return await getContainerInspect(id)
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

export async function removeContainerById(id: string): Promise<any> {
  return await removeContainer(id, { force: true })
}
