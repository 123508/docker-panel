package service

import (
	"context"
	"docker-panel/docker_cli_wrapper"
	"fmt"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/image"
)

// ========== 数据结构定义 ==========

// ImageListItem 镜像列表项
type ImageListItem struct {
	ID          string            `json:"id"`
	ParentID    string            `json:"parent_id"`
	RepoTags    []string          `json:"repo_tags"`
	RepoDigests []string          `json:"repo_digests"`
	Created     int64             `json:"created"`
	Size        int64             `json:"size"`
	VirtualSize int64             `json:"virtual_size"`
	SharedSize  int64             `json:"shared_size"`
	Labels      map[string]string `json:"labels"`
	Containers  int               `json:"containers"`
}

// ImageDetail 镜像详情
type ImageDetail struct {
	ID              string          `json:"id"`
	RepoTags        []string        `json:"repo_tags"`
	RepoDigests     []string        `json:"repo_digests"`
	Parent          string          `json:"parent"`
	Comment         string          `json:"comment"`
	Created         string          `json:"created"`
	Container       string          `json:"container"`
	ContainerConfig ContainerConfig `json:"container_config"`
	DockerVersion   string          `json:"docker_version"`
	Author          string          `json:"author"`
	Config          ContainerConfig `json:"config"`
	Architecture    string          `json:"architecture"`
	Os              string          `json:"os"`
	Size            int64           `json:"size"`
	VirtualSize     int64           `json:"virtual_size"`
}

// ImageListRequest 镜像列表请求
type ImageListRequest struct {
	All     bool   `form:"all"`
	Filters string `form:"filters"`
}

// ImagePullRequest 拉取镜像请求
type ImagePullRequest struct {
	Image string `json:"image" binding:"required"`
}

// ImageRemoveRequest 删除镜像请求
type ImageRemoveRequest struct {
	Force         bool `form:"force"`
	PruneChildren bool `form:"prune_children"`
}

// ImageTagRequest 镜像打标签请求
type ImageTagRequest struct {
	Source string `json:"source" binding:"required"`
	Target string `json:"target" binding:"required"`
}

// ImagePushRequest 推送镜像请求
type ImagePushRequest struct {
	Image string `json:"image" binding:"required"`
}

// ImageImportRequest 导入镜像请求
type ImageImportRequest struct {
	Source     string `json:"source"`
	Repository string `json:"repository"`
	Tag        string `json:"tag"`
	Message    string `json:"message"`
}

// ========== ImageService ==========

// ImageService 镜像业务逻辑
type ImageService struct {
	docker docker_cli_wrapper.DockerClientInterface
}

// NewImageService 创建 ImageService
func NewImageService(docker docker_cli_wrapper.DockerClientInterface) *ImageService {
	return &ImageService{docker: docker}
}

// ImageList 获取镜像列表
func (s *ImageService) ImageList(ctx context.Context, req ImageListRequest) ([]ImageListItem, error) {
	images, err := s.docker.ImageList(ctx, image.ListOptions{All: req.All})
	if err != nil {
		return nil, fmt.Errorf("docker list images: %w", err)
	}

	items := make([]ImageListItem, 0, len(images))
	for _, img := range images {
		items = append(items, ImageListItem{
			ID:          img.ID,
			ParentID:    img.ParentID,
			RepoTags:    img.RepoTags,
			RepoDigests: img.RepoDigests,
			Created:     img.Created,
			Size:        img.Size,
			VirtualSize: img.VirtualSize,
			SharedSize:  img.SharedSize,
			Labels:      img.Labels,
			Containers:  int(img.Containers),
		})
	}
	return items, nil
}

// ImageInspect 获取镜像详情
func (s *ImageService) ImageInspect(ctx context.Context, imageID string) (*ImageDetail, error) {
	info, err := s.docker.ImageInspect(ctx, imageID)
	if err != nil {
		return nil, fmt.Errorf("docker inspect image: %w", err)
	}

	detail := &ImageDetail{
		ID:            info.ID,
		RepoTags:      info.RepoTags,
		RepoDigests:   info.RepoDigests,
		Parent:        info.Parent,
		Comment:       info.Comment,
		Created:       info.Created,
		DockerVersion: info.DockerVersion,
		Author:        info.Author,
		Architecture:  info.Architecture,
		Os:            info.Os,
		Size:          info.Size,
		VirtualSize:   info.VirtualSize,
	}
	return detail, nil
}

// ImagePull 拉取镜像，返回流式响应
func (s *ImageService) ImagePull(ctx context.Context, imageRef string) (io.ReadCloser, error) {
	rc, err := s.docker.ImagePull(ctx, imageRef, image.PullOptions{})
	if err != nil {
		return nil, fmt.Errorf("docker pull image: %w", err)
	}
	return rc, nil
}

// ImageRemove 删除镜像
func (s *ImageService) ImageRemove(ctx context.Context, imageID string, force bool, pruneChildren bool) ([]image.DeleteResponse, error) {
	result, err := s.docker.ImageRemove(ctx, imageID, image.RemoveOptions{
		Force:         force,
		PruneChildren: pruneChildren,
	})
	if err != nil {
		return nil, fmt.Errorf("docker remove image: %w", err)
	}
	return result, nil
}

// ImageBuild 构建镜像
func (s *ImageService) ImageBuild(ctx context.Context, buildContext io.Reader, options types.ImageBuildOptions) (types.ImageBuildResponse, error) {
	resp, err := s.docker.ImageBuild(ctx, buildContext, options)
	if err != nil {
		return types.ImageBuildResponse{}, fmt.Errorf("docker build image: %w", err)
	}
	return resp, nil
}

// ImageSave 导出镜像
func (s *ImageService) ImageSave(ctx context.Context, imageIDs []string) (io.ReadCloser, error) {
	rc, err := s.docker.ImageSave(ctx, imageIDs)
	if err != nil {
		return nil, fmt.Errorf("docker save image: %w", err)
	}
	return rc, nil
}

// ImageLoad 导入镜像
func (s *ImageService) ImageLoad(ctx context.Context, input io.Reader, quiet bool) (image.LoadResponse, error) {
	resp, err := s.docker.ImageLoad(ctx, input, quiet)
	if err != nil {
		return image.LoadResponse{}, fmt.Errorf("docker load image: %w", err)
	}
	return resp, nil
}

// ImageTag 给镜像打标签
func (s *ImageService) ImageTag(ctx context.Context, source, target string) error {
	if err := s.docker.ImageTag(ctx, source, target); err != nil {
		return fmt.Errorf("docker tag image: %w", err)
	}
	return nil
}

// ImagePush 推送镜像
func (s *ImageService) ImagePush(ctx context.Context, imageRef string) (io.ReadCloser, error) {
	rc, err := s.docker.ImagePush(ctx, imageRef, image.PushOptions{})
	if err != nil {
		return nil, fmt.Errorf("docker push image: %w", err)
	}
	return rc, nil
}
