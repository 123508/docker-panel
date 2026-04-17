package service

// Response is the standard API response shape.
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// PageResponse is the standard paged API response shape.
type PageResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Total   int64       `json:"total"`
	Page    int         `json:"page"`
	Size    int         `json:"size"`
}

func Success(data interface{}) Response {
	return Response{Code: 0, Message: "operation success", Data: data}
}

func Error(code int, message string) Response {
	return Response{Code: code, Message: message, Data: nil}
}

const (
	ErrCodeSystem       = 1000
	ErrCodeInvalidParam = 1001
	ErrCodeNotFound     = 1004
	ErrCodeUnauthorized = 1401
	ErrCodeForbidden    = 1403

	ErrCodeDockerClient  = 2000
	ErrCodeDockerAPI     = 2001
	ErrCodeDockerTimeout = 2002

	ErrCodeContainerNotFound = 3001
	ErrCodeContainerRunning  = 3002
	ErrCodeContainerStopped  = 3003
	ErrCodeContainerConflict = 3004

	ErrCodeImageNotFound = 4001
	ErrCodeImageInUse    = 4002

	ErrCodeNetworkNotFound = 5001
	ErrCodeNetworkInUse    = 5002

	ErrCodeVolumeNotFound = 6001
	ErrCodeVolumeInUse    = 6002
)
