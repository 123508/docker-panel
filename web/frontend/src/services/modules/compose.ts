import { call } from "../common"

/* ==================== 数据结构定义 ==================== */

// Compose服务
export interface ComposeService {
    name: string
    image: string
    containers: string[]
    replicas: number
    status: string
}

// Compose项目
export interface ComposeProject {
    name: string
    file_path: string
    status: string
    services: ComposeService[]
    created_at: string
    updated_at: string
}

/* ==================== 请求类型定义 ==================== */

export interface ComposeUploadRequest {
    name: string
    content: string
}

export interface ComposeUpRequest {
    build?: boolean
    force_recreate?: boolean
    no_recreate?: boolean
    no_build?: boolean
    timeout?: number
    remove_orphans?: boolean
    services?: string[]
}

export interface ComposeStopRequest {
    timeout?: number
    services?: string[]
}

export interface ComposeRestartRequest {
    timeout?: number
    services?: string[]
}

export interface ComposeDownRequest {
    remove_volumes?: boolean
    remove_images?: string
    remove_orphans?: boolean
    timeout?: number
}

export interface ComposeScaleRequest {
    services: Record<string, number>
}

export interface ComposeLogsRequest {
    follow?: boolean
    tail?: string
    since?: string
    until?: string
    timestamps?: boolean
    services?: string[]
}

/* ==================== API ==================== */

// 7.2.1 上传compose文件
export async function uploadComposeProject(data: ComposeUploadRequest): Promise<any> {
    const res = await call("POST", "/api/v1/compose/projects", data).unwrap()
    return res
}

// 7.2.2 启动compose项目
export async function upComposeProject(
    name: string,
    data?: ComposeUpRequest
): Promise<any> {
    const res = await call("POST", `/api/v1/compose/projects/${name}/up`, data).unwrap()
    return res
}

// 7.2.3 停止compose项目
export async function stopComposeProject(
    name: string,
    data?: ComposeStopRequest
): Promise<any> {
    const res = await call("POST", `/api/v1/compose/projects/${name}/stop`, data).unwrap()
    return res
}

// 7.2.4 重启compose项目
export async function restartComposeProject(
    name: string,
    data?: ComposeRestartRequest
): Promise<any> {
    const res = await call("POST", `/api/v1/compose/projects/${name}/restart`, data).unwrap()
    return res
}

// 7.2.5 删除compose项目
export async function downComposeProject(
    name: string,
    data?: ComposeDownRequest
): Promise<any> {
    const res = await call("DELETE", `/api/v1/compose/projects/${name}`, data).unwrap()
    return res
}

// 7.2.6 查看compose状态
export async function getComposePS(name: string): Promise<any> {
    const res = await call("GET", `/api/v1/compose/projects/${name}/ps`, null).unwrap()
    return res
}

// 7.2.7 查看compose日志（流式）
export async function getComposeLogs(
    name: string,
    params?: ComposeLogsRequest
): Promise<any> {
    const res = await call("GET", `/api/v1/compose/projects/${name}/logs`, params).unwrap()
    return res
}

// 7.2.8 扩展服务数量
export async function scaleCompose(
    name: string,
    data: ComposeScaleRequest
): Promise<any> {
    const res = await call("POST", `/api/v1/compose/projects/${name}/scale`, data).unwrap()
    return res
}