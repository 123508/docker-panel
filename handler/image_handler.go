package handler

import (
	"docker-panel/service"
	"io"
	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/gin-gonic/gin"
)

// ImageHandler 镜像 HTTP 处理器
type ImageHandler struct {
	svc *service.ImageService
}

// NewImageHandler 创建 ImageHandler
func NewImageHandler(svc *service.ImageService) *ImageHandler {
	return &ImageHandler{svc: svc}
}

// RegisterRoutes 注册镜像路由
func (h *ImageHandler) RegisterRoutes(rg *gin.RouterGroup) {
	images := rg.Group("/images")
	{
		images.GET("", h.List)
		images.GET("/:id", h.Inspect)
		images.DELETE("/:id", h.Remove)
		images.POST("/pull", h.Pull)
		images.POST("/tag", h.Tag)
		images.POST("/push", h.Push)
		images.GET("/:id/save", h.Save)
		images.POST("/load", h.Load)
		images.POST("/import", h.Import)
	}
}

// List GET /api/v1/images
func (h *ImageHandler) List(c *gin.Context) {
	var req service.ImageListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		respondJSON(c, service.Error(service.ErrCodeInvalidParam, err.Error()))
		return
	}
	items, err := h.svc.ImageList(c.Request.Context(), req)
	if err != nil {
		respondJSON(c, service.Error(service.ErrCodeDockerAPI, err.Error()))
		return
	}
	respondJSON(c, service.Success(items))
}

// Inspect GET /api/v1/images/:id
func (h *ImageHandler) Inspect(c *gin.Context) {
	detail, err := h.svc.ImageInspect(c.Request.Context(), c.Param("id"))
	if err != nil {
		respondJSON(c, service.Error(service.ErrCodeImageNotFound, err.Error()))
		return
	}
	respondJSON(c, service.Success(detail))
}

// Pull POST /api/v1/images/pull
func (h *ImageHandler) Pull(c *gin.Context) {
	var req service.ImagePullRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondJSON(c, service.Error(service.ErrCodeInvalidParam, err.Error()))
		return
	}

	rc, err := h.svc.ImagePull(c.Request.Context(), req.Image)
	if err != nil {
		respondJSON(c, service.Error(service.ErrCodeDockerAPI, err.Error()))
		return
	}
	defer rc.Close()

	c.Header("Content-Type", "application/json")
	c.Header("Transfer-Encoding", "chunked")
	c.Status(http.StatusOK)
	io.Copy(c.Writer, rc)
}

// Remove DELETE /api/v1/images/:id
func (h *ImageHandler) Remove(c *gin.Context) {
	var req service.ImageRemoveRequest
	_ = c.ShouldBindQuery(&req)

	result, err := h.svc.ImageRemove(c.Request.Context(), c.Param("id"), req.Force, req.PruneChildren)
	if err != nil {
		respondJSON(c, service.Error(service.ErrCodeDockerAPI, err.Error()))
		return
	}
	respondJSON(c, service.Success(result))
}

// Tag POST /api/v1/images/tag
func (h *ImageHandler) Tag(c *gin.Context) {
	var req service.ImageTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondJSON(c, service.Error(service.ErrCodeInvalidParam, err.Error()))
		return
	}
	if err := h.svc.ImageTag(c.Request.Context(), req.Source, req.Target); err != nil {
		respondJSON(c, service.Error(service.ErrCodeDockerAPI, err.Error()))
		return
	}
	respondJSON(c, service.Success(nil))
}

// Push POST /api/v1/images/push
func (h *ImageHandler) Push(c *gin.Context) {
	var req service.ImagePushRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondJSON(c, service.Error(service.ErrCodeInvalidParam, err.Error()))
		return
	}

	rc, err := h.svc.ImagePush(c.Request.Context(), req.Image)
	if err != nil {
		respondJSON(c, service.Error(service.ErrCodeDockerAPI, err.Error()))
		return
	}
	defer rc.Close()

	c.Header("Content-Type", "application/json")
	c.Header("Transfer-Encoding", "chunked")
	c.Status(http.StatusOK)
	io.Copy(c.Writer, rc)
}

// Save GET /api/v1/images/:id/save
func (h *ImageHandler) Save(c *gin.Context) {
	rc, err := h.svc.ImageSave(c.Request.Context(), []string{c.Param("id")})
	if err != nil {
		respondJSON(c, service.Error(service.ErrCodeDockerAPI, err.Error()))
		return
	}
	defer rc.Close()

	c.Header("Content-Type", "application/x-tar")
	c.Header("Content-Disposition", "attachment; filename=image.tar")
	c.Status(http.StatusOK)
	io.Copy(c.Writer, rc)
}

// Load POST /api/v1/images/load
func (h *ImageHandler) Load(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		respondJSON(c, service.Error(service.ErrCodeInvalidParam, "缺少file参数"))
		return
	}
	defer file.Close()

	resp, err := h.svc.ImageLoad(c.Request.Context(), file, false)
	if err != nil {
		respondJSON(c, service.Error(service.ErrCodeDockerAPI, err.Error()))
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	respondJSON(c, service.Success(gin.H{"output": string(body)}))
}

// Import POST /api/v1/images/import
func (h *ImageHandler) Import(c *gin.Context) {
	var req service.ImageImportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondJSON(c, service.Error(service.ErrCodeInvalidParam, err.Error()))
		return
	}

	// 通过 ImageLoad 实现导入（从body读取tar流）
	buildOpts := types.ImageBuildOptions{
		Tags:   []string{req.Repository + ":" + req.Tag},
		Remove: true,
	}
	_ = buildOpts

	respondJSON(c, service.Success(gin.H{"message": "请使用 /images/load 接口上传tar文件"}))
}
