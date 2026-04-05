package service

import (
	"context"
	"docker-panel/docker_cli_wrapper"
	"fmt"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/volume"
)

// ========== 数据结构定义 ==========

// VolumeListItem 数据卷列表项
type VolumeListItem struct {
	Name       string            `json:"name"`
	Driver     string            `json:"driver"`
	Mountpoint string            `json:"mountpoint"`
	CreatedAt  string            `json:"created_at"`
	Status     map[string]string `json:"status"`
	Labels     map[string]string `json:"labels"`
	Scope      string            `json:"scope"`
}

// VolumeCreateRequest 创建数据卷请求
type VolumeCreateRequest struct {
	Name       string            `json:"name"`
	Driver     string            `json:"driver"`
	DriverOpts map[string]string `json:"driver_opts"`
	Labels     map[string]string `json:"labels"`
}

// VolumeRemoveRequest 删除数据卷请求
type VolumeRemoveRequest struct {
	Force bool `form:"force"`
}

// ========== VolumeService ==========

// VolumeService 数据卷业务逻辑
type VolumeService struct {
	docker *docker_cli_wrapper.DockerClient
}

// NewVolumeService 创建 VolumeService
func NewVolumeService(docker *docker_cli_wrapper.DockerClient) *VolumeService {
	return &VolumeService{docker: docker}
}

// VolumeList 列出数据卷
func (s *VolumeService) VolumeList(ctx context.Context) ([]VolumeListItem, error) {
	resp, err := s.docker.VolumeList(ctx, filters.NewArgs())
	if err != nil {
		return nil, fmt.Errorf("docker list volumes: %w", err)
	}

	items := make([]VolumeListItem, 0, len(resp.Volumes))
	for _, v := range resp.Volumes {
		items = append(items, VolumeListItem{
			Name:       v.Name,
			Driver:     v.Driver,
			Mountpoint: v.Mountpoint,
			CreatedAt:  v.CreatedAt,
			Labels:     v.Labels,
			Scope:      v.Scope,
		})
	}
	return items, nil
}

// VolumeInspect 获取数据卷详情
func (s *VolumeService) VolumeInspect(ctx context.Context, volumeID string) (*volume.Volume, error) {
	v, err := s.docker.VolumeInspect(ctx, volumeID)
	if err != nil {
		return nil, fmt.Errorf("docker inspect volume: %w", err)
	}
	return &v, nil
}

// VolumeCreate 创建数据卷
func (s *VolumeService) VolumeCreate(ctx context.Context, req VolumeCreateRequest) (*volume.Volume, error) {
	v, err := s.docker.VolumeCreate(ctx, volume.CreateOptions{
		Name:       req.Name,
		Driver:     req.Driver,
		DriverOpts: req.DriverOpts,
		Labels:     req.Labels,
	})
	if err != nil {
		return nil, fmt.Errorf("docker create volume: %w", err)
	}
	return &v, nil
}

// VolumeRemove 删除数据卷
func (s *VolumeService) VolumeRemove(ctx context.Context, volumeID string, force bool) error {
	if err := s.docker.VolumeRemove(ctx, volumeID, force); err != nil {
		return fmt.Errorf("docker remove volume: %w", err)
	}
	return nil
}
