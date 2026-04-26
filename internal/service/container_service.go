package service

import (
	"context"
	"docker-panel/internal/docker"
	"fmt"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/go-connections/nat"
)

// ========== 数据结构定义 ==========

// ContainerListItem 容器列表项
type ContainerListItem struct {
	ID      string            `json:"id"`
	Names   []string          `json:"names"`
	Image   string            `json:"image"`
	ImageID string            `json:"image_id"`
	Command string            `json:"command"`
	Created int64             `json:"created"`
	State   string            `json:"state"`
	Status  string            `json:"status"`
	Ports   []PortBinding     `json:"ports"`
	Labels  map[string]string `json:"labels"`
	Mounts  []MountPoint      `json:"mounts"`
}

// PortBinding 端口绑定
type PortBinding struct {
	IP          string `json:"ip"`
	PrivatePort uint16 `json:"private_port"`
	PublicPort  uint16 `json:"public_port"`
	Type        string `json:"type"`
}

// MountPoint 挂载点
type MountPoint struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Driver      string `json:"driver"`
	Mode        string `json:"mode"`
	RW          bool   `json:"rw"`
}

// ContainerState 容器状态
type ContainerState struct {
	Status     string `json:"status"`
	Running    bool   `json:"running"`
	Paused     bool   `json:"paused"`
	Restarting bool   `json:"restarting"`
	OOMKilled  bool   `json:"oom_killed"`
	Dead       bool   `json:"dead"`
	Pid        int    `json:"pid"`
	ExitCode   int    `json:"exit_code"`
	Error      string `json:"error"`
	StartedAt  string `json:"started_at"`
	FinishedAt string `json:"finished_at"`
}

// ContainerConfig 容器配置
type ContainerConfig struct {
	Hostname     string              `json:"hostname"`
	Domainname   string              `json:"domainname"`
	User         string              `json:"user"`
	AttachStdin  bool                `json:"attach_stdin"`
	AttachStdout bool                `json:"attach_stdout"`
	AttachStderr bool                `json:"attach_stderr"`
	Tty          bool                `json:"tty"`
	OpenStdin    bool                `json:"open_stdin"`
	StdinOnce    bool                `json:"stdin_once"`
	Env          []string            `json:"env"`
	Cmd          []string            `json:"cmd"`
	Image        string              `json:"image"`
	Volumes      map[string]struct{} `json:"volumes"`
	WorkingDir   string              `json:"working_dir"`
	Entrypoint   []string            `json:"entrypoint"`
	Labels       map[string]string   `json:"labels"`
	ExposedPorts map[string]struct{} `json:"exposed_ports"`
}

// NetworkSettings 网络设置
type NetworkSettings struct {
	Bridge               string                      `json:"bridge"`
	SandboxID            string                      `json:"sandbox_id"`
	HairpinMode          bool                        `json:"hairpin_mode"`
	LinkLocalIPv6Address string                      `json:"link_local_ipv6_address"`
	Ports                map[string][]PortBinding    `json:"ports"`
	SandboxKey           string                      `json:"sandbox_key"`
	IPAddress            string                      `json:"ip_address"`
	MacAddress           string                      `json:"mac_address"`
	Networks             map[string]EndpointSettings `json:"networks"`
}

// EndpointSettings 端点设置
type EndpointSettings struct {
	NetworkID  string `json:"network_id"`
	EndpointID string `json:"endpoint_id"`
	Gateway    string `json:"gateway"`
	IPAddress  string `json:"ip_address"`
	MacAddress string `json:"mac_address"`
}

// RestartPolicy 重启策略
type RestartPolicy struct {
	Name              string `json:"name"`
	MaximumRetryCount int    `json:"maximum_retry_count"`
}

// HostConfig 主机配置
type HostConfig struct {
	Binds           []string                 `json:"binds"`
	NetworkMode     string                   `json:"network_mode"`
	PortBindings    map[string][]PortBinding `json:"port_bindings"`
	RestartPolicy   RestartPolicy            `json:"restart_policy"`
	AutoRemove      bool                     `json:"auto_remove"`
	VolumeDriver    string                   `json:"volume_driver"`
	VolumesFrom     []string                 `json:"volumes_from"`
	CapAdd          []string                 `json:"cap_add"`
	CapDrop         []string                 `json:"cap_drop"`
	Dns             []string                 `json:"dns"`
	ExtraHosts      []string                 `json:"extra_hosts"`
	Privileged      bool                     `json:"privileged"`
	PublishAllPorts bool                     `json:"publish_all_ports"`
	ReadonlyRootfs  bool                     `json:"readonly_rootfs"`
	Memory          int64                    `json:"memory"`
	MemorySwap      int64                    `json:"memory_swap"`
	NanoCPUs        int64                    `json:"nano_cpus"`
	CPUShares       int64                    `json:"cpu_shares"`
}

