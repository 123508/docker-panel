package handler

import (
	service2 "docker-panel/internal/service"
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
	svc    *service2.ContainerService
	recent *service2.RecentContainers // 最近操作容器 LRU，用于仪表盘活跃容器列表
}

// NewContainerHandler 创建 ContainerHandler
func NewContainerHandler(svc *service2.ContainerService, recent *service2.RecentContainers) *ContainerHandler {
	return &ContainerHandler{svc: svc, recent: recent}
}

// List GET /api/v1/containers
func (h *ContainerHandler) List(c *gin.Context) {
	var req service2.ContainerListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		respondJSON(c, service2.Error(service2.ErrCodeInvalidParam, err.Error()))
		return
	}
	items, err := h.svc.ContainerList(c.Request.Context(), req)
	if err != nil {
		respondJSON(c, service2.Error(service2.ErrCodeDockerAPI, err.Error()))
		return
	}
	respondJSON(c, service2.Success(items))
}

// RecentContainers GET /api/v1/containers/recent
// 返回最近操作过的容器列表（LRU），用于仪表盘活跃容器展示。
func (h *ContainerHandler) RecentContainers(c *gin.Context) {
	ids := h.recent.GetIDs()
	if len(ids) == 0 {
		respondJSON(c, service2.Success([]service2.ContainerListItem{}))
		return
	}
	items, err := h.svc.GetContainersByIDs(c.Request.Context(), ids)
	if err != nil {
		respondJSON(c, service2.Error(service2.ErrCodeDockerAPI, err.Error()))
		return
	}
	respondJSON(c, service2.Success(items))
}

// Inspect GET /api/v1/containers/:id
func (h *ContainerHandler) Inspect(c *gin.Context) {
	detail, err := h.svc.ContainerInspect(c.Request.Context(), c.Param("id"))
	if err != nil {
		respondJSON(c, service2.Error(service2.ErrCodeContainerNotFound, err.Error()))
		return
	}
	respondJSON(c, service2.Success(detail))
}

// Create POST /api/v1/containers
func (h *ContainerHandler) Create(c *gin.Context) {
	var req service2.ContainerCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondJSON(c, service2.Error(service2.ErrCodeInvalidParam, err.Error()))
		return
	}
	id, err := h.svc.ContainerCreate(c.Request.Context(), req)
	if err != nil {
		respondJSON(c, service2.Error(service2.ErrCodeDockerAPI, err.Error()))
		return
	}
	h.recent.Record(id) // 将新创建的容器记录到 LRU 中
	respondJSON(c, service2.Success(gin.H{"id": id}))
}

// Start POST /api/v1/containers/:id/start
func (h *ContainerHandler) Start(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.ContainerStart(c.Request.Context(), id); err != nil {
		respondJSON(c, service2.Error(service2.ErrCodeDockerAPI, err.Error()))
		return
	}
	h.recent.Record(id) // 操作成功后将容器记录到 LRU 中
	respondJSON(c, service2.Success(nil))
}

// Stop POST /api/v1/containers/:id/stop
func (h *ContainerHandler) Stop(c *gin.Context) {
	var req service2.ContainerStopRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
	}
	id := c.Param("id")
	if err := h.svc.ContainerStop(c.Request.Context(), id, req.Timeout); err != nil {
		respondJSON(c, service2.Error(service2.ErrCodeDockerAPI, err.Error()))
		return
	}
	h.recent.Record(id) // 操作成功后将容器记录到 LRU 中
	respondJSON(c, service2.Success(nil))
}

// Kill POST /api/v1/containers/:id/kill
func (h *ContainerHandler) Kill(c *gin.Context) {
	var req service2.ContainerKillRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
	}
	id := c.Param("id")
	if err := h.svc.ContainerKill(c.Request.Context(), id, req.Signal); err != nil {
		respondJSON(c, service2.Error(service2.ErrCodeDockerAPI, err.Error()))
		return
	}
	h.recent.Record(id) // 操作成功后将容器记录到 LRU 中
	respondJSON(c, service2.Success(nil))
}

// Restart POST /api/v1/containers/:id/restart
func (h *ContainerHandler) Restart(c *gin.Context) {
	var req service2.ContainerRestartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
	}
	id := c.Param("id")
	if err := h.svc.ContainerRestart(c.Request.Context(), id, req.Timeout); err != nil {
		respondJSON(c, service2.Error(service2.ErrCodeDockerAPI, err.Error()))
		return
	}
	h.recent.Record(id) // 操作成功后将容器记录到 LRU 中
	respondJSON(c, service2.Success(nil))
}

// Pause POST /api/v1/containers/:id/pause
func (h *ContainerHandler) Pause(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.ContainerPause(c.Request.Context(), id); err != nil {
		respondJSON(c, service2.Error(service2.ErrCodeDockerAPI, err.Error()))
		return
	}
	h.recent.Record(id) // 操作成功后将容器记录到 LRU 中
	respondJSON(c, service2.Success(nil))
}

// Unpause POST /api/v1/containers/:id/unpause
func (h *ContainerHandler) Unpause(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.ContainerUnpause(c.Request.Context(), id); err != nil {
		respondJSON(c, service2.Error(service2.ErrCodeDockerAPI, err.Error()))
		return
	}
	h.recent.Record(id) // 操作成功后将容器记录到 LRU 中
	respondJSON(c, service2.Success(nil))
}

// Remove DELETE /api/v1/containers/:id
func (h *ContainerHandler) Remove(c *gin.Context) {
	var req service2.ContainerRemoveRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		_ = c.Error(err)
	}
	id := c.Param("id")
	if err := h.svc.ContainerRemove(c.Request.Context(), id, req.Force, req.RemoveVolumes); err != nil {
		respondJSON(c, service2.Error(service2.ErrCodeDockerAPI, err.Error()))
		return
	}
	h.recent.Remove(id) // 容器已删除，从 LRU 中移除记录
	respondJSON(c, service2.Success(nil))
}

// Rename POST /api/v1/containers/:id/rename
func (h *ContainerHandler) Rename(c *gin.Context) {
	var req service2.ContainerRenameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondJSON(c, service2.Error(service2.ErrCodeInvalidParam, err.Error()))
		return
	}
	id := c.Param("id")
	if err := h.svc.ContainerRename(c.Request.Context(), id, req.NewName); err != nil {
		respondJSON(c, service2.Error(service2.ErrCodeDockerAPI, err.Error()))
		return
	}
	h.recent.Record(id) // 操作成功后将容器记录到 LRU 中
	respondJSON(c, service2.Success(nil))
}

// Logs GET /api/v1/containers/:id/logs
func (h *ContainerHandler) Logs(c *gin.Context) {
	var opts service2.ContainerLogsOptions
	if err := c.ShouldBindQuery(&opts); err != nil {
		respondJSON(c, service2.Error(service2.ErrCodeInvalidParam, err.Error()))
		return
	}
	// 默认展示 stdout+stderr
	if !opts.ShowStdout && !opts.ShowStderr {
		opts.ShowStdout = true
		opts.ShowStderr = true
	}

	rc, err := h.svc.ContainerLogs(c.Request.Context(), c.Param("id"), opts)
	if err != nil {
		respondJSON(c, service2.Error(service2.ErrCodeDockerAPI, err.Error()))
		return
	}
	defer func() { _ = rc.Close() }()

	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.Header("Transfer-Encoding", "chunked")
	c.Status(http.StatusOK)
	if _, err := io.Copy(c.Writer, rc); err != nil {
		_ = c.Error(err)
	}
}

