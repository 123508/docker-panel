package service

import (
	"context"
	mock_docker "docker-panel/docker_cli_wrapper/mock"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/image"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestImageService_ImageList_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_docker.NewMockDockerClient(ctrl)
	service := NewImageService(mockClient)
	ctx := context.Background()
	req := ImageListRequest{All: true}

	mockImages := []image.Summary{
		{ID: "image1", RepoTags: []string{"ubuntu:latest"}, Size: 100},
		{ID: "image2", RepoTags: []string{"nginx:alpine"}, Size: 200},
	}

	mockClient.EXPECT().ImageList(ctx, image.ListOptions{All: true}).Return(mockImages, nil)

	items, err := service.ImageList(ctx, req)

	assert.NoError(t, err)
	assert.Len(t, items, 2)
	assert.Equal(t, "image1", items[0].ID)
	assert.Equal(t, "image2", items[1].ID)
}

func TestImageService_ImageList_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_docker.NewMockDockerClient(ctrl)
	service := NewImageService(mockClient)
	ctx := context.Background()

	mockClient.EXPECT().ImageList(ctx, gomock.Any()).Return(nil, errors.New("docker error"))

	items, err := service.ImageList(ctx, ImageListRequest{})

	assert.Error(t, err)
	assert.Nil(t, items)
	assert.Contains(t, err.Error(), "docker list images")
}

func TestImageService_ImageInspect_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_docker.NewMockDockerClient(ctrl)
	service := NewImageService(mockClient)
	ctx := context.Background()

	mockInfo := types.ImageInspect{
		ID:       "image1",
		RepoTags: []string{"ubuntu:latest"},
		Author:   "test",
	}

	mockClient.EXPECT().ImageInspect(ctx, "image1").Return(mockInfo, nil)

	detail, err := service.ImageInspect(ctx, "image1")

	assert.NoError(t, err)
	assert.NotNil(t, detail)
	assert.Equal(t, "image1", detail.ID)
	assert.Equal(t, "test", detail.Author)
}

func TestImageService_ImageInspect_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_docker.NewMockDockerClient(ctrl)
	service := NewImageService(mockClient)
	ctx := context.Background()

	mockClient.EXPECT().ImageInspect(ctx, "image1").Return(types.ImageInspect{}, errors.New("not found"))

	detail, err := service.ImageInspect(ctx, "image1")

	assert.Error(t, err)
	assert.Nil(t, detail)
	assert.Contains(t, err.Error(), "docker inspect image")
}

func TestImageService_ImagePull_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_docker.NewMockDockerClient(ctrl)
	service := NewImageService(mockClient)
	ctx := context.Background()

	mockRC := io.NopCloser(strings.NewReader("pulling..."))
	mockClient.EXPECT().ImagePull(ctx, "ubuntu:latest", image.PullOptions{}).Return(mockRC, nil)

	rc, err := service.ImagePull(ctx, "ubuntu:latest")

	assert.NoError(t, err)
	assert.NotNil(t, rc)
}

func TestImageService_ImagePull_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_docker.NewMockDockerClient(ctrl)
	service := NewImageService(mockClient)
	ctx := context.Background()

	mockClient.EXPECT().ImagePull(ctx, "ubuntu:latest", image.PullOptions{}).Return(nil, errors.New("pull error"))

	rc, err := service.ImagePull(ctx, "ubuntu:latest")

	assert.Error(t, err)
	assert.Nil(t, rc)
}

func TestImageService_ImageRemove_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_docker.NewMockDockerClient(ctrl)
	service := NewImageService(mockClient)
	ctx := context.Background()

	mockResp := []image.DeleteResponse{{Deleted: "image1"}}
	mockClient.EXPECT().ImageRemove(ctx, "image1", image.RemoveOptions{Force: true, PruneChildren: false}).Return(mockResp, nil)

	resp, err := service.ImageRemove(ctx, "image1", true, false)

	assert.NoError(t, err)
	assert.Equal(t, mockResp, resp)
}

func TestImageService_ImageRemove_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_docker.NewMockDockerClient(ctrl)
	service := NewImageService(mockClient)
	ctx := context.Background()

	mockClient.EXPECT().ImageRemove(ctx, "image1", gomock.Any()).Return(nil, errors.New("remove error"))

	resp, err := service.ImageRemove(ctx, "image1", true, false)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestImageService_ImageBuild_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_docker.NewMockDockerClient(ctrl)
	service := NewImageService(mockClient)
	ctx := context.Background()

	mockCtx := strings.NewReader("dockerfile")
	mockOpts := types.ImageBuildOptions{Tags: []string{"myimage"}}
	mockResp := types.ImageBuildResponse{Body: io.NopCloser(strings.NewReader("building..."))}

	mockClient.EXPECT().ImageBuild(ctx, mockCtx, mockOpts).Return(mockResp, nil)

	resp, err := service.ImageBuild(ctx, mockCtx, mockOpts)

	assert.NoError(t, err)
	assert.NotNil(t, resp.Body)
}

