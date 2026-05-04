import { createContainer } from './container'
import { getImageList } from './image'
import { getNetworkList } from './network'
import { getVolumeList } from './volume'

export async function getCreateContainerImageOptions(): Promise<any[]> {
  const res = await getImageList({ all: true })
  return Array.isArray(res) ? res : []
}

export async function getCreateContainerNetworkOptions(): Promise<any[]> {
  const res = await getNetworkList()
  return Array.isArray(res) ? res : []
}

export async function getCreateContainerVolumeOptions(): Promise<any[]> {
  const res = await getVolumeList()
  return Array.isArray(res) ? res : []
}

export async function createNewContainer(payload: any): Promise<any> {
  return await createContainer(payload)
}