// Exec POST /api/v1/containers/:id/exec
func (h *ContainerHandler) Exec(c *gin.Context) {
	var req service2.ContainerExecRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondJSON(c, service2.Error(service2.ErrCodeInvalidParam, err.Error()))
		return
	}
	req.AttachStdout = true
	req.AttachStderr = true

	id := c.Param("id")
	hijack, err := h.svc.ContainerExec(c.Request.Context(), id, req)
	if err != nil {
		respondJSON(c, service2.Error(service2.ErrCodeDockerAPI, err.Error()))
		return
	}
	defer hijack.Close()

	output, err := io.ReadAll(hijack.Reader)
	if err != nil {
		respondJSON(c, service2.Error(service2.ErrCodeDockerAPI, err.Error()))
		return
	}
	h.recent.Record(id) // exec 成功后将容器记录到 LRU 中
	respondJSON(c, service2.Success(gin.H{"output": string(output)}))
}

// Terminal WebSocket /api/v1/containers/:id/terminal
func (h *ContainerHandler) Terminal(c *gin.Context) {
	containerID := c.Param("id")

	ws, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer func() { _ = ws.Close() }()

	ctx := c.Request.Context()

	// 创建 exec 会话（sh 优先，降级 bash）
	execID, err := h.svc.ContainerExecCreate(ctx, containerID, service2.ContainerExecRequest{
		Cmd: []string{"sh"},
	})
	if err != nil {
		_ = ws.WriteMessage(websocket.TextMessage, []byte("error: "+err.Error()))
		return
	}

	hijack, err := h.svc.ContainerExecAttach(ctx, execID)
	if err != nil {
		_ = ws.WriteMessage(websocket.TextMessage, []byte("error: "+err.Error()))
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
			if _, err := hijack.Conn.Write(msg); err != nil {
				return
			}
		}
	}()

	// Docker stdout -> WebSocket
	buf := make([]byte, 4096)
	for {
		n, err := hijack.Reader.Read(buf)
		if n > 0 {
			if err := ws.WriteMessage(websocket.BinaryMessage, buf[:n]); err != nil {
				break
			}
		}
		if err != nil {
			break
		}
	}
}

// Top GET /api/v1/containers/:id/top
func (h *ContainerHandler) Top(c *gin.Context) {
	var req service2.ContainerTopRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		_ = c.Error(err)
	}
	result, err := h.svc.ContainerTop(c.Request.Context(), c.Param("id"), req.PsArgs)
	if err != nil {
		respondJSON(c, service2.Error(service2.ErrCodeDockerAPI, err.Error()))
		return
	}
	respondJSON(c, service2.Success(result))
}

// Ports GET /api/v1/containers/:id/ports
func (h *ContainerHandler) Ports(c *gin.Context) {
	detail, err := h.svc.ContainerInspect(c.Request.Context(), c.Param("id"))
	if err != nil {
		respondJSON(c, service2.Error(service2.ErrCodeContainerNotFound, err.Error()))
		return
	}
	respondJSON(c, service2.Success(gin.H{"ports": detail.NetworkSettings.Ports}))
}

// Mounts GET /api/v1/containers/:id/mounts
func (h *ContainerHandler) Mounts(c *gin.Context) {
	detail, err := h.svc.ContainerInspect(c.Request.Context(), c.Param("id"))
	if err != nil {
		respondJSON(c, service2.Error(service2.ErrCodeContainerNotFound, err.Error()))
		return
	}
	respondJSON(c, service2.Success(gin.H{"mounts": detail.Mounts}))
}

// Export GET /api/v1/containers/:id/export
func (h *ContainerHandler) Export(c *gin.Context) {
	rc, err := h.svc.ContainerExport(c.Request.Context(), c.Param("id"))
	if err != nil {
		respondJSON(c, service2.Error(service2.ErrCodeDockerAPI, err.Error()))
		return
	}
	defer func() { _ = rc.Close() }()

	c.Header("Content-Type", "application/x-tar")
	c.Header("Content-Disposition", "attachment; filename=container.tar")
	c.Status(http.StatusOK)
	if _, err := io.Copy(c.Writer, rc); err != nil {
		_ = c.Error(err)
	}
}
