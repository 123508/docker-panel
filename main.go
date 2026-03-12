package main

import (
	"docker-panel/internal/config"
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

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	registerAPIRoutes(r)
	setupStaticFiles(r)

	cfg := config.AppConfig
	addr := cfg.Server.BindIP + ":" + cfg.Server.BindPort

	fmt.Printf("\nDocker Panel 启动成功！\n")
	fmt.Printf("管理员账号: %s\n", cfg.User.AdminUsername)
	fmt.Printf("绑定地址: %s\n", addr)
	fmt.Printf("访问地址: http://localhost:%s\n\n", cfg.Server.BindPort)

	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}

func registerAPIRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.GET("/health", healthHandler)
		api.GET("/containers", containersHandler)
		api.GET("/images", imagesHandler)
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

func containersHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"containers": []gin.H{
			{"id": "container-1", "name": "web-server", "status": "running"},
			{"id": "container-2", "name": "database", "status": "stopped"},
		},
	})
}

func imagesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"images": []gin.H{
			{"id": "image-1", "name": "nginx:latest", "size": "142MB"},
			{"id": "image-2", "name": "postgres:14", "size": "376MB"},
		},
	})
}
