import {http} from "./api";

// 响应基础结构
export interface Response<T = any>{
    code: number         // 业务状态码: 0-成功, 非0-失败
    message: string      // 提示信息
    data?: T           // 业务数据
}

// 分页响应结构
export interface PageResponse<T = any> extends Response<T> {
    total: number       // 总数
    page: number        // 当前页
    size: number        // 每页大小
}

export class CustomError extends Error{
    code: number

    constructor(code: number,message: string){
        super(message);
        this.code=code
        this.name="CustomError"

        Object.setPrototypeOf(this, CustomError.prototype);
    }

    toString():string{
        return `code:${this.code} , message:${this.message}`
    }
}

export function call<T>(
    method: string,
    endpoint: string,
    data: any
) {
    const promise = (async () => {
        try {
            const resp = await http({ method, url: endpoint, data })
            const res = resp.data

            if (res.code !== 0) {
                return [new CustomError(res.code, res.message), null]
            }

            return [null, res]
        } catch (e: any) {
            return [new CustomError(e.code ?? 1000, e.message ?? '系统错误'), null]
        }
    })()

    return {
        // 原始模式
        raw: () => promise,

        // 直接解包
        async unwrap(): Promise<T> {
            const [err, res] = await promise

            if (err) {
                throw err
            }

            return res?.data ?? null
        }
    }
}

export enum ErrCode {
    System = 1000,              // 系统错误
    InvalidParam = 1001,        // 参数错误
    NotFound = 1004,            // 资源不存在
    Unauthorized = 1401,        // 未授权
    Forbidden = 1403,           // 禁止访问

    DockerClient = 2000,        // Docker客户端错误
    DockerAPI = 2001,           // Docker API错误
    DockerTimeout = 2002,       // Docker操作超时

    ContainerNotFound = 3001,   // 容器不存在
    ContainerRunning = 3002,    // 容器运行中
    ContainerStopped = 3003,    // 容器已停止
    ContainerConflict = 3004,   // 容器名称冲突

    ImageNotFound = 4001,       // 镜像不存在
    ImageInUse = 4002,          // 镜像使用中

    NetworkNotFound = 5001,     // 网络不存在
    NetworkInUse = 5002,        // 网络使用中

    VolumeNotFound = 6001,      // 数据卷不存在
    VolumeInUse = 6002          // 数据卷使用中
}

const ErrCodeMessage: Record<number, string> = {
    [ErrCode.System]: "系统错误",
    [ErrCode.InvalidParam]: "参数错误",
    [ErrCode.NotFound]: "资源不存在",
    [ErrCode.Unauthorized]: "未授权",
    [ErrCode.Forbidden]: "禁止访问",

    [ErrCode.DockerClient]: "Docker客户端错误",
    [ErrCode.DockerAPI]: "Docker API错误",
    [ErrCode.DockerTimeout]: "Docker操作超时",

    [ErrCode.ContainerNotFound]: "容器不存在",
    [ErrCode.ContainerRunning]: "容器运行中",
    [ErrCode.ContainerStopped]: "容器已停止",
    [ErrCode.ContainerConflict]: "容器名称冲突",

    [ErrCode.ImageNotFound]: "镜像不存在",
    [ErrCode.ImageInUse]: "镜像使用中",

    [ErrCode.NetworkNotFound]: "网络不存在",
    [ErrCode.NetworkInUse]: "网络使用中",

    [ErrCode.VolumeNotFound]: "数据卷不存在",
    [ErrCode.VolumeInUse]: "数据卷使用中",
}