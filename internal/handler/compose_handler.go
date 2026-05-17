package handler

import (
	service2 "docker-panel/internal/service"

	"github.com/gin-gonic/gin"
)

// ComposeHandler compose HTTP 处理器
type ComposeHandler struct {
	svc *service2.ComposeService
}

// NewComposeHandler 创建 ComposeHandler
func NewComposeHandler(svc *service2.ComposeService) *ComposeHandler {
	return &ComposeHandler{svc: svc}
}

// List GET /api/v1/compose/projects
func (h *ComposeHandler) List(c *gin.Context) {
	items, err := h.svc.ListProjects(c.Request.Context())
	if err != nil {
		respondJSON(c, service2.Error(service2.ErrCodeDockerAPI, err.Error()))
		return
	}
	respondJSON(c, service2.Success(items))
}

// Upload POST /api/v1/compose/projects
func (h *ComposeHandler) Upload(c *gin.Context) {
	var req service2.ComposeUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondJSON(c, service2.Error(service2.ErrCodeInvalidParam, err.Error()))
		return
	}
	if err := h.svc.UploadProject(c.Request.Context(), req); err != nil {
		respondJSON(c, service2.Error(service2.ErrCodeDockerAPI, err.Error()))
		return
	}
	respondJSON(c, service2.Success(nil))
}

// Up POST /api/v1/compose/projects/:name/up
func (h *ComposeHandler) Up(c *gin.Context) {
	var req service2.ComposeUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// 允许空 body
		req = service2.ComposeUpRequest{}
	}
	if err := h.svc.Up(c.Request.Context(), c.Param("name"), &req); err != nil {
		respondJSON(c, service2.Error(service2.ErrCodeDockerAPI, err.Error()))
		return
	}
	respondJSON(c, service2.Success(nil))
}

// Stop POST /api/v1/compose/projects/:name/stop
func (h *ComposeHandler) Stop(c *gin.Context) {
	var req service2.ComposeStopRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		req = service2.ComposeStopRequest{}
	}
	if err := h.svc.Stop(c.Request.Context(), c.Param("name"), &req); err != nil {
		respondJSON(c, service2.Error(service2.ErrCodeDockerAPI, err.Error()))
		return
	}
	respondJSON(c, service2.Success(nil))
}

// Restart POST /api/v1/compose/projects/:name/restart
func (h *ComposeHandler) Restart(c *gin.Context) {
	var req service2.ComposeRestartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		req = service2.ComposeRestartRequest{}
	}
	if err := h.svc.Restart(c.Request.Context(), c.Param("name"), &req); err != nil {
		respondJSON(c, service2.Error(service2.ErrCodeDockerAPI, err.Error()))
		return
	}
	respondJSON(c, service2.Success(nil))
}

// Down DELETE /api/v1/compose/projects/:name
func (h *ComposeHandler) Down(c *gin.Context) {
	var req service2.ComposeDownRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		req = service2.ComposeDownRequest{}
	}
	if err := h.svc.Down(c.Request.Context(), c.Param("name"), &req); err != nil {
		respondJSON(c, service2.Error(service2.ErrCodeDockerAPI, err.Error()))
		return
	}
	respondJSON(c, service2.Success(nil))
}

// PS GET /api/v1/compose/projects/:name/ps
func (h *ComposeHandler) PS(c *gin.Context) {
	items, err := h.svc.PS(c.Request.Context(), c.Param("name"))
	if err != nil {
		respondJSON(c, service2.Error(service2.ErrCodeDockerAPI, err.Error()))
		return
	}
	respondJSON(c, service2.Success(items))
}

// Logs GET /api/v1/compose/projects/:name/logs
func (h *ComposeHandler) Logs(c *gin.Context) {
	var req service2.ComposeLogsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		req = service2.ComposeLogsRequest{}
	}
	output, err := h.svc.Logs(c.Request.Context(), c.Param("name"), &req)
	if err != nil {
		respondJSON(c, service2.Error(service2.ErrCodeDockerAPI, err.Error()))
		return
	}
	respondJSON(c, service2.Success(gin.H{"logs": output}))
}

// Scale POST /api/v1/compose/projects/:name/scale
func (h *ComposeHandler) Scale(c *gin.Context) {
	var req service2.ComposeScaleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondJSON(c, service2.Error(service2.ErrCodeInvalidParam, err.Error()))
		return
	}
	if err := h.svc.Scale(c.Request.Context(), c.Param("name"), req); err != nil {
		respondJSON(c, service2.Error(service2.ErrCodeDockerAPI, err.Error()))
		return
	}
	respondJSON(c, service2.Success(nil))
}