func TestImageService_ImageBuild_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_docker.NewMockDockerClient(ctrl)
	service := NewImageService(mockClient)
	ctx := context.Background()

	mockCtx := strings.NewReader("dockerfile")
	mockOpts := types.ImageBuildOptions{}

	mockClient.EXPECT().ImageBuild(ctx, mockCtx, mockOpts).Return(types.ImageBuildResponse{}, errors.New("build error"))

	resp, err := service.ImageBuild(ctx, mockCtx, mockOpts)

	assert.Error(t, err)
	assert.Nil(t, resp.Body)
}

func TestImageService_ImageSave_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_docker.NewMockDockerClient(ctrl)
	service := NewImageService(mockClient)
	ctx := context.Background()

	mockRC := io.NopCloser(strings.NewReader("saving..."))
	mockClient.EXPECT().ImageSave(ctx, []string{"image1"}).Return(mockRC, nil)

	rc, err := service.ImageSave(ctx, []string{"image1"})

	assert.NoError(t, err)
	assert.NotNil(t, rc)
}

func TestImageService_ImageSave_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_docker.NewMockDockerClient(ctrl)
	service := NewImageService(mockClient)
	ctx := context.Background()

	mockClient.EXPECT().ImageSave(ctx, []string{"image1"}).Return(nil, errors.New("save error"))

	rc, err := service.ImageSave(ctx, []string{"image1"})

	assert.Error(t, err)
	assert.Nil(t, rc)
}

func TestImageService_ImageLoad_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_docker.NewMockDockerClient(ctrl)
	service := NewImageService(mockClient)
	ctx := context.Background()

	mockInput := strings.NewReader("load...")
	mockResp := image.LoadResponse{Body: io.NopCloser(strings.NewReader("loaded"))}

	mockClient.EXPECT().ImageLoad(ctx, mockInput, false).Return(mockResp, nil)

	resp, err := service.ImageLoad(ctx, mockInput, false)

	assert.NoError(t, err)
	assert.NotNil(t, resp.Body)
}

func TestImageService_ImageLoad_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_docker.NewMockDockerClient(ctrl)
	service := NewImageService(mockClient)
	ctx := context.Background()

	mockInput := strings.NewReader("load...")

	mockClient.EXPECT().ImageLoad(ctx, mockInput, false).Return(image.LoadResponse{}, errors.New("load error"))

	resp, err := service.ImageLoad(ctx, mockInput, false)

	assert.Error(t, err)
	assert.Nil(t, resp.Body)
}

func TestImageService_ImageTag_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_docker.NewMockDockerClient(ctrl)
	service := NewImageService(mockClient)
	ctx := context.Background()

	mockClient.EXPECT().ImageTag(ctx, "source", "target").Return(nil)

	err := service.ImageTag(ctx, "source", "target")

	assert.NoError(t, err)
}

func TestImageService_ImageTag_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_docker.NewMockDockerClient(ctrl)
	service := NewImageService(mockClient)
	ctx := context.Background()

	mockClient.EXPECT().ImageTag(ctx, "source", "target").Return(errors.New("tag error"))

	err := service.ImageTag(ctx, "source", "target")

	assert.Error(t, err)
}

func TestImageService_ImagePush_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_docker.NewMockDockerClient(ctrl)
	service := NewImageService(mockClient)
	ctx := context.Background()

	mockRC := io.NopCloser(strings.NewReader("pushing..."))
	mockClient.EXPECT().ImagePush(ctx, "ubuntu:latest", image.PushOptions{}).Return(mockRC, nil)

	rc, err := service.ImagePush(ctx, "ubuntu:latest")

	assert.NoError(t, err)
	assert.NotNil(t, rc)
}

func TestImageService_ImagePush_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_docker.NewMockDockerClient(ctrl)
	service := NewImageService(mockClient)
	ctx := context.Background()

	mockClient.EXPECT().ImagePush(ctx, "ubuntu:latest", image.PushOptions{}).Return(nil, errors.New("push error"))

	rc, err := service.ImagePush(ctx, "ubuntu:latest")

	assert.Error(t, err)
	assert.Nil(t, rc)
}