// ContainerDetail 容器详情
type ContainerDetail struct {
	ID              string          `json:"id"`
	Created         string          `json:"created"`
	Path            string          `json:"path"`
	Args            []string        `json:"args"`
	State           ContainerState  `json:"state"`
	Image           string          `json:"image"`
	Name            string          `json:"name"`
	RestartCount    int             `json:"restart_count"`
	Driver          string          `json:"driver"`
	Platform        string          `json:"platform"`
	Mounts          []MountPoint    `json:"mounts"`
	Config          ContainerConfig `json:"config"`
	NetworkSettings NetworkSettings `json:"network_settings"`
	HostConfig      HostConfig      `json:"host_config"`
}

// ContainerTopResponse 容器进程信息
type ContainerTopResponse struct {
	Titles    []string   `json:"titles"`
	Processes [][]string `json:"processes"`
}

// PortBindingRequest 端口绑定请求
type PortBindingRequest struct {
	HostIP   string `json:"host_ip"`
	HostPort string `json:"host_port"`
}

// BindOptions 绑定挂载选项
type BindOptions struct {
	Propagation string `json:"propagation"`
}

// VolumeOptions 卷挂载选项
type VolumeOptions struct {
	NoCopy bool              `json:"no_copy"`
	Labels map[string]string `json:"labels"`
}

// MountRequest 挂载请求
type MountRequest struct {
	Type          string         `json:"type"`
	Source        string         `json:"source"`
	Target        string         `json:"target"`
	ReadOnly      bool           `json:"read_only"`
	BindOptions   *BindOptions   `json:"bind_options,omitempty"`
	VolumeOptions *VolumeOptions `json:"volume_options,omitempty"`
}

// HostConfigRequest 主机配置请求
type HostConfigRequest struct {
	PortBindings  map[string][]PortBindingRequest `json:"port_bindings"`
	Binds         []string                        `json:"binds"`
	Mounts        []MountRequest                  `json:"mounts"`
	NetworkMode   string                          `json:"network_mode"`
	RestartPolicy RestartPolicy                   `json:"restart_policy"`
	Memory        int64                           `json:"memory"`
	NanoCPUs      int64                           `json:"nano_cpus"`
	Privileged    bool                            `json:"privileged"`
}

// NetworkingConfigRequest 网络配置请求
type NetworkingConfigRequest struct {
	EndpointsConfig map[string]*network.EndpointSettings `json:"endpoints_config"`
}

// ContainerCreateRequest 容器创建请求
type ContainerCreateRequest struct {
	Name             string                   `json:"name" binding:"required"`
	Image            string                   `json:"image" binding:"required"`
	Cmd              []string                 `json:"cmd"`
	Entrypoint       []string                 `json:"entrypoint"`
	Env              []string                 `json:"env"`
	WorkingDir       string                   `json:"working_dir"`
	ExposedPorts     map[string]struct{}      `json:"exposed_ports"`
	HostConfig       HostConfigRequest        `json:"host_config"`
	NetworkingConfig *NetworkingConfigRequest `json:"networking_config"`
	Labels           map[string]string        `json:"labels"`
}

// ContainerListRequest 容器列表请求
type ContainerListRequest struct {
	All     bool   `form:"all"`
	Limit   int    `form:"limit"`
	Filters string `form:"filters"`
}

// ContainerStopRequest 停止容器请求
type ContainerStopRequest struct {
	Timeout int `json:"timeout"`
}

// ContainerKillRequest 强制停止容器请求
type ContainerKillRequest struct {
	Signal string `json:"signal"`
}

// ContainerRestartRequest 重启容器请求
type ContainerRestartRequest struct {
	Timeout int `json:"timeout"`
}

// ContainerRemoveRequest 删除容器请求
type ContainerRemoveRequest struct {
	Force         bool `form:"force"`
	RemoveVolumes bool `form:"remove_volumes"`
}

// ContainerRenameRequest 重命名容器请求
type ContainerRenameRequest struct {
	NewName string `json:"new_name" binding:"required"`
}

// ContainerLogsOptions 容器日志选项
type ContainerLogsOptions struct {
	ShowStdout bool   `form:"show_stdout"`
	ShowStderr bool   `form:"show_stderr"`
	Since      string `form:"since"`
	Until      string `form:"until"`
	Timestamps bool   `form:"timestamps"`
	Follow     bool   `form:"follow"`
	Tail       string `form:"tail"`
}

// ContainerExecRequest 容器执行命令请求
type ContainerExecRequest struct {
	AttachStdin  bool     `json:"attach_stdin"`
	AttachStdout bool     `json:"attach_stdout"`
	AttachStderr bool     `json:"attach_stderr"`
	Tty          bool     `json:"tty"`
	Cmd          []string `json:"cmd" binding:"required"`
	Env          []string `json:"env"`
	WorkingDir   string   `json:"working_dir"`
	User         string   `json:"user"`
}

// ContainerTopRequest 容器进程请求
type ContainerTopRequest struct {
	PsArgs string `form:"ps_args"`
}

// ========== ContainerService ==========

// ContainerService 容器业务逻辑
type ContainerService struct {
	docker docker.DockerClientInterface
}

// NewContainerService 创建 ContainerService
func NewContainerService(docker docker.DockerClientInterface) *ContainerService {
	return &ContainerService{docker: docker}
}

// ContainerList 获取容器列表
func (s *ContainerService) ContainerList(ctx context.Context, req ContainerListRequest) ([]ContainerListItem, error) {
	opts := container.ListOptions{
		All:   req.All,
		Limit: req.Limit,
	}

	containers, err := s.docker.ContainerList(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("docker list containers: %w", err)
	}

	items := make([]ContainerListItem, 0, len(containers))
	for _, c := range containers {
		item := ContainerListItem{
			ID:      c.ID,
			Names:   c.Names,
			Image:   c.Image,
			ImageID: c.ImageID,
			Command: c.Command,
			Created: c.Created,
			State:   c.State,
			Status:  c.Status,
			Labels:  c.Labels,
		}
		for _, p := range c.Ports {
			item.Ports = append(item.Ports, PortBinding{
				IP:          p.IP,
				PrivatePort: p.PrivatePort,
				PublicPort:  p.PublicPort,
				Type:        p.Type,
			})
		}
		for _, m := range c.Mounts {
			item.Mounts = append(item.Mounts, MountPoint{
				Type:        string(m.Type),
				Name:        m.Name,
				Source:      m.Source,
				Destination: m.Destination,
				Driver:      m.Driver,
				Mode:        m.Mode,
				RW:          m.RW,
			})
		}
		items = append(items, item)
	}
	return items, nil
}

// ContainerInspect 获取容器详情
func (s *ContainerService) ContainerInspect(ctx context.Context, containerID string) (*ContainerDetail, error) {
	info, err := s.docker.ContainerInspect(ctx, containerID)
	if err != nil {
		return nil, fmt.Errorf("docker inspect container: %w", err)
	}

	detail := &ContainerDetail{
		ID:           info.ID,
		Created:      info.Created,
		Path:         info.Path,
		Args:         info.Args,
		Image:        info.Image,
		Name:         info.Name,
		RestartCount: info.RestartCount,
		Driver:       info.Driver,
		Platform:     info.Platform,
	}

	if info.State != nil {
		detail.State = ContainerState{
			Status:     info.State.Status,
			Running:    info.State.Running,
			Paused:     info.State.Paused,
			Restarting: info.State.Restarting,
			OOMKilled:  info.State.OOMKilled,
			Dead:       info.State.Dead,
			Pid:        info.State.Pid,
			ExitCode:   info.State.ExitCode,
			Error:      info.State.Error,
			StartedAt:  info.State.StartedAt,
			FinishedAt: info.State.FinishedAt,
		}
	}

	if info.Config != nil {
		detail.Config = ContainerConfig{
			Hostname:   info.Config.Hostname,
			Domainname: info.Config.Domainname,
			User:       info.Config.User,
			Tty:        info.Config.Tty,
			OpenStdin:  info.Config.OpenStdin,
			Env:        info.Config.Env,
			Image:      info.Config.Image,
			WorkingDir: info.Config.WorkingDir,
			Labels:     info.Config.Labels,
		}
		if info.Config.Cmd != nil {
			detail.Config.Cmd = []string(info.Config.Cmd)
		}
		if info.Config.Entrypoint != nil {
			detail.Config.Entrypoint = []string(info.Config.Entrypoint)
		}
	}

	for _, m := range info.Mounts {
		detail.Mounts = append(detail.Mounts, MountPoint{
			Type:        string(m.Type),
			Name:        m.Name,
			Source:      m.Source,
			Destination: m.Destination,
			Driver:      m.Driver,
			Mode:        m.Mode,
			RW:          m.RW,
		})
	}

	return detail, nil
}

// ContainerCreate 创建容器
func (s *ContainerService) ContainerCreate(ctx context.Context, req ContainerCreateRequest) (string, error) {
	// 构建容器配置
	cfg := &container.Config{
		Image:      req.Image,
		Env:        req.Env,
		WorkingDir: req.WorkingDir,
		Labels:     req.Labels,
	}
	if len(req.Cmd) > 0 {
		cfg.Cmd = strslice.StrSlice(req.Cmd)
	}
	if len(req.Entrypoint) > 0 {
		cfg.Entrypoint = strslice.StrSlice(req.Entrypoint)
	}

	// 构建端口暴露
	if len(req.HostConfig.PortBindings) > 0 {
		cfg.ExposedPorts = make(nat.PortSet)
		for portProto := range req.HostConfig.PortBindings {
			cfg.ExposedPorts[nat.Port(portProto)] = struct{}{}
		}
	}

	// 构建主机配置
	hostCfg := &container.HostConfig{
		Binds:       req.HostConfig.Binds,
		NetworkMode: container.NetworkMode(req.HostConfig.NetworkMode),
		Privileged:  req.HostConfig.Privileged,
		Resources: container.Resources{
			Memory:   req.HostConfig.Memory,
			NanoCPUs: req.HostConfig.NanoCPUs,
		},
	}

	// 重启策略
	hostCfg.RestartPolicy = container.RestartPolicy{
		Name:              container.RestartPolicyMode(req.HostConfig.RestartPolicy.Name),
		MaximumRetryCount: req.HostConfig.RestartPolicy.MaximumRetryCount,
	}

	// 端口绑定
	if len(req.HostConfig.PortBindings) > 0 {
		hostCfg.PortBindings = make(nat.PortMap)
		for portProto, bindings := range req.HostConfig.PortBindings {
			var natBindings []nat.PortBinding
			for _, b := range bindings {
				natBindings = append(natBindings, nat.PortBinding{
					HostIP:   b.HostIP,
					HostPort: b.HostPort,
				})
			}
			hostCfg.PortBindings[nat.Port(portProto)] = natBindings
		}
	}

	// 挂载
	for _, m := range req.HostConfig.Mounts {
		mnt := mount.Mount{
			Type:     mount.Type(m.Type),
			Source:   m.Source,
			Target:   m.Target,
			ReadOnly: m.ReadOnly,
		}
		if m.BindOptions != nil {
			mnt.BindOptions = &mount.BindOptions{
				Propagation: mount.Propagation(m.BindOptions.Propagation),
			}
		}
		if m.VolumeOptions != nil {
			mnt.VolumeOptions = &mount.VolumeOptions{
				NoCopy: m.VolumeOptions.NoCopy,
				Labels: m.VolumeOptions.Labels,
			}
		}
		hostCfg.Mounts = append(hostCfg.Mounts, mnt)
	}

	// 网络配置
	var netCfg *network.NetworkingConfig
	if req.NetworkingConfig != nil {
		netCfg = &network.NetworkingConfig{
			EndpointsConfig: req.NetworkingConfig.EndpointsConfig,
		}
	}

	resp, err := s.docker.ContainerCreate(ctx, cfg, hostCfg, netCfg, req.Name)
	if err != nil {
		return "", fmt.Errorf("docker create container: %w", err)
	}
	return resp.ID, nil
}

// ContainerStart 启动容器
func (s *ContainerService) ContainerStart(ctx context.Context, containerID string) error {
	return s.docker.ContainerStart(ctx, containerID, container.StartOptions{})
}

// ContainerStop 停止容器
func (s *ContainerService) ContainerStop(ctx context.Context, containerID string, timeout int) error {
	opts := container.StopOptions{}
	if timeout > 0 {
		t := timeout
		opts.Timeout = &t
	}
	return s.docker.ContainerStop(ctx, containerID, opts)
}

// ContainerKill 强制停止容器
func (s *ContainerService) ContainerKill(ctx context.Context, containerID string, signal string) error {
	if signal == "" {
		signal = "SIGKILL"
	}
	return s.docker.ContainerKill(ctx, containerID, signal)
}

// ContainerRestart 重启容器
func (s *ContainerService) ContainerRestart(ctx context.Context, containerID string, timeout int) error {
	opts := container.StopOptions{}
	if timeout > 0 {
		t := timeout
		opts.Timeout = &t
	}
	return s.docker.ContainerRestart(ctx, containerID, opts)
}

// ContainerPause 暂停容器
func (s *ContainerService) ContainerPause(ctx context.Context, containerID string) error {
	return s.docker.ContainerPause(ctx, containerID)
}

