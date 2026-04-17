package docker

import (
	"context"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
)

// DockerClientInterface 定义 Docker 客户端操作接口，便于测试 mock
type DockerClientInterface interface {
	ContainerList(ctx context.Context, options container.ListOptions) ([]types.Container, error)
	ContainerInspect(ctx context.Context, containerID string) (types.ContainerJSON, error)
	ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, containerName string) (container.CreateResponse, error)
	ContainerStart(ctx context.Context, containerID string, options container.StartOptions) error
	ContainerStop(ctx context.Context, containerID string, options container.StopOptions) error
	ContainerKill(ctx context.Context, containerID string, signal string) error
	ContainerRestart(ctx context.Context, containerID string, options container.StopOptions) error
	ContainerPause(ctx context.Context, containerID string) error
	ContainerUnpause(ctx context.Context, containerID string) error
	ContainerRemove(ctx context.Context, containerID string, options container.RemoveOptions) error
	ContainerRename(ctx context.Context, containerID string, newName string) error
	ContainerLogs(ctx context.Context, containerID string, options container.LogsOptions) (io.ReadCloser, error)
	ContainerExecCreate(ctx context.Context, containerID string, config container.ExecOptions) (container.ExecCreateResponse, error)
	ContainerExecAttach(ctx context.Context, execID string, config container.ExecAttachOptions) (types.HijackedResponse, error)
	ContainerTop(ctx context.Context, containerID string, arguments []string) (container.ContainerTopOKBody, error)
	ContainerExport(ctx context.Context, containerID string) (io.ReadCloser, error)
	ImageList(ctx context.Context, options image.ListOptions) ([]image.Summary, error)
	ImageInspect(ctx context.Context, imageID string) (types.ImageInspect, error)
	ImagePull(ctx context.Context, refStr string, options image.PullOptions) (io.ReadCloser, error)
	ImageRemove(ctx context.Context, imageID string, options image.RemoveOptions) ([]image.DeleteResponse, error)
	ImageBuild(ctx context.Context, buildContext io.Reader, options types.ImageBuildOptions) (types.ImageBuildResponse, error)
	ImageSave(ctx context.Context, imageIDs []string) (io.ReadCloser, error)
	ImageLoad(ctx context.Context, input io.Reader, quiet bool) (image.LoadResponse, error)
	ImageTag(ctx context.Context, source, target string) error
	ImagePush(ctx context.Context, ref string, options image.PushOptions) (io.ReadCloser, error)
	VolumeList(ctx context.Context, filter filters.Args) (volume.ListResponse, error)
	VolumeInspect(ctx context.Context, volumeID string) (volume.Volume, error)
	VolumeCreate(ctx context.Context, options volume.CreateOptions) (volume.Volume, error)
	VolumeRemove(ctx context.Context, volumeID string, force bool) error
	NetworkList(ctx context.Context, options network.ListOptions) ([]network.Inspect, error)
	NetworkInspect(ctx context.Context, networkID string, options network.InspectOptions) (network.Inspect, error)
	NetworkCreate(ctx context.Context, name string, options network.CreateOptions) (network.CreateResponse, error)
	NetworkRemove(ctx context.Context, networkID string) error
	NetworkConnect(ctx context.Context, networkID, containerID string, config *network.EndpointSettings) error
	NetworkDisconnect(ctx context.Context, networkID, containerID string, force bool) error
	Close() error
}

// DockerClient 封装 Docker SDK，提供统一的操作接口
type DockerClient struct {
	cli *client.Client
}

// NewDockerClient 创建并返回一个新的 DockerClient 实例
func NewDockerClient() (*DockerClient, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &DockerClient{cli: cli}, nil
}

// Close 关闭 Docker 客户端连接
func (d *DockerClient) Close() error {
	return d.cli.Close()
}

func (d *DockerClient) Version(ctx context.Context) (types.Version, error) {
	return d.cli.ServerVersion(ctx)
}

// ========== 容器操作 ==========

