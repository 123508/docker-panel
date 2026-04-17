import { call } from "../common"

/* ==================== 数据结构定义 ==================== */

// 镜像列表项
export interface ImageListItem {
    id: string
    parent_id: string
    repo_tags: string[]
    repo_digests: string[]
    created: number
    size: number
    virtual_size: number
    shared_size: number
    labels: Record<string, string>
    containers: number
}

// 图形驱动
export interface GraphDriver {
    name: string
    data: Record<string, string>
}

// 根文件系统
export interface RootFS {
    type: string
    layers: string[]
}

// 镜像元数据
export interface ImageMetadata {
    last_tag_time: string
}

// 容器配置（复用容器定义，可按需精简）
export interface ContainerConfig {
    hostname: string
    domainname: string
    user: string
    env: string[]
    cmd: string[]
    image: string
    working_dir: string
    entrypoint: string[]
    labels: Record<string, string>
}

// 镜像详情
export interface ImageDetail {
    id: string
    repo_tags: string[]
    repo_digests: string[]
    parent: string
    comment: string
    created: string
    container: string
    container_config: ContainerConfig
    docker_version: string
    author: string
    config: ContainerConfig
    architecture: string
    os: string
    size: number
    virtual_size: number
    graph_driver: GraphDriver
    root_fs: RootFS
    metadata: ImageMetadata
}

/* ==================== 请求类型定义 ==================== */

export interface ImageListRequest {
    all?: boolean
    filters?: string
}

export interface ImagePullRequest {
    image: string
}

export interface ImageRemoveRequest {
    force?: boolean
    no_prune?: boolean
}

export interface ImageBuildRequest {
    dockerfile: string
    tags?: string[]
    build_args?: Record<string, string>
    no_cache?: boolean
    remove?: boolean
    force_rm?: boolean
    target?: string
}

export interface ImageTagRequest {
    repository: string
    tag: string
}

export interface ImagePushRequest {
    tag?: string
}

/* ==================== API ==================== */

// 4.2.1 镜像列表
export async function getImageList(params?: ImageListRequest): Promise<any> {
    const res = await call("GET", "/api/v1/images", params).unwrap()
    return res
}

// 4.2.2 镜像详情
export async function getImageInspect(id: string): Promise<any> {
    const res = await call("GET", `/api/v1/images/${id}`, null).unwrap()
    return res
}

// 4.2.3 拉取镜像（流式）
export async function pullImage(data: ImagePullRequest): Promise<any> {
    const res = await call("POST", "/api/v1/images/pull", data).unwrap()
    return res
}

// 4.2.4 删除镜像
export async function removeImage(id: string, params?: ImageRemoveRequest): Promise<any> {
    const res = await call("DELETE", `/api/v1/images/${id}`, params).unwrap()
    return res
}

// 4.2.5 构建镜像（流式）
export async function buildImage(data: ImageBuildRequest): Promise<any> {
    const res = await call("POST", "/api/v1/images/build", data).unwrap()
    return res
}

// 4.2.6 导出镜像（下载）
export function saveImage(id: string): string {
    return `/api/v1/images/${id}/save`
}

// 4.2.7 导入镜像
export async function loadImage(data: FormData): Promise<any> {
    const res = await call("POST", "/api/v1/images/load", data).unwrap()
    return res
}

// 4.2.8 打标签
export async function tagImage(id: string, data: ImageTagRequest): Promise<any> {
    const res = await call("POST", `/api/v1/images/${id}/tag`, data).unwrap()
    return res
}

// 4.2.9 推送镜像（流式）
export async function pushImage(id: string, data?: ImagePushRequest): Promise<any> {
    const res = await call("POST", `/api/v1/images/${id}/push`, data).unwrap()
    return res
}