package handler

import (
	"docker-panel/service"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var wsUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// ContainerHandler 容器 HTTP 处理器
type ContainerHandler struct {
	svc *service.ContainerService
}

// NewContainerHandler 创建 ContainerHandler
func NewContainerHandler(svc *service.ContainerService) *ContainerHandler {
	return &ContainerHandler{svc: svc}
}

// RegisterRoutes 注册容器路由
func (h *ContainerHandler) RegisterRoutes(rg *gin.RouterGroup) {
	containers := rg.Group("/containers")
	{
		containers.GET("", h.List)
		containers.POST("", h.Create)
		containers.GET("/:id", h.Inspect)
		containers.DELETE("/:id", h.Remove)
		containers.POST("/:id/start", h.Start)
		containers.POST("/:id/stop", h.Stop)
		containers.POST("/:id/kill", h.Kill)
		containers.POST("/:id/restart", h.Restart)
		containers.POST("/:id/pause", h.Pause)
		containers.POST("/:id/unpause", h.Unpause)
		containers.POST("/:id/rename", h.Rename)
		containers.GET("/:id/logs", h.Logs)
		containers.POST("/:id/exec", h.Exec)
		containers.GET("/:id/terminal", h.Terminal)
		containers.GET("/:id/top", h.Top)
		containers.GET("/:id/ports", h.Ports)
		containers.GET("/:id/mounts", h.Mounts)
		containers.GET("/:id/export", h.Export)
	}
}

// List GET /api/v1/containers
func (h *ContainerHandler) List(c *gin.Context) {
	var req service.ContainerListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		respondJSON(c, service.Error(service.ErrCodeInvalidParam, err.Error()))
		return
	}
	items, err := h.svc.ContainerList(c.Request.Context(), req)
	if err != nil {
		respondJSON(c, service.Error(service.ErrCodeDockerAPI, err.Error()))
		return
	}
	respondJSON(c, service.Success(items))
}

// Inspect GET /api/v1/containers/:id
func (h *ContainerHandler) Inspect(c *gin.Context) {
	detail, err := h.svc.ContainerInspect(c.Request.Context(), c.Param("id"))
	if err != nil {
		respondJSON(c, service.Error(service.ErrCodeContainerNotFound, err.Error()))
		return
	}
	respondJSON(c, service.Success(detail))
}

// Create POST /api/v1/containers
func (h *ContainerHandler) Create(c *gin.Context) {
	var req service.ContainerCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondJSON(c, service.Error(service.ErrCodeInvalidParam, err.Error()))
		return
	}
	id, err := h.svc.ContainerCreate(c.Request.Context(), req)
	if err != nil {
		respondJSON(c, service.Error(service.ErrCodeDockerAPI, err.Error()))
		return
	}
	respondJSON(c, service.Success(gin.H{"id": id}))
}

// Start POST /api/v1/containers/:id/start
func (h *ContainerHandler) Start(c *gin.Context) {
	if err := h.svc.ContainerStart(c.Request.Context(), c.Param("id")); err != nil {
		respondJSON(c, service.Error(service.ErrCodeDockerAPI, err.Error()))
		return
	}
	respondJSON(c, service.Success(nil))
}

// Stop POST /api/v1/containers/:id/stop
func (h *ContainerHandler) Stop(c *gin.Context) {
	var req service.ContainerStopRequest
	_ = c.ShouldBindJSON(&req)
	if err := h.svc.ContainerStop(c.Request.Context(), c.Param("id"), req.Timeout); err != nil {
		respondJSON(c, service.Error(service.ErrCodeDockerAPI, err.Error()))
		return
	}
	respondJSON(c, service.Success(nil))
}

// Kill POST /api/v1/containers/:id/kill
func (h *ContainerHandler) Kill(c *gin.Context) {
	var req service.ContainerKillRequest
	_ = c.ShouldBindJSON(&req)
	if err := h.svc.ContainerKill(c.Request.Context(), c.Param("id"), req.Signal); err != nil {
		respondJSON(c, service.Error(service.ErrCodeDockerAPI, err.Error()))
		return
	}
	respondJSON(c, service.Success(nil))
}

// Restart POST /api/v1/containers/:id/restart
func (h *ContainerHandler) Restart(c *gin.Context) {
	var req service.ContainerRestartRequest
	_ = c.ShouldBindJSON(&req)
	if err := h.svc.ContainerRestart(c.Request.Context(), c.Param("id"), req.Timeout); err != nil {
		respondJSON(c, service.Error(service.ErrCodeDockerAPI, err.Error()))
		return
	}
	respondJSON(c, service.Success(nil))
}

// Pause POST /api/v1/containers/:id/pause
func (h *ContainerHandler) Pause(c *gin.Context) {
	if err := h.svc.ContainerPause(c.Request.Context(), c.Param("id")); err != nil {
		respondJSON(c, service.Error(service.ErrCodeDockerAPI, err.Error()))
		return
	}
	respondJSON(c, service.Success(nil))
}

// Unpause POST /api/v1/containers/:id/unpause
func (h *ContainerHandler) Unpause(c *gin.Context) {
	if err := h.svc.ContainerUnpause(c.Request.Context(), c.Param("id")); err != nil {
		respondJSON(c, service.Error(service.ErrCodeDockerAPI, err.Error()))
		return
	}
	respondJSON(c, service.Success(nil))
}

// Remove DELETE /api/v1/containers/:id
func (h *ContainerHandler) Remove(c *gin.Context) {
	var req service.ContainerRemoveRequest
	_ = c.ShouldBindQuery(&req)
	if err := h.svc.ContainerRemove(c.Request.Context(), c.Param("id"), req.Force, req.RemoveVolumes); err != nil {
		respondJSON(c, service.Error(service.ErrCodeDockerAPI, err.Error()))
		return
	}
	respondJSON(c, service.Success(nil))
}

// Rename POST /api/v1/containers/:id/rename
func (h *ContainerHandler) Rename(c *gin.Context) {
	var req service.ContainerRenameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondJSON(c, service.Error(service.ErrCodeInvalidParam, err.Error()))
		return
	}
	if err := h.svc.ContainerRename(c.Request.Context(), c.Param("id"), req.NewName); err != nil {
		respondJSON(c, service.Error(service.ErrCodeDockerAPI, err.Error()))
		return
	}
	respondJSON(c, service.Success(nil))
}

// Logs GET /api/v1/containers/:id/logs
func (h *ContainerHandler) Logs(c *gin.Context) {
	var opts service.ContainerLogsOptions
	if err := c.ShouldBindQuery(&opts); err != nil {
		respondJSON(c, service.Error(service.ErrCodeInvalidParam, err.Error()))
		return
	}
	// 默认展示 stdout+stderr
	if !opts.ShowStdout && !opts.ShowStderr {
		opts.ShowStdout = true
		opts.ShowStderr = true
	}

	rc, err := h.svc.ContainerLogs(c.Request.Context(), c.Param("id"), opts)
	if err != nil {
		respondJSON(c, service.Error(service.ErrCodeDockerAPI, err.Error()))
		return
	}
	defer rc.Close()

	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.Header("Transfer-Encoding", "chunked")
	c.Status(http.StatusOK)
	io.Copy(c.Writer, rc)
}

// Exec POST /api/v1/containers/:id/exec
func (h *ContainerHandler) Exec(c *gin.Context) {
	var req service.ContainerExecRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondJSON(c, service.Error(service.ErrCodeInvalidParam, err.Error()))
		return
	}
	req.AttachStdout = true
	req.AttachStderr = true

	hijack, err := h.svc.ContainerExec(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		respondJSON(c, service.Error(service.ErrCodeDockerAPI, err.Error()))
		return
	}
	defer hijack.Close()

	output, err := io.ReadAll(hijack.Reader)
	if err != nil {
		respondJSON(c, service.Error(service.ErrCodeDockerAPI, err.Error()))
		return
	}
	respondJSON(c, service.Success(gin.H{"output": string(output)}))
}

// Terminal WebSocket /api/v1/containers/:id/terminal
func (h *ContainerHandler) Terminal(c *gin.Context) {
	containerID := c.Param("id")

	ws, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()

	ctx := c.Request.Context()

	// 创建 exec 会话（sh 优先，降级 bash）
	execID, err := h.svc.ContainerExecCreate(ctx, containerID, service.ContainerExecRequest{
		Cmd: []string{"sh"},
	})
	if err != nil {
		ws.WriteMessage(websocket.TextMessage, []byte("error: "+err.Error()))
		return
	}

	hijack, err := h.svc.ContainerExecAttach(ctx, execID)
	if err != nil {
		ws.WriteMessage(websocket.TextMessage, []byte("error: "+err.Error()))
		return
	}
	defer hijack.Close()

	// WebSocket -> Docker stdin
	go func() {
		for {
			_, msg, err := ws.ReadMessage()
			if err != nil {
				return
			}
			hijack.Conn.Write(msg)
		}
	}()

	// Docker stdout -> WebSocket
	buf := make([]byte, 4096)
	for {
		n, err := hijack.Reader.Read(buf)
		if n > 0 {
			ws.WriteMessage(websocket.BinaryMessage, buf[:n])
		}
		if err != nil {
			break
		}
	}
}

// Top GET /api/v1/containers/:id/top
func (h *ContainerHandler) Top(c *gin.Context) {
	var req service.ContainerTopRequest
	_ = c.ShouldBindQuery(&req)
	result, err := h.svc.ContainerTop(c.Request.Context(), c.Param("id"), req.PsArgs)
	if err != nil {
		respondJSON(c, service.Error(service.ErrCodeDockerAPI, err.Error()))
		return
	}
	respondJSON(c, service.Success(result))
}

// Ports GET /api/v1/containers/:id/ports
func (h *ContainerHandler) Ports(c *gin.Context) {
	detail, err := h.svc.ContainerInspect(c.Request.Context(), c.Param("id"))
	if err != nil {
		respondJSON(c, service.Error(service.ErrCodeContainerNotFound, err.Error()))
		return
	}
	respondJSON(c, service.Success(gin.H{"ports": detail.NetworkSettings.Ports}))
}

// Mounts GET /api/v1/containers/:id/mounts
func (h *ContainerHandler) Mounts(c *gin.Context) {
	detail, err := h.svc.ContainerInspect(c.Request.Context(), c.Param("id"))
	if err != nil {
		respondJSON(c, service.Error(service.ErrCodeContainerNotFound, err.Error()))
		return
	}
	respondJSON(c, service.Success(gin.H{"mounts": detail.Mounts}))
}

// Export GET /api/v1/containers/:id/export
func (h *ContainerHandler) Export(c *gin.Context) {
	rc, err := h.svc.ContainerExport(c.Request.Context(), c.Param("id"))
	if err != nil {
		respondJSON(c, service.Error(service.ErrCodeDockerAPI, err.Error()))
		return
	}
	defer rc.Close()

	c.Header("Content-Type", "application/x-tar")
	c.Header("Content-Disposition", "attachment; filename=container.tar")
	c.Status(http.StatusOK)
	io.Copy(c.Writer, rc)
}
