import { call } from "../common"

/* ==================== 数据结构定义 ==================== */

// 数据卷列表项
export interface VolumeListItem {
    created_at: string
    driver: string
    labels: Record<string, string>
    mountpoint: string
    name: string
    options: Record<string, string>
    scope: string
}

// 数据卷使用数据
export interface VolumeUsageData {
    size: number
    ref_count: number
}

// 数据卷详情
export interface VolumeDetail {
    name: string
    driver: string
    mountpoint: string
    created_at: string
    labels: Record<string, string>
    options: Record<string, string>
    scope: string
    status: Record<string, any>
    usage_data: VolumeUsageData
}

// Volume 被哪些容器使用
export interface VolumeContainersResponse {
    volume_name: string
    containers: string[]
}

/* ==================== 请求类型定义 ==================== */

export interface VolumeListRequest {
    filters?: string
}

export interface VolumeCreateRequest {
    name: string
    driver?: string
    driver_opts?: Record<string, string>
    labels?: Record<string, string>
}

export interface VolumeRemoveRequest {
    force?: boolean
}

/* ==================== API ==================== */

// 5.2.1 Volume列表
export async function getVolumeList(params?: VolumeListRequest): Promise<any> {
    const res = await call("GET", "/api/v1/volumes", params).unwrap()
    return res
}

// 5.2.2 Volume详情
export async function getVolumeInspect(name: string): Promise<any> {
    const res = await call("GET", `/api/v1/volumes/${name}`, null).unwrap()
    return res
}

// 5.2.3 创建Volume
export async function createVolume(data: VolumeCreateRequest): Promise<any> {
    const res = await call("POST", "/api/v1/volumes", data).unwrap()
    return res
}

// 5.2.4 删除Volume
export async function removeVolume(name: string, params?: VolumeRemoveRequest): Promise<any> {
    const res = await call("DELETE", `/api/v1/volumes/${name}`, params).unwrap()
    return res
}

// 5.2.5 查看Volume被哪些容器使用
export async function getVolumeContainers(name: string): Promise<any> {
    const res = await call("GET", `/api/v1/volumes/${name}/containers`, null).unwrap()
    return res
}