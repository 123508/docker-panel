package service

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// ========== 数据结构定义 ==========

// ComposeServiceItem compose 项目中的服务信息
type ComposeServiceItem struct {
	Name       string   `json:"name"`
	Image      string   `json:"image"`
	Containers []string `json:"containers"`
	Replicas   int      `json:"replicas"`
	Status     string   `json:"status"`
}

// ComposeProjectItem compose 项目列表项
type ComposeProjectItem struct {
	Name      string               `json:"name"`
	FilePath  string               `json:"file_path"`
	Status    string               `json:"status"`
	Services  []ComposeServiceItem `json:"services"`
	CreatedAt string               `json:"created_at"`
	UpdatedAt string               `json:"updated_at"`
}

// ComposeUploadRequest 上传 compose 文件请求
type ComposeUploadRequest struct {
	Name    string `json:"name" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// ComposeUpRequest 启动请求
type ComposeUpRequest struct {
	Build         bool     `json:"build"`
	ForceRecreate bool     `json:"force_recreate"`
	NoRecreate    bool     `json:"no_recreate"`
	NoBuild       bool     `json:"no_build"`
	Timeout       int      `json:"timeout"`
	RemoveOrphans bool     `json:"remove_orphans"`
	Services      []string `json:"services"`
}

// ComposeStopRequest 停止请求
type ComposeStopRequest struct {
	Timeout  int      `json:"timeout"`
	Services []string `json:"services"`
}

// ComposeRestartRequest 重启请求
type ComposeRestartRequest struct {
	Timeout  int      `json:"timeout"`
	Services []string `json:"services"`
}

// ComposeDownRequest 删除请求
type ComposeDownRequest struct {
	RemoveVolumes bool   `json:"remove_volumes"`
	RemoveImages  string `json:"remove_images"`
	RemoveOrphans bool   `json:"remove_orphans"`
	Timeout       int    `json:"timeout"`
}

// ComposeScaleRequest 扩缩容请求
type ComposeScaleRequest struct {
	Services map[string]int `json:"services" binding:"required"`
}

// ComposeLogsRequest 日志请求
type ComposeLogsRequest struct {
	Follow     bool     `json:"follow" form:"follow"`
	Tail       string   `json:"tail" form:"tail"`
	Since      string   `json:"since" form:"since"`
	Until      string   `json:"until" form:"until"`
	Timestamps bool     `json:"timestamps" form:"timestamps"`
	Services   []string `json:"services" form:"services"`
}

// projectRecord 持久化的项目记录
type projectRecord struct {
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// ========== ComposeService ==========

// ComposeService compose 业务逻辑
type ComposeService struct {
	projectsDir string
	mu          sync.RWMutex
}

// NewComposeService 创建 ComposeService
func NewComposeService() *ComposeService {
	dir := "compose_projects"
	return &ComposeService{projectsDir: dir}
}

// projectsFilePath 返回项目注册表 JSON 文件路径
func (s *ComposeService) projectsFilePath() string {
	return filepath.Join(s.projectsDir, "projects.json")
}

// projectDir 返回单个项目的目录路径
func (s *ComposeService) projectDir(name string) string {
	return filepath.Join(s.projectsDir, name)
}

// composeFilePath 返回项目的 compose 文件路径
func (s *ComposeService) composeFilePath(name string) string {
	return filepath.Join(s.projectDir(name), "docker-compose.yml")
}

// ========== 项目注册表操作 ==========

func (s *ComposeService) loadRegistry() ([]projectRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	p := s.projectsFilePath()
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return []projectRecord{}, nil
	}
	data, err := os.ReadFile(p)
	if err != nil {
		return nil, fmt.Errorf("读取项目注册表: %w", err)
	}
	var records []projectRecord
	if len(data) == 0 {
		return []projectRecord{}, nil
	}
	if err := json.Unmarshal(data, &records); err != nil {
		return nil, fmt.Errorf("解析项目注册表: %w", err)
	}
	return records, nil
}

func (s *ComposeService) saveRegistry(records []projectRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := os.MkdirAll(s.projectsDir, 0755); err != nil {
		return fmt.Errorf("创建 compose 项目目录: %w", err)
	}
	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化项目注册表: %w", err)
	}
	if err := os.WriteFile(s.projectsFilePath(), data, 0644); err != nil {
		return fmt.Errorf("写入项目注册表: %w", err)
	}
	return nil
}

// ========== docker compose 命令执行 ==========

func (s *ComposeService) runCompose(ctx context.Context, projectName string, args ...string) ([]byte, error) {
	dir := s.projectDir(projectName)
	cmd := exec.CommandContext(ctx, "docker", append([]string{"compose"}, args...)...)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("docker compose %s: %s (%w)", strings.Join(args, " "), string(output), err)
	}
	return output, nil
}

// parseComposePSOutput 解析 docker compose ps --format json 的 JSON Lines 输出
func parseComposePSOutput(data []byte) ([]ComposeServiceItem, error) {
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	items := make([]ComposeServiceItem, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var item ComposeServiceItem
		if err := json.Unmarshal([]byte(line), &item); err != nil {
			continue
		}
		items = append(items, item)
	}
	return items, nil
}

// ========== 公开方法 ==========

// ListProjects 列出所有 compose 项目
func (s *ComposeService) ListProjects(ctx context.Context) ([]ComposeProjectItem, error) {
	records, err := s.loadRegistry()
	if err != nil {
		return nil, err
	}

	items := make([]ComposeProjectItem, 0, len(records))
	for _, rec := range records {
		item := ComposeProjectItem{
			Name:      rec.Name,
			FilePath:  s.composeFilePath(rec.Name),
			Status:    "stopped",
			Services:  []ComposeServiceItem{},
			CreatedAt: rec.CreatedAt,
			UpdatedAt: rec.UpdatedAt,
		}

		psOutput, err := s.runCompose(ctx, rec.Name, "-p", rec.Name, "ps", "--format", "json")
		if err == nil && len(psOutput) > 0 {
			if svcList, err := parseComposePSOutput(psOutput); err == nil {
				item.Services = svcList
				running := 0
				for _, svc := range svcList {
					if strings.Contains(strings.ToLower(svc.Status), "running") ||
						strings.Contains(strings.ToLower(svc.Status), "up") {
						running++
					}
				}
				if running > 0 {
					item.Status = "running"
				}
			}
		}

		items = append(items, item)
	}
	return items, nil
}

// UploadProject 上传 compose 文件并注册项目
func (s *ComposeService) UploadProject(ctx context.Context, req ComposeUploadRequest) error {
	records, err := s.loadRegistry()
	if err != nil {
		return err
	}

	for _, r := range records {
		if r.Name == req.Name {
			return fmt.Errorf("项目 %s 已存在", req.Name)
		}
	}

	dir := s.projectDir(req.Name)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建项目目录: %w", err)
	}

	if err := os.WriteFile(s.composeFilePath(req.Name), []byte(req.Content), 0644); err != nil {
		return fmt.Errorf("写入 compose 文件: %w", err)
	}

	now := time.Now().Format(time.RFC3339)
	records = append(records, projectRecord{
		Name:      req.Name,
		CreatedAt: now,
		UpdatedAt: now,
	})

	return s.saveRegistry(records)
}

// Up 启动 compose 项目
func (s *ComposeService) Up(ctx context.Context, name string, req *ComposeUpRequest) error {
	args := []string{"-p", name, "up", "-d"}

	if req != nil {
		if req.Build {
			args = append(args, "--build")
		}
		if req.ForceRecreate {
			args = append(args, "--force-recreate")
		}
		if req.NoRecreate {
			args = append(args, "--no-recreate")
		}
		if req.NoBuild {
			args = append(args, "--no-build")
		}
		if req.RemoveOrphans {
			args = append(args, "--remove-orphans")
		}
		if req.Timeout > 0 {
			args = append(args, "--timeout", fmt.Sprintf("%d", req.Timeout))
		}
	}
	if req != nil && len(req.Services) > 0 {
		args = append(args, req.Services...)
	}

	if _, err := s.runCompose(ctx, name, args...); err != nil {
		return err
	}

	return s.touchProject(name)
}

// Stop 停止 compose 项目
func (s *ComposeService) Stop(ctx context.Context, name string, req *ComposeStopRequest) error {
	args := []string{"-p", name, "stop"}

	if req != nil && req.Timeout > 0 {
		args = append(args, "--timeout", fmt.Sprintf("%d", req.Timeout))
	}
	if req != nil && len(req.Services) > 0 {
		args = append(args, req.Services...)
	}

	if _, err := s.runCompose(ctx, name, args...); err != nil {
		return err
	}
	return nil
}

// Restart 重启 compose 项目
func (s *ComposeService) Restart(ctx context.Context, name string, req *ComposeRestartRequest) error {
	args := []string{"-p", name, "restart"}

	if req != nil && req.Timeout > 0 {
		args = append(args, "--timeout", fmt.Sprintf("%d", req.Timeout))
	}
	if req != nil && len(req.Services) > 0 {
		args = append(args, req.Services...)
	}

	if _, err := s.runCompose(ctx, name, args...); err != nil {
		return err
	}
	return nil
}

// Down 删除 compose 项目（docker compose down + 移除注册表记录）
func (s *ComposeService) Down(ctx context.Context, name string, req *ComposeDownRequest) error {
	dir := s.projectDir(name)
	if _, err := os.Stat(dir); err == nil {
		args := []string{"-p", name, "down"}

		if req != nil {
			if req.RemoveVolumes {
				args = append(args, "--volumes")
			}
			if req.RemoveImages != "" {
				args = append(args, "--rmi", req.RemoveImages)
			}
			if req.RemoveOrphans {
				args = append(args, "--remove-orphans")
			}
			if req.Timeout > 0 {
				args = append(args, "--timeout", fmt.Sprintf("%d", req.Timeout))
			}
		}

		if _, err := s.runCompose(ctx, name, args...); err != nil {
			return err
		}
		_ = os.RemoveAll(dir)
	}

	records, err := s.loadRegistry()
	if err != nil {
		return err
	}

	filtered := make([]projectRecord, 0, len(records))
	for _, r := range records {
		if r.Name != name {
			filtered = append(filtered, r)
		}
	}

	return s.saveRegistry(filtered)
}

// PS 查看 compose 项目服务状态
func (s *ComposeService) PS(ctx context.Context, name string) ([]ComposeServiceItem, error) {
	output, err := s.runCompose(ctx, name, "-p", name, "ps", "--format", "json")
	if err != nil {
		return nil, err
	}

	return parseComposePSOutput(output)
}

// Logs 查看 compose 项目日志
func (s *ComposeService) Logs(ctx context.Context, name string, req *ComposeLogsRequest) (string, error) {
	args := []string{"-p", name, "logs"}

	if req != nil {
		if req.Follow {
			args = append(args, "--follow")
		}
		if req.Tail != "" {
			args = append(args, "--tail", req.Tail)
		}
		if req.Since != "" {
			args = append(args, "--since", req.Since)
		}
		if req.Until != "" {
			args = append(args, "--until", req.Until)
		}
		if req.Timestamps {
			args = append(args, "--timestamps")
		}
		if len(req.Services) > 0 {
			args = append(args, req.Services...)
		}
	}

	output, err := s.runCompose(ctx, name, args...)
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// Scale 扩缩容 compose 服务
func (s *ComposeService) Scale(ctx context.Context, name string, req ComposeScaleRequest) error {
	args := []string{"-p", name, "up", "-d"}

	for svc, count := range req.Services {
		args = append(args, "--scale", fmt.Sprintf("%s=%d", svc, count))
	}

	if _, err := s.runCompose(ctx, name, args...); err != nil {
		return err
	}
	return s.touchProject(name)
}

// touchProject 更新项目的 updated_at 时间戳
func (s *ComposeService) touchProject(name string) error {
	records, err := s.loadRegistry()
	if err != nil {
		return err
	}

	for i := range records {
		if records[i].Name == name {
			records[i].UpdatedAt = time.Now().Format(time.RFC3339)
			return s.saveRegistry(records)
		}
	}
	return nil
}
