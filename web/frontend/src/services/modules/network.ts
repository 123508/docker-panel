import { call } from "../common"

/* ==================== 数据结构定义 ==================== */

// IPAM配置
export interface IPAMConfig {
    subnet: string
    ip_range: string
    gateway: string
    aux_address: Record<string, string>
}

// IP地址管理
export interface IPAM {
    driver: string
    options: Record<string, string>
    config: IPAMConfig[]
}

// 网络容器
export interface NetworkContainer {
    name: string
    endpoint_id: string
    mac_address: string
    ipv4_address: string
    ipv6_address: string
}

// 网络列表项
export interface NetworkListItem {
    id: string
    name: string
    created: string
    scope: string
    driver: string
    enable_ipv6: boolean
    internal: boolean
    attachable: boolean
    ingress: boolean
    ipam: IPAM
    containers: Record<string, NetworkContainer>
    options: Record<string, string>
    labels: Record<string, string>
}

// 网络详情
export interface NetworkDetail {
    name: string
    id: string
    created: string
    scope: string
    driver: string
    enable_ipv6: boolean
    ipam: IPAM
    internal: boolean
    attachable: boolean
    ingress: boolean
    containers: Record<string, NetworkContainer>
    options: Record<string, string>
    labels: Record<string, string>
}

/* ==================== 请求类型定义 ==================== */

export interface NetworkListRequest {
    filters?: string
}

export interface NetworkInspectRequest {
    verbose?: boolean
    scope?: string
}

export interface NetworkCreateRequest {
    name: string
    driver?: string
    scope?: string
    enable_ipv6?: boolean
    ipam?: IPAM
    internal?: boolean
    attachable?: boolean
    ingress?: boolean
    options?: Record<string, string>
    labels?: Record<string, string>
}

export interface NetworkConnectRequest {
    container: string
    endpoint_config?: any
}

export interface NetworkDisconnectRequest {
    container: string
    force?: boolean
}

/* ==================== API ==================== */

// 6.2.1 网络列表
export async function getNetworkList(params?: NetworkListRequest): Promise<any> {
    const res = await call("GET", "/api/v1/networks", params).unwrap()
    return res
}

// 6.2.2 网络详情
export async function getNetworkInspect(
    id: string,
    params?: NetworkInspectRequest
): Promise<any> {
    const res = await call("GET", `/api/v1/networks/${id}`, params).unwrap()
    return res
}

// 6.2.3 创建网络
export async function createNetwork(data: NetworkCreateRequest): Promise<any> {
    const res = await call("POST", "/api/v1/networks", data).unwrap()
    return res
}

// 6.2.4 删除网络
export async function removeNetwork(id: string): Promise<any> {
    const res = await call("DELETE", `/api/v1/networks/${id}`, null).unwrap()
    return res
}

// 6.2.5 连接容器到网络
export async function connectNetwork(
    id: string,
    data: NetworkConnectRequest
): Promise<any> {
    const res = await call("POST", `/api/v1/networks/${id}/connect`, data).unwrap()
    return res
}

// 6.2.6 断开容器网络
export async function disconnectNetwork(
    id: string,
    data: NetworkDisconnectRequest
): Promise<any> {
    const res = await call("POST", `/api/v1/networks/${id}/disconnect`, data).unwrap()
    return res
}