package handler

import (
	service2 "docker-panel/internal/service"

	"github.com/gin-gonic/gin"
)

// VolumeContainersResponse 卷关联容器响应
type VolumeContainersResponse struct {
	VolumeName string   `json:"volume_name"`
	Containers []string `json:"containers"`
}

// VolumeHandler 数据卷 HTTP 处理器
type VolumeHandler struct {
	svc *service2.VolumeService
}

// NewVolumeHandler 创建 VolumeHandler
func NewVolumeHandler(svc *service2.VolumeService) *VolumeHandler {
	return &VolumeHandler{svc: svc}
}

// List GET /api/v1/volumes
func (h *VolumeHandler) List(c *gin.Context) {
	items, err := h.svc.VolumeList(c.Request.Context())
	if err != nil {
		respondJSON(c, service2.Error(service2.ErrCodeDockerAPI, err.Error()))
		return
	}
	respondJSON(c, service2.Success(items))
}

// Inspect GET /api/v1/volumes/:name
func (h *VolumeHandler) Inspect(c *gin.Context) {
	v, err := h.svc.VolumeInspect(c.Request.Context(), c.Param("name"))
	if err != nil {
		respondJSON(c, service2.Error(service2.ErrCodeVolumeNotFound, err.Error()))
		return
	}
	respondJSON(c, service2.Success(v))
}

// Create POST /api/v1/volumes
func (h *VolumeHandler) Create(c *gin.Context) {
	var req service2.VolumeCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondJSON(c, service2.Error(service2.ErrCodeInvalidParam, err.Error()))
		return
	}
	v, err := h.svc.VolumeCreate(c.Request.Context(), req)
	if err != nil {
		respondJSON(c, service2.Error(service2.ErrCodeDockerAPI, err.Error()))
		return
	}
	respondJSON(c, service2.Success(v))
}

// Containers GET /api/v1/volumes/:name/containers
func (h *VolumeHandler) Containers(c *gin.Context) {
	names, err := h.svc.VolumeContainers(c.Request.Context(), c.Param("name"))
	if err != nil {
		respondJSON(c, service2.Error(service2.ErrCodeDockerAPI, err.Error()))
		return
	}
	respondJSON(c, service2.Success(VolumeContainersResponse{
		VolumeName: c.Param("name"),
		Containers: names,
	}))
}

// Remove DELETE /api/v1/volumes/:name
func (h *VolumeHandler) Remove(c *gin.Context) {
	var req service2.VolumeRemoveRequest
	_ = c.ShouldBindQuery(&req)
	if err := h.svc.VolumeRemove(c.Request.Context(), c.Param("name"), req.Force); err != nil {
		respondJSON(c, service2.Error(service2.ErrCodeDockerAPI, err.Error()))
		return
	}
	respondJSON(c, service2.Success(nil))
}
