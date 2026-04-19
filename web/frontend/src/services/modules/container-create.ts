import { createContainer } from './container'
import { getImageList } from './image'
import { getNetworkList } from './network'

export async function getCreateContainerImageOptions(): Promise<any[]> {
  const res = await getImageList({ all: true })
  return Array.isArray(res) ? res : []
}

export async function getCreateContainerNetworkOptions(): Promise<any[]> {
  const res = await getNetworkList()
  return Array.isArray(res) ? res : []
}

export async function createNewContainer(payload: any): Promise<any> {
  return await createContainer(payload)
}
