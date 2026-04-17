package handler

import (
	service2 "docker-panel/internal/service"

	"github.com/gin-gonic/gin"
)

// NetworkHandler 网络 HTTP 处理器
type NetworkHandler struct {
	svc *service2.NetworkService
}

// NewNetworkHandler 创建 NetworkHandler
func NewNetworkHandler(svc *service2.NetworkService) *NetworkHandler {
	return &NetworkHandler{svc: svc}
}

// List GET /api/v1/networks
func (h *NetworkHandler) List(c *gin.Context) {
	items, err := h.svc.NetworkList(c.Request.Context())
	if err != nil {
		respondJSON(c, service2.Error(service2.ErrCodeDockerAPI, err.Error()))
		return
	}
	respondJSON(c, service2.Success(items))
}

// Inspect GET /api/v1/networks/:id
func (h *NetworkHandler) Inspect(c *gin.Context) {
	n, err := h.svc.NetworkInspect(c.Request.Context(), c.Param("id"))
	if err != nil {
		respondJSON(c, service2.Error(service2.ErrCodeNetworkNotFound, err.Error()))
		return
	}
	respondJSON(c, service2.Success(n))
}

// Create POST /api/v1/networks
func (h *NetworkHandler) Create(c *gin.Context) {
	var req service2.NetworkCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondJSON(c, service2.Error(service2.ErrCodeInvalidParam, err.Error()))
		return
	}
	id, err := h.svc.NetworkCreate(c.Request.Context(), req)
	if err != nil {
		respondJSON(c, service2.Error(service2.ErrCodeDockerAPI, err.Error()))
		return
	}
	respondJSON(c, service2.Success(gin.H{"id": id}))
}

// Remove DELETE /api/v1/networks/:id
func (h *NetworkHandler) Remove(c *gin.Context) {
	if err := h.svc.NetworkRemove(c.Request.Context(), c.Param("id")); err != nil {
		respondJSON(c, service2.Error(service2.ErrCodeDockerAPI, err.Error()))
		return
	}
	respondJSON(c, service2.Success(nil))
}

// Connect POST /api/v1/networks/:id/connect
func (h *NetworkHandler) Connect(c *gin.Context) {
	var req service2.NetworkConnectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondJSON(c, service2.Error(service2.ErrCodeInvalidParam, err.Error()))
		return
	}
	if err := h.svc.NetworkConnect(c.Request.Context(), c.Param("id"), req.ContainerID); err != nil {
		respondJSON(c, service2.Error(service2.ErrCodeDockerAPI, err.Error()))
		return
	}
	respondJSON(c, service2.Success(nil))
}

// Disconnect POST /api/v1/networks/:id/disconnect
func (h *NetworkHandler) Disconnect(c *gin.Context) {
	var req service2.NetworkDisconnectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondJSON(c, service2.Error(service2.ErrCodeInvalidParam, err.Error()))
		return
	}
	if err := h.svc.NetworkDisconnect(c.Request.Context(), c.Param("id"), req.ContainerID, req.Force); err != nil {
		respondJSON(c, service2.Error(service2.ErrCodeDockerAPI, err.Error()))
		return
	}
	respondJSON(c, service2.Success(nil))
}
