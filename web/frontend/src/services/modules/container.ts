import {call} from "../common";

/* ==================== 数据结构定义 ==================== */

// 容器列表项
export interface ContainerListItem {
    id: string;                   // 容器ID
    names: string[];              // 容器名称
    image: string;                // 镜像名称
    image_id: string;             // 镜像ID
    command: string;              // 启动命令
    created: number;              // 创建时间戳
    state: string;                // 状态: created/running/paused/exited/dead
    status: string;               // 状态描述: Up 2 hours
    ports: PortBinding[];         // 端口映射
    labels: Record<string, string>; // 标签
    mounts: MountPoint[];         // 挂载点
}

// 端口绑定
export interface PortBinding {
    ip: string;             // 绑定IP
    private_port: number;   // 容器内端口
    public_port: number;    // 主机端口
    type: string;           // 协议类型: tcp/udp
}

// 挂载点
export interface MountPoint {
    type: string;           // 类型: bind/volume/tmpfs
    name: string;           // 卷名称(仅volume类型)
    source: string;         // 源路径
    destination: string;    // 目标路径
    driver: string;         // 驱动(仅volume类型)
    mode: string;           // 模式: ro/rw
    rw: boolean;            // 是否可写
}

// 容器详情
export interface ContainerDetail {
    id: string;
    created: string;
    path: string;
    args: string[];
    state: ContainerState;
    image: string;
    name: string;
    restart_count: number;
    driver: string;
    platform: string;
    mounts: MountPoint[];
    config: ContainerConfig;
    network_settings: NetworkSettings;
    host_config: HostConfig;
}

// 容器状态
export interface ContainerState {
    status: string;       // 状态
    running: boolean;
    paused: boolean;
    restarting: boolean;
    oom_killed: boolean;
    dead: boolean;
    pid: number;
    exit_code: number;
    error: string;
    started_at: string;
    finished_at: string;
}

// 容器配置
export interface ContainerConfig {
    hostname: string;
    domainname: string;
    user: string;
    attach_stdin: boolean;
    attach_stdout: boolean;
    attach_stderr: boolean;
    tty: boolean;
    open_stdin: boolean;
    stdin_once: boolean;
    env: string[];
    cmd: string[];
    image: string;
    volumes: Record<string, unknown>;
    working_dir: string;
    entrypoint: string[];
    labels: Record<string, string>;
    exposed_ports: Record<string, unknown>;
}

// 网络设置
export interface NetworkSettings {
    bridge: string;
    sandbox_id: string;
    hairpin_mode: boolean;
    link_local_ipv6_address: string;
    link_local_ipv6_prefix_len: number;
    ports: Record<string, PortBinding[]>;
    sandbox_key: string;
    ip_address: string;
    ip_prefix_len: number;
    ipv6_gateway: string;
    gateway: string;
    mac_address: string;
    networks: Record<string, EndpointSettings>;
}

// 端点设置
export interface EndpointSettings {
    ipam_config?: EndpointIPAMConfig | null;
    links: string[];
    aliases: string[];
    network_id: string;
    endpoint_id: string;
    gateway: string;
    ip_address: string;
    ip_prefix_len: number;
    ipv6_gateway: string;
    global_ipv6_address: string;
    global_ipv6_prefix_len: number;
    mac_address: string;
}

// 主机配置
export interface HostConfig {
    binds: string[];
    network_mode: string;
    port_bindings: Record<string, PortBinding[]>;
    restart_policy: RestartPolicy;
    auto_remove: boolean;
    volume_driver: string;
    volumes_from: string[];
    cap_add: string[];
    cap_drop: string[];
    dns: string[];
    dns_options: string[];
    dns_search: string[];
    extra_hosts: string[];
    links: string[];
    privileged: boolean;
    publish_all_ports: boolean;
    readonly_rootfs: boolean;
    memory: number;
    memory_swap: number;
    nano_cpus: number;
    cpu_shares: number;
}

// 重启策略
export interface RestartPolicy {
    name: 'no' | 'always' | 'on-failure' | 'unless-stopped';
    maximum_retry_count: number;
}

// 容器创建请求
export interface ContainerCreateRequest {
    name: string;
    image: string;
    cmd?: string[];
    entrypoint?: string[];
    env?: string[];
    working_dir?: string;
    exposed_ports?: Record<string, unknown>;
    host_config: HostConfigRequest;
    networking_config: NetworkingConfigRequest;
    labels?: Record<string, string>;
}

// 主机配置请求
export interface HostConfigRequest {
    port_bindings?: Record<string, PortBindingRequest[]>;
    binds?: string[];
    mounts?: MountRequest[];
    network_mode?: string;
    restart_policy?: RestartPolicy;
    memory?: number;
    nano_cpus?: number;
    privileged?: boolean;
}

// 端口绑定请求
export interface PortBindingRequest {
    host_ip?: string;
    host_port?: string;
}

// 挂载请求
export interface MountRequest {
    type: string;      // bind/volume/tmpfs
    source: string;
    target: string;
    read_only?: boolean;
    bind_options?: BindOptions;
    volume_options?: VolumeOptions;
}

// 绑定挂载选项
export interface BindOptions {
    propagation?: string; // rprivate/private/rshared/shared/rslave/slave
}

// 卷挂载选项
export interface VolumeOptions {
    no_copy?: boolean;
    labels?: Record<string, string>;
    driver_config?: DriverConfig;
}

// 驱动配置
export interface DriverConfig {
    name: string;
    options?: Record<string, string>;
}

// 网络配置请求
export interface NetworkingConfigRequest {
    endpoints_config?: Record<string, EndpointSettings>;
}

// 容器日志选项
export interface ContainerLogsOptions {
    show_stdout?: boolean;
    show_stderr?: boolean;
    since?: string;
    until?: string;
    timestamps?: boolean;
    follow?: boolean;
    tail?: string;
}

// 容器执行命令请求
export interface ContainerExecRequest {
    attach_stdin?: boolean;
    attach_stdout?: boolean;
    attach_stderr?: boolean;
    tty?: boolean;
    cmd: string[];
    env?: string[];
    working_dir?: string;
    user?: string;
}

// 容器进程信息
export interface ContainerTopResponse {
    titles: string[];
    processes: string[][];
}

// 端点IPAM配置（可选）
export interface EndpointIPAMConfig {
    ipv4_address?: string;
    ipv6_address?: string;
    link_local_ips?: string[];
}

/* ==================== 请求类型定义 ==================== */

export interface ContainerListRequest {
    all?: boolean
    limit?: number
    filters?: string
}

export interface ContainerStopRequest {
    timeout?: number
}

export interface ContainerKillRequest {
    signal?: string
}

export interface ContainerRestartRequest {
    timeout?: number
}

export interface ContainerRemoveRequest {
    force?: boolean
    remove_volumes?: boolean
}

export interface ContainerRenameRequest {
    new_name: string
}

export interface ContainerTopRequest {
    ps_args?: string
}

export interface ContainerExecRequest {
    cmd: string[]
    env?: string[]
    workdir?: string
}

/* ==================== API ==================== */

// 3.3.1 容器列表
export async function getContainerList(params?: ContainerListRequest): Promise<any> {
    const res = await call("GET", "/api/v1/containers", params).unwrap()
    return res
}

// 3.3.2 容器详情
export async function getContainerInspect(id: string): Promise<any> {
    const res = await call("GET", `/api/v1/containers/${id}`, null).unwrap()
    return res
}

// 3.3.3 创建容器
export async function createContainer(data: any): Promise<any> {
    const res = await call("POST", "/api/v1/containers", data).unwrap()
    return res
}

// 3.3.4 启动容器
export async function startContainer(id: string): Promise<any> {
    const res = await call("POST", `/api/v1/containers/${id}/start`, null).unwrap()
    return res
}

// 3.3.5 停止容器
export async function stopContainer(id: string, data?: ContainerStopRequest): Promise<any> {
    const res = await call("POST", `/api/v1/containers/${id}/stop`, data).unwrap()
    return res
}

// 3.3.6 强制停止容器
export async function killContainer(id: string, data?: ContainerKillRequest): Promise<any> {
    const res = await call("POST", `/api/v1/containers/${id}/kill`, data).unwrap()
    return res
}

// 3.3.7 重启容器
export async function restartContainer(id: string, data?: ContainerRestartRequest): Promise<any> {
    const res = await call("POST", `/api/v1/containers/${id}/restart`, data).unwrap()
    return res
}

// 3.3.8 暂停容器
export async function pauseContainer(id: string): Promise<any> {
    const res = await call("POST", `/api/v1/containers/${id}/pause`, null).unwrap()
    return res
}

// 3.3.9 恢复容器
export async function unpauseContainer(id: string): Promise<any> {
    const res = await call("POST", `/api/v1/containers/${id}/unpause`, null).unwrap()
    return res
}

// 3.3.10 删除容器
export async function removeContainer(id: string, params?: ContainerRemoveRequest): Promise<any> {
    const res = await call("DELETE", `/api/v1/containers/${id}`, params).unwrap()
    return res
}

// 3.3.11 重命名容器
export async function renameContainer(id: string, data: ContainerRenameRequest): Promise<any> {
    const res = await call("POST", `/api/v1/containers/${id}/rename`, data).unwrap()
    return res
}

// 3.3.12 容器日志
export async function getContainerLogs(id: string, params?: any): Promise<any> {
    const res = await call("GET", `/api/v1/containers/${id}/logs`, params).unwrap()
    return res
}

// 3.3.13 exec
export async function execContainer(id: string, data: ContainerExecRequest): Promise<any> {
    const res = await call("POST", `/api/v1/containers/${id}/exec`, data).unwrap()
    return res
}

// 3.3.14 终端（WebSocket 地址生成）
export function getContainerTerminalWS(id: string): string {
    return `/api/v1/containers/${id}/terminal`
}

// 3.3.15 进程列表
export async function getContainerTop(id: string, params?: ContainerTopRequest): Promise<any> {
    const res = await call("GET", `/api/v1/containers/${id}/top`, params).unwrap()
    return res
}

// 3.3.16 端口映射
export async function getContainerPorts(id: string): Promise<any> {
    const res = await call("GET", `/api/v1/containers/${id}/ports`, null).unwrap()
    return res
}

// 3.3.17 导出容器
export function exportContainer(id: string): string {
    return `/api/v1/containers/${id}/export`
}

// 3.3.20 挂载信息
export async function getContainerMounts(id: string): Promise<any> {
    const res = await call("GET", `/api/v1/containers/${id}/mounts`, null).unwrap()
    return res
}