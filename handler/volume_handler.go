package handler

import (
	"docker-panel/service"

	"github.com/gin-gonic/gin"
)

// VolumeHandler 数据卷 HTTP 处理器
type VolumeHandler struct {
	svc *service.VolumeService
}

// NewVolumeHandler 创建 VolumeHandler
func NewVolumeHandler(svc *service.VolumeService) *VolumeHandler {
	return &VolumeHandler{svc: svc}
}

// RegisterRoutes 注册数据卷路由
func (h *VolumeHandler) RegisterRoutes(rg *gin.RouterGroup) {
	volumes := rg.Group("/volumes")
	{
		volumes.GET("", h.List)
		volumes.POST("", h.Create)
		volumes.GET("/:name", h.Inspect)
		volumes.DELETE("/:name", h.Remove)
	}
}

// List GET /api/v1/volumes
func (h *VolumeHandler) List(c *gin.Context) {
	items, err := h.svc.VolumeList(c.Request.Context())
	if err != nil {
		respondJSON(c, service.Error(service.ErrCodeDockerAPI, err.Error()))
		return
	}
	respondJSON(c, service.Success(items))
}

// Inspect GET /api/v1/volumes/:name
func (h *VolumeHandler) Inspect(c *gin.Context) {
	v, err := h.svc.VolumeInspect(c.Request.Context(), c.Param("name"))
	if err != nil {
		respondJSON(c, service.Error(service.ErrCodeVolumeNotFound, err.Error()))
		return
	}
	respondJSON(c, service.Success(v))
}

// Create POST /api/v1/volumes
func (h *VolumeHandler) Create(c *gin.Context) {
	var req service.VolumeCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondJSON(c, service.Error(service.ErrCodeInvalidParam, err.Error()))
		return
	}
	v, err := h.svc.VolumeCreate(c.Request.Context(), req)
	if err != nil {
		respondJSON(c, service.Error(service.ErrCodeDockerAPI, err.Error()))
		return
	}
	respondJSON(c, service.Success(v))
}

// Remove DELETE /api/v1/volumes/:name
func (h *VolumeHandler) Remove(c *gin.Context) {
	var req service.VolumeRemoveRequest
	_ = c.ShouldBindQuery(&req)
	if err := h.svc.VolumeRemove(c.Request.Context(), c.Param("name"), req.Force); err != nil {
		respondJSON(c, service.Error(service.ErrCodeDockerAPI, err.Error()))
		return
	}
	respondJSON(c, service.Success(nil))
}