// ContainerUnpause 恢复容器
func (s *ContainerService) ContainerUnpause(ctx context.Context, containerID string) error {
	return s.docker.ContainerUnpause(ctx, containerID)
}

// ContainerRemove 删除容器
func (s *ContainerService) ContainerRemove(ctx context.Context, containerID string, force bool, removeVolumes bool) error {
	return s.docker.ContainerRemove(ctx, containerID, container.RemoveOptions{
		Force:         force,
		RemoveVolumes: removeVolumes,
	})
}

// ContainerRename 重命名容器
func (s *ContainerService) ContainerRename(ctx context.Context, containerID string, newName string) error {
	return s.docker.ContainerRename(ctx, containerID, newName)
}

// ContainerLogs 获取容器日志
func (s *ContainerService) ContainerLogs(ctx context.Context, containerID string, opts ContainerLogsOptions) (io.ReadCloser, error) {
	tail := opts.Tail
	if tail == "" {
		tail = "100"
	}
	return s.docker.ContainerLogs(ctx, containerID, container.LogsOptions{
		ShowStdout: opts.ShowStdout,
		ShowStderr: opts.ShowStderr,
		Since:      opts.Since,
		Until:      opts.Until,
		Timestamps: opts.Timestamps,
		Follow:     opts.Follow,
		Tail:       tail,
	})
}

// ContainerExec 在容器中执行命令，返回 hijacked 连接
func (s *ContainerService) ContainerExec(ctx context.Context, containerID string, req ContainerExecRequest) (types.HijackedResponse, error) {
	execID, err := s.docker.ContainerExecCreate(ctx, containerID, container.ExecOptions{
		AttachStdin:  req.AttachStdin,
		AttachStdout: req.AttachStdout,
		AttachStderr: req.AttachStderr,
		Tty:          req.Tty,
		Cmd:          req.Cmd,
		Env:          req.Env,
		WorkingDir:   req.WorkingDir,
		User:         req.User,
	})
	if err != nil {
		return types.HijackedResponse{}, fmt.Errorf("docker exec create: %w", err)
	}

	return s.docker.ContainerExecAttach(ctx, execID.ID, container.ExecAttachOptions{
		Tty: req.Tty,
	})
}

// ContainerExecCreate 创建 exec 会话（用于 terminal）
func (s *ContainerService) ContainerExecCreate(ctx context.Context, containerID string, req ContainerExecRequest) (string, error) {
	resp, err := s.docker.ContainerExecCreate(ctx, containerID, container.ExecOptions{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
		Cmd:          req.Cmd,
	})
	if err != nil {
		return "", fmt.Errorf("docker exec create: %w", err)
	}
	return resp.ID, nil
}

// ContainerExecAttach 附加到 exec 会话（用于 terminal）
func (s *ContainerService) ContainerExecAttach(ctx context.Context, execID string) (types.HijackedResponse, error) {
	return s.docker.ContainerExecAttach(ctx, execID, container.ExecAttachOptions{Tty: true})
}

// ContainerTop 查看容器进程
func (s *ContainerService) ContainerTop(ctx context.Context, containerID string, psArgs string) (*ContainerTopResponse, error) {
	args := []string{}
	if psArgs != "" {
		args = []string{psArgs}
	}
	result, err := s.docker.ContainerTop(ctx, containerID, args)
	if err != nil {
		return nil, fmt.Errorf("docker container top: %w", err)
	}
	return &ContainerTopResponse{
		Titles:    result.Titles,
		Processes: result.Processes,
	}, nil
}

// ContainerExport 导出容器
func (s *ContainerService) ContainerExport(ctx context.Context, containerID string) (io.ReadCloser, error) {
	return s.docker.ContainerExport(ctx, containerID)
}

// GetContainersByIDs 根据容器 ID 列表查询容器简要信息。
// 先获取全量容器列表，再按传入的 ID 顺序筛选，保留请求的先后次序。
func (s *ContainerService) GetContainersByIDs(ctx context.Context, ids []string) ([]ContainerListItem, error) {
	// 获取全部容器
	all, err := s.ContainerList(ctx, ContainerListRequest{All: true})
	if err != nil {
		return nil, err
	}

	// 构建 ID 到容器信息的映射
	ordered := make(map[string]ContainerListItem, len(all))
	for _, item := range all {
		ordered[item.ID] = item
	}

	// 按传入的 ID 顺序返回结果，保持 LRU 中的先后次序
	result := make([]ContainerListItem, 0, len(ids))
	for _, id := range ids {
		if item, ok := ordered[id]; ok {
			result = append(result, item)
		}
	}
	return result, nil
}
