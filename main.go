package main

import (
	"context"
	"docker-panel/internal/config"
	"docker-panel/internal/docker"
	"docker-panel/internal/handler"
	"docker-panel/internal/service"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"

	"github.com/docker/docker/api/types"

	"github.com/gin-gonic/gin"
)

//go:embed web/dist
var staticFS embed.FS

func printMsg(version types.Version, cfg *config.Config, addr string, r *gin.Engine, debug bool) {

	if !debug {
		return
	}

	fmt.Printf("Docker Engine 版本: %s\n", version.Version)
	fmt.Printf("Docker API 版本: %s\n", version.APIVersion)
	fmt.Printf("编译 Docker 的 Go 版本: %s\n", version.GoVersion)
	fmt.Printf("系统信息: %s/%s\n", version.Os, version.Arch)
	fmt.Printf("\nDocker Panel 启动成功！\n")
	fmt.Printf("管理员账号: %s\n", cfg.User.AdminUsername)
	fmt.Printf("绑定地址: %s\n", addr)
	fmt.Printf("访问地址: http://localhost:%s\n\n", cfg.Server.BindPort)

	for _, route := range r.Routes() {
		fmt.Printf("%-6s %-30s %s\n", route.Method, route.Path, route.Handler)
	}
}

func main() {
	if err := config.InitConfig(); err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}

	// 初始化 Docker 客户端
	dockerClient, err := docker.NewDockerClient()
	if err != nil {
		log.Fatalf("Failed to connect to Docker: %v", err)
	}
	defer dockerClient.Close()

	// 初始化 Handler 层
	deps := handler.Dependencies{
		ContainerSvc: service.NewContainerService(dockerClient),
		ImageSvc:     service.NewImageService(dockerClient),
		VolumeSvc:    service.NewVolumeService(dockerClient),
		NetworkSvc:   service.NewNetworkService(dockerClient),
	}

	// 初始化路由
	r := handler.NewRouter(deps)
	setupStaticFiles(r)

	cfg := config.AppConfig
	addr := cfg.Server.BindIP + ":" + cfg.Server.BindPort
	version, err := dockerClient.Version(context.Background())
	if err != nil {
		log.Fatalf("Failed to get Docker version: %v", err)
	}

	printMsg(version, cfg, addr, r, false)

	if err := r.Run(addr); err != nil {
		log.Fatal(err)
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
