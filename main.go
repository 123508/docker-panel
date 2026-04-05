package main

import (
	"docker-panel/docker_cli_wrapper"
	"docker-panel/handler"
	"docker-panel/internal/config"
	"docker-panel/service"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed web/dist
var staticFS embed.FS

func main() {
	if err := config.InitConfig(); err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}

	// 初始化 Docker 客户端
	dockerClient, err := docker_cli_wrapper.NewDockerClient()
	if err != nil {
		log.Fatalf("Failed to connect to Docker: %v", err)
	}
	defer dockerClient.Close()

	// 初始化 Service 层
	containerSvc := service.NewContainerService(dockerClient)
	imageSvc := service.NewImageService(dockerClient)
	volumeSvc := service.NewVolumeService(dockerClient)
	networkSvc := service.NewNetworkService(dockerClient)

	// 初始化 Handler 层
	containerHandler := handler.NewContainerHandler(containerSvc)
	imageHandler := handler.NewImageHandler(imageSvc)
	volumeHandler := handler.NewVolumeHandler(volumeSvc)
	networkHandler := handler.NewNetworkHandler(networkSvc)

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(handler.RecoveryMiddleware())
	r.Use(handler.LoggerMiddleware())
	r.Use(handler.CORSMiddleware())

	registerAPIRoutes(r, containerHandler, imageHandler, volumeHandler, networkHandler)
	setupStaticFiles(r)

	cfg := config.AppConfig
	addr := cfg.Server.BindIP + ":" + cfg.Server.BindPort

	fmt.Printf("\nDocker Panel 启动成功！\n")
	fmt.Printf("管理员账号: %s\n", cfg.User.AdminUsername)
	fmt.Printf("绑定地址: %s\n", addr)
	fmt.Printf("访问地址: http://localhost:%s\n\n", cfg.Server.BindPort)

	for _, route := range r.Routes() {
		fmt.Printf("%-6s %-30s %s\n", route.Method, route.Path, route.Handler)
	}

	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}

func registerAPIRoutes(
	r *gin.Engine,
	containerH *handler.ContainerHandler,
	imageH *handler.ImageHandler,
	volumeH *handler.VolumeHandler,
	networkH *handler.NetworkHandler,
) {
	api := r.Group("/api")
	{
		api.GET("/health", healthHandler)
	}

	v1 := api.Group("/v1")
	{
		containerH.RegisterRoutes(v1)
		imageH.RegisterRoutes(v1)
		volumeH.RegisterRoutes(v1)
		networkH.RegisterRoutes(v1)
	}
}

func setupStaticFiles(r *gin.Engine) {
	distFS, err := fs.Sub(staticFS, "web/dist")
	if err != nil {
		log.Fatal(err)
	}

	fileServer := http.FileServer(http.FS(distFS))

	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path[1:]
		if path == "" {
			path = "index.html"
		}

		if !fileExists(distFS, path) {
			c.Request.URL.Path = "/"
		}

		fileServer.ServeHTTP(c.Writer, c.Request)
	})
}

func fileExists(fsys fs.FS, path string) bool {
	if path == "" {
		path = "index.html"
	}
	_, err := fsys.Open(path)
	return err == nil
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
