package service

import (
	"context"
	"docker-panel/internal/docker"
	"fmt"

	"github.com/docker/docker/api/types/network"
)

// ========== 数据结构定义 ==========

// NetworkListItem 网络列表项
type NetworkListItem struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	Driver     string            `json:"driver"`
	Scope      string            `json:"scope"`
	IPAM       NetworkIPAM       `json:"ipam"`
	Internal   bool              `json:"internal"`
	Attachable bool              `json:"attachable"`
	Ingress    bool              `json:"ingress"`
	Labels     map[string]string `json:"labels"`
	Created    string            `json:"created"`
}

// NetworkIPAM IPAM配置
type NetworkIPAM struct {
	Driver string       `json:"driver"`
	Config []IPAMConfig `json:"config"`
}

// IPAMConfig IPAM子网配置
type IPAMConfig struct {
	Subnet  string `json:"subnet"`
	Gateway string `json:"gateway"`
}

// NetworkCreateRequest 创建网络请求
type NetworkCreateRequest struct {
	Name       string            `json:"name" binding:"required"`
	Driver     string            `json:"driver"`
	Internal   bool              `json:"internal"`
	Attachable bool              `json:"attachable"`
	Labels     map[string]string `json:"labels"`
	IPAM       *NetworkIPAM      `json:"ipam"`
}

// NetworkConnectRequest 连接容器到网络请求
type NetworkConnectRequest struct {
	ContainerID string `json:"container_id" binding:"required"`
}

// NetworkDisconnectRequest 断开容器与网络请求
type NetworkDisconnectRequest struct {
	ContainerID string `json:"container_id" binding:"required"`
	Force       bool   `json:"force"`
}

// ========== NetworkService ==========

// NetworkService 网络业务逻辑
type NetworkService struct {
	docker docker.DockerClientInterface
}

// NewNetworkService 创建 NetworkService
func NewNetworkService(docker docker.DockerClientInterface) *NetworkService {
	return &NetworkService{docker: docker}
}

// NetworkList 列出网络
func (s *NetworkService) NetworkList(ctx context.Context) ([]NetworkListItem, error) {
	networks, err := s.docker.NetworkList(ctx, network.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("docker list networks: %w", err)
	}

	items := make([]NetworkListItem, 0, len(networks))
	for _, n := range networks {
		item := NetworkListItem{
			ID:         n.ID,
			Name:       n.Name,
			Driver:     n.Driver,
			Scope:      n.Scope,
			Internal:   n.Internal,
			Attachable: n.Attachable,
			Ingress:    n.Ingress,
			Labels:     n.Labels,
			Created:    n.Created.String(),
			IPAM: NetworkIPAM{
				Driver: n.IPAM.Driver,
			},
		}
		for _, cfg := range n.IPAM.Config {
			item.IPAM.Config = append(item.IPAM.Config, IPAMConfig{
				Subnet:  cfg.Subnet,
				Gateway: cfg.Gateway,
			})
		}
		items = append(items, item)
	}
	return items, nil
}

// NetworkInspect 获取网络详情
func (s *NetworkService) NetworkInspect(ctx context.Context, networkID string) (*NetworkListItem, error) {
	n, err := s.docker.NetworkInspect(ctx, networkID, network.InspectOptions{})
	if err != nil {
		return nil, fmt.Errorf("docker inspect network: %w", err)
	}

	item := &NetworkListItem{
		ID:         n.ID,
		Name:       n.Name,
		Driver:     n.Driver,
		Scope:      n.Scope,
		Internal:   n.Internal,
		Attachable: n.Attachable,
		Ingress:    n.Ingress,
		Labels:     n.Labels,
		Created:    n.Created.String(),
		IPAM: NetworkIPAM{
			Driver: n.IPAM.Driver,
		},
	}
	for _, cfg := range n.IPAM.Config {
		item.IPAM.Config = append(item.IPAM.Config, IPAMConfig{
			Subnet:  cfg.Subnet,
			Gateway: cfg.Gateway,
		})
	}
	return item, nil
}

// NetworkCreate 创建网络
func (s *NetworkService) NetworkCreate(ctx context.Context, req NetworkCreateRequest) (string, error) {
	opts := network.CreateOptions{
		Driver:     req.Driver,
		Internal:   req.Internal,
		Attachable: req.Attachable,
		Labels:     req.Labels,
	}

	if req.IPAM != nil {
		ipamCfg := make([]network.IPAMConfig, 0, len(req.IPAM.Config))
		for _, c := range req.IPAM.Config {
			ipamCfg = append(ipamCfg, network.IPAMConfig{
				Subnet:  c.Subnet,
				Gateway: c.Gateway,
			})
		}
		opts.IPAM = &network.IPAM{
			Driver: req.IPAM.Driver,
			Config: ipamCfg,
		}
	}

	resp, err := s.docker.NetworkCreate(ctx, req.Name, opts)
	if err != nil {
		return "", fmt.Errorf("docker create network: %w", err)
	}
	return resp.ID, nil
}

// NetworkRemove 删除网络
func (s *NetworkService) NetworkRemove(ctx context.Context, networkID string) error {
	if err := s.docker.NetworkRemove(ctx, networkID); err != nil {
		return fmt.Errorf("docker remove network: %w", err)
	}
	return nil
}

// NetworkConnect 连接容器到网络
func (s *NetworkService) NetworkConnect(ctx context.Context, networkID, containerID string) error {
	if err := s.docker.NetworkConnect(ctx, networkID, containerID, nil); err != nil {
		return fmt.Errorf("docker network connect: %w", err)
	}
	return nil
}

// NetworkDisconnect 断开容器与网络连接
func (s *NetworkService) NetworkDisconnect(ctx context.Context, networkID, containerID string, force bool) error {
	if err := s.docker.NetworkDisconnect(ctx, networkID, containerID, force); err != nil {
		return fmt.Errorf("docker network disconnect: %w", err)
	}
	return nil
}
