package handler

import (
	"docker-panel/internal/config"
	"docker-panel/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Dependencies 持有所有路由注册所需的依赖
type Dependencies struct {
	ContainerSvc *service.ContainerService
	ImageSvc     *service.ImageService
	VolumeSvc    *service.VolumeService
	NetworkSvc   *service.NetworkService
}

// NewRouter 创建并配置 gin.Engine，注册所有中间件和 HTTP 路由
func NewRouter(deps Dependencies) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(RecoveryMiddleware())
	r.Use(LoggerMiddleware())
	r.Use(CORSMiddleware())

	containerH := NewContainerHandler(deps.ContainerSvc)
	imageH := NewImageHandler(deps.ImageSvc)
	volumeH := NewVolumeHandler(deps.VolumeSvc)
	networkH := NewNetworkHandler(deps.NetworkSvc)

	apiGroup := r.Group("/api")
	{
		apiGroup.GET("/health", healthHandler)
	}

	v1 := apiGroup.Group("/v1")
	{
		RegisterContainerRoutes(v1, containerH)
		RegisterImageRoutes(v1, imageH)
		RegisterVolumeRoutes(v1, volumeH)
		RegisterNetworkRoutes(v1, networkH)
	}

	return r
}

// RegisterContainerRoutes 注册容器路由
func RegisterContainerRoutes(rg *gin.RouterGroup, h *ContainerHandler) {
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

// RegisterImageRoutes 注册镜像路由
func RegisterImageRoutes(rg *gin.RouterGroup, h *ImageHandler) {
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

// RegisterVolumeRoutes 注册数据卷路由
func RegisterVolumeRoutes(rg *gin.RouterGroup, h *VolumeHandler) {
	volumes := rg.Group("/volumes")
	{
		volumes.GET("", h.List)
		volumes.POST("", h.Create)
		volumes.GET("/:name", h.Inspect)
		volumes.DELETE("/:name", h.Remove)
	}
}

// RegisterNetworkRoutes 注册网络路由
func RegisterNetworkRoutes(rg *gin.RouterGroup, h *NetworkHandler) {
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

func healthHandler(c *gin.Context) {
	cfg := config.AppConfig
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Docker Panel is running",
		"config": gin.H{
			"admin_user": cfg.User.AdminUsername,
			"bind_addr":  cfg.Server.BindIP + ":" + cfg.Server.BindPort,
		},
	})
}