// ContainerList 列出容器
func (d *DockerClient) ContainerList(ctx context.Context, options container.ListOptions) ([]types.Container, error) {
	return d.cli.ContainerList(ctx, options)
}

// ContainerInspect 获取容器详情
func (d *DockerClient) ContainerInspect(ctx context.Context, containerID string) (types.ContainerJSON, error) {
	return d.cli.ContainerInspect(ctx, containerID)
}

// ContainerCreate 创建容器
func (d *DockerClient) ContainerCreate(
	ctx context.Context,
	config *container.Config,
	hostConfig *container.HostConfig,
	networkingConfig *network.NetworkingConfig,
	containerName string,
) (container.CreateResponse, error) {
	return d.cli.ContainerCreate(ctx, config, hostConfig, networkingConfig, nil, containerName)
}

// ContainerStart 启动容器
func (d *DockerClient) ContainerStart(ctx context.Context, containerID string, options container.StartOptions) error {
	return d.cli.ContainerStart(ctx, containerID, options)
}

// ContainerStop 停止容器
func (d *DockerClient) ContainerStop(ctx context.Context, containerID string, options container.StopOptions) error {
	return d.cli.ContainerStop(ctx, containerID, options)
}

// ContainerKill 强制停止容器
func (d *DockerClient) ContainerKill(ctx context.Context, containerID string, signal string) error {
	return d.cli.ContainerKill(ctx, containerID, signal)
}

// ContainerRestart 重启容器
func (d *DockerClient) ContainerRestart(ctx context.Context, containerID string, options container.StopOptions) error {
	return d.cli.ContainerRestart(ctx, containerID, options)
}

// ContainerPause 暂停容器
func (d *DockerClient) ContainerPause(ctx context.Context, containerID string) error {
	return d.cli.ContainerPause(ctx, containerID)
}

// ContainerUnpause 恢复容器
func (d *DockerClient) ContainerUnpause(ctx context.Context, containerID string) error {
	return d.cli.ContainerUnpause(ctx, containerID)
}

// ContainerRemove 删除容器
func (d *DockerClient) ContainerRemove(ctx context.Context, containerID string, options container.RemoveOptions) error {
	return d.cli.ContainerRemove(ctx, containerID, options)
}

// ContainerRename 重命名容器
func (d *DockerClient) ContainerRename(ctx context.Context, containerID string, newName string) error {
	return d.cli.ContainerRename(ctx, containerID, newName)
}

// ContainerLogs 获取容器日志
func (d *DockerClient) ContainerLogs(ctx context.Context, containerID string, options container.LogsOptions) (io.ReadCloser, error) {
	return d.cli.ContainerLogs(ctx, containerID, options)
}

// ContainerExecCreate 创建 exec 会话
func (d *DockerClient) ContainerExecCreate(ctx context.Context, containerID string, config container.ExecOptions) (container.ExecCreateResponse, error) {
	return d.cli.ContainerExecCreate(ctx, containerID, config)
}

// ContainerExecAttach 附加到 exec 会话
func (d *DockerClient) ContainerExecAttach(ctx context.Context, execID string, config container.ExecAttachOptions) (types.HijackedResponse, error) {
	return d.cli.ContainerExecAttach(ctx, execID, config)
}

// ContainerTop 查看容器进程
func (d *DockerClient) ContainerTop(ctx context.Context, containerID string, arguments []string) (container.ContainerTopOKBody, error) {
	return d.cli.ContainerTop(ctx, containerID, arguments)
}

// ContainerExport 导出容器
func (d *DockerClient) ContainerExport(ctx context.Context, containerID string) (io.ReadCloser, error) {
	return d.cli.ContainerExport(ctx, containerID)
}

// ========== 镜像操作 ==========

// ImageList 列出镜像
func (d *DockerClient) ImageList(ctx context.Context, options image.ListOptions) ([]image.Summary, error) {
	return d.cli.ImageList(ctx, options)
}

// ImageInspect 获取镜像详情
func (d *DockerClient) ImageInspect(ctx context.Context, imageID string) (types.ImageInspect, error) {
	return d.cli.ImageInspect(ctx, imageID)
}

// ImagePull 拉取镜像
func (d *DockerClient) ImagePull(ctx context.Context, refStr string, options image.PullOptions) (io.ReadCloser, error) {
	return d.cli.ImagePull(ctx, refStr, options)
}

// ImageRemove 删除镜像
func (d *DockerClient) ImageRemove(ctx context.Context, imageID string, options image.RemoveOptions) ([]image.DeleteResponse, error) {
	return d.cli.ImageRemove(ctx, imageID, options)
}

// ImageBuild 构建镜像
func (d *DockerClient) ImageBuild(ctx context.Context, buildContext io.Reader, options types.ImageBuildOptions) (types.ImageBuildResponse, error) {
	return d.cli.ImageBuild(ctx, buildContext, options)
}

// ImageSave 导出镜像
func (d *DockerClient) ImageSave(ctx context.Context, imageIDs []string) (io.ReadCloser, error) {
	return d.cli.ImageSave(ctx, imageIDs)
}

// ImageLoad 导入镜像
func (d *DockerClient) ImageLoad(ctx context.Context, input io.Reader, quiet bool) (image.LoadResponse, error) {
	return d.cli.ImageLoad(ctx, input, client.ImageLoadWithQuiet(quiet))
}

// ImageTag 给镜像打标签
func (d *DockerClient) ImageTag(ctx context.Context, source, target string) error {
	return d.cli.ImageTag(ctx, source, target)
}

// ImagePush 推送镜像
func (d *DockerClient) ImagePush(ctx context.Context, ref string, options image.PushOptions) (io.ReadCloser, error) {
	return d.cli.ImagePush(ctx, ref, options)
}

// ========== 数据卷操作 ==========

// VolumeList 列出数据卷
func (d *DockerClient) VolumeList(ctx context.Context, filter filters.Args) (volume.ListResponse, error) {
	return d.cli.VolumeList(ctx, volume.ListOptions{Filters: filter})
}

// VolumeInspect 获取数据卷详情
func (d *DockerClient) VolumeInspect(ctx context.Context, volumeID string) (volume.Volume, error) {
	return d.cli.VolumeInspect(ctx, volumeID)
}

// VolumeCreate 创建数据卷
func (d *DockerClient) VolumeCreate(ctx context.Context, options volume.CreateOptions) (volume.Volume, error) {
	return d.cli.VolumeCreate(ctx, options)
}

// VolumeRemove 删除数据卷
func (d *DockerClient) VolumeRemove(ctx context.Context, volumeID string, force bool) error {
	return d.cli.VolumeRemove(ctx, volumeID, force)
}

// ========== 网络操作 ==========

// NetworkList 列出网络
func (d *DockerClient) NetworkList(ctx context.Context, options network.ListOptions) ([]network.Inspect, error) {
	return d.cli.NetworkList(ctx, options)
}

// NetworkInspect 获取网络详情
func (d *DockerClient) NetworkInspect(ctx context.Context, networkID string, options network.InspectOptions) (network.Inspect, error) {
	return d.cli.NetworkInspect(ctx, networkID, options)
}

// NetworkCreate 创建网络
func (d *DockerClient) NetworkCreate(ctx context.Context, name string, options network.CreateOptions) (network.CreateResponse, error) {
	return d.cli.NetworkCreate(ctx, name, options)
}

// NetworkRemove 删除网络
func (d *DockerClient) NetworkRemove(ctx context.Context, networkID string) error {
	return d.cli.NetworkRemove(ctx, networkID)
}

// NetworkConnect 将容器连接到网络
func (d *DockerClient) NetworkConnect(ctx context.Context, networkID, containerID string, config *network.EndpointSettings) error {
	return d.cli.NetworkConnect(ctx, networkID, containerID, config)
}

// NetworkDisconnect 将容器从网络断开
func (d *DockerClient) NetworkDisconnect(ctx context.Context, networkID, containerID string, force bool) error {
	return d.cli.NetworkDisconnect(ctx, networkID, containerID, force)
}
