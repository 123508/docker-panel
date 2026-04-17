package handler

import (
	"docker-panel/service"

	"github.com/gin-gonic/gin"
)

// NetworkHandler 网络 HTTP 处理器
type NetworkHandler struct {
	svc *service.NetworkService
}

// NewNetworkHandler 创建 NetworkHandler
func NewNetworkHandler(svc *service.NetworkService) *NetworkHandler {
	return &NetworkHandler{svc: svc}
}

// RegisterRoutes 注册网络路由
func (h *NetworkHandler) RegisterRoutes(rg *gin.RouterGroup) {
	networks := rg.Group("/networks")
	{
		networks.GET("", h.List)
		networks.POST("", h.Create)
		networks.GET("/:id", h.Inspect)
		networks.DELETE("/:id", h.Remove)
		networks.POST("/:id/connect", h.Connect)
		networks.POST("/:id/disconnect", h.Disconnect)
	}
}

// List GET /api/v1/networks
func (h *NetworkHandler) List(c *gin.Context) {
	items, err := h.svc.NetworkList(c.Request.Context())
	if err != nil {
		respondJSON(c, service.Error(service.ErrCodeDockerAPI, err.Error()))
		return
	}
	respondJSON(c, service.Success(items))
}

// Inspect GET /api/v1/networks/:id
func (h *NetworkHandler) Inspect(c *gin.Context) {
	n, err := h.svc.NetworkInspect(c.Request.Context(), c.Param("id"))
	if err != nil {
		respondJSON(c, service.Error(service.ErrCodeNetworkNotFound, err.Error()))
		return
	}
	respondJSON(c, service.Success(n))
}

// Create POST /api/v1/networks
func (h *NetworkHandler) Create(c *gin.Context) {
	var req service.NetworkCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondJSON(c, service.Error(service.ErrCodeInvalidParam, err.Error()))
		return
	}
	id, err := h.svc.NetworkCreate(c.Request.Context(), req)
	if err != nil {
		respondJSON(c, service.Error(service.ErrCodeDockerAPI, err.Error()))
		return
	}
	respondJSON(c, service.Success(gin.H{"id": id}))
}

// Remove DELETE /api/v1/networks/:id
func (h *NetworkHandler) Remove(c *gin.Context) {
	if err := h.svc.NetworkRemove(c.Request.Context(), c.Param("id")); err != nil {
		respondJSON(c, service.Error(service.ErrCodeDockerAPI, err.Error()))
		return
	}
	respondJSON(c, service.Success(nil))
}

// Connect POST /api/v1/networks/:id/connect
func (h *NetworkHandler) Connect(c *gin.Context) {
	var req service.NetworkConnectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondJSON(c, service.Error(service.ErrCodeInvalidParam, err.Error()))
		return
	}
	if err := h.svc.NetworkConnect(c.Request.Context(), c.Param("id"), req.ContainerID); err != nil {
		respondJSON(c, service.Error(service.ErrCodeDockerAPI, err.Error()))
		return
	}
	respondJSON(c, service.Success(nil))
}

// Disconnect POST /api/v1/networks/:id/disconnect
func (h *NetworkHandler) Disconnect(c *gin.Context) {
	var req service.NetworkDisconnectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondJSON(c, service.Error(service.ErrCodeInvalidParam, err.Error()))
		return
	}
	if err := h.svc.NetworkDisconnect(c.Request.Context(), c.Param("id"), req.ContainerID, req.Force); err != nil {
		respondJSON(c, service.Error(service.ErrCodeDockerAPI, err.Error()))
		return
	}
	respondJSON(c, service.Success(nil))
}
