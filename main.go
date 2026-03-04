package main

import (
	"docker-panel/internal/config"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"strings"
)

//go:embed web/dist
var staticFS embed.FS

func main() {
	if err := config.InitConfig(); err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}

	http.HandleFunc("/api/containers", containersHandler)
	http.HandleFunc("/api/images", imagesHandler)
	http.HandleFunc("/api/health", healthHandler)

	http.Handle("/", spaHandler())

	cfg := config.AppConfig
	addr := cfg.Server.BindIP + ":" + cfg.Server.BindPort

	fmt.Printf("\nDocker Panel 启动成功！\n")
	fmt.Printf("管理员账号: %s\n", cfg.User.AdminUsername)
	fmt.Printf("绑定地址: %s\n", addr)
	fmt.Printf("访问地址: http://localhost:%s\n\n", cfg.Server.BindPort)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}

func spaHandler() http.Handler {
	distFS, err := fs.Sub(staticFS, "web/dist")
	if err != nil {
		log.Fatal(err)
	}

	fileServer := http.FileServer(http.FS(distFS))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}

		path := strings.TrimPrefix(r.URL.Path, "/")
		if path == "" {
			path = "index.html"
		}

		_, err := distFS.Open(path)
		if err != nil {
			r.URL.Path = "/"
		}

		fileServer.ServeHTTP(w, r)
	})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cfg := config.AppConfig
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "ok",
		"message": "Docker Panel is running",
		"config": map[string]string{
			"admin_user": cfg.User.AdminUsername,
			"bind_addr":  cfg.Server.BindIP + ":" + cfg.Server.BindPort,
		},
	})
}

func containersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"containers": []map[string]string{
			{"id": "container-1", "name": "web-server", "status": "running"},
			{"id": "container-2", "name": "database", "status": "stopped"},
		},
	})
}

func imagesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"images": []map[string]string{
			{"id": "image-1", "name": "nginx:latest", "size": "142MB"},
			{"id": "image-2", "name": "postgres:14", "size": "376MB"},
		},
	})
}
