package service

// Response 统一响应基础结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// PageResponse 分页响应结构
type PageResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Total   int64       `json:"total"`
	Page    int         `json:"page"`
	Size    int         `json:"size"`
}

// Success 成功响应
func Success(data interface{}) Response {
	return Response{Code: 0, Message: "操作成功", Data: data}
}

// Error 错误响应
func Error(code int, message string) Response {
	return Response{Code: code, Message: message, Data: nil}
}

// 系统级错误 (1000-1999)
const (
	ErrCodeSystem       = 1000 // 系统错误
	ErrCodeInvalidParam = 1001 // 参数错误
	ErrCodeNotFound     = 1004 // 资源不存在
	ErrCodeUnauthorized = 1401 // 未授权
	ErrCodeForbidden    = 1403 // 禁止访问

	// Docker操作错误 (2000-2999)
	ErrCodeDockerClient  = 2000 // Docker客户端错误
	ErrCodeDockerAPI     = 2001 // Docker API错误
	ErrCodeDockerTimeout = 2002 // Docker操作超时

	// 容器相关错误 (3000-3999)
	ErrCodeContainerNotFound = 3001 // 容器不存在
	ErrCodeContainerRunning  = 3002 // 容器运行中
	ErrCodeContainerStopped  = 3003 // 容器已停止
	ErrCodeContainerConflict = 3004 // 容器名称冲突

	// 镜像相关错误 (4000-4999)
	ErrCodeImageNotFound = 4001 // 镜像不存在
	ErrCodeImageInUse    = 4002 // 镜像使用中

	// 网络相关错误 (5000-5999)
	ErrCodeNetworkNotFound = 5001 // 网络不存在
	ErrCodeNetworkInUse    = 5002 // 网络使用中

	// Volume相关错误 (6000-6999)
	ErrCodeVolumeNotFound = 6001 // 数据卷不存在
	ErrCodeVolumeInUse    = 6002 // 数据卷使用中
)
