package service

import (
	"context"
	mock_docker "docker-panel/docker_cli_wrapper/mock"
	"errors"
	"testing"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/volume"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestVolumeService_VolumeList_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_docker.NewMockDockerClient(ctrl)
	service := NewVolumeService(mockClient)
	ctx := context.Background()

	mockResp := volume.ListResponse{
		Volumes: []*volume.Volume{
			{Name: "vol1", Driver: "local"},
			{Name: "vol2", Driver: "local"},
		},
	}

	mockClient.EXPECT().VolumeList(ctx, filters.NewArgs()).Return(mockResp, nil)

	items, err := service.VolumeList(ctx)

	assert.NoError(t, err)
	assert.Len(t, items, 2)
	assert.Equal(t, "vol1", items[0].Name)
	assert.Equal(t, "vol2", items[1].Name)
}

func TestVolumeService_VolumeList_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_docker.NewMockDockerClient(ctrl)
	service := NewVolumeService(mockClient)
	ctx := context.Background()

	mockClient.EXPECT().VolumeList(ctx, filters.NewArgs()).Return(volume.ListResponse{}, errors.New("docker error"))

	items, err := service.VolumeList(ctx)

	assert.Error(t, err)
	assert.Nil(t, items)
	assert.Contains(t, err.Error(), "docker list volumes")
}

func TestVolumeService_VolumeInspect_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_docker.NewMockDockerClient(ctrl)
	service := NewVolumeService(mockClient)
	ctx := context.Background()

	mockVol := volume.Volume{
		Name:   "vol1",
		Driver: "local",
	}

	mockClient.EXPECT().VolumeInspect(ctx, "vol1").Return(mockVol, nil)

	detail, err := service.VolumeInspect(ctx, "vol1")

	assert.NoError(t, err)
	assert.NotNil(t, detail)
	assert.Equal(t, "vol1", detail.Name)
}

func TestVolumeService_VolumeInspect_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_docker.NewMockDockerClient(ctrl)
	service := NewVolumeService(mockClient)
	ctx := context.Background()

	mockClient.EXPECT().VolumeInspect(ctx, "vol1").Return(volume.Volume{}, errors.New("inspect error"))

	detail, err := service.VolumeInspect(ctx, "vol1")

	assert.Error(t, err)
	assert.Nil(t, detail)
	assert.Contains(t, err.Error(), "docker inspect volume")
}

func TestVolumeService_VolumeCreate_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_docker.NewMockDockerClient(ctrl)
	service := NewVolumeService(mockClient)
	ctx := context.Background()

	req := VolumeCreateRequest{
		Name:   "vol1",
		Driver: "local",
	}

	mockVol := volume.Volume{
		Name:   "vol1",
		Driver: "local",
	}

	mockClient.EXPECT().VolumeCreate(ctx, volume.CreateOptions{Name: "vol1", Driver: "local"}).Return(mockVol, nil)

	detail, err := service.VolumeCreate(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, detail)
	assert.Equal(t, "vol1", detail.Name)
}

func TestVolumeService_VolumeCreate_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_docker.NewMockDockerClient(ctrl)
	service := NewVolumeService(mockClient)
	ctx := context.Background()

	req := VolumeCreateRequest{
		Name:   "vol1",
		Driver: "local",
	}

	mockClient.EXPECT().VolumeCreate(ctx, gomock.Any()).Return(volume.Volume{}, errors.New("create error"))

	detail, err := service.VolumeCreate(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, detail)
	assert.Contains(t, err.Error(), "docker create volume")
}

func TestVolumeService_VolumeRemove_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_docker.NewMockDockerClient(ctrl)
	service := NewVolumeService(mockClient)
	ctx := context.Background()

	mockClient.EXPECT().VolumeRemove(ctx, "vol1", true).Return(nil)

	err := service.VolumeRemove(ctx, "vol1", true)

	assert.NoError(t, err)
}

func TestVolumeService_VolumeRemove_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_docker.NewMockDockerClient(ctrl)
	service := NewVolumeService(mockClient)
	ctx := context.Background()

	mockClient.EXPECT().VolumeRemove(ctx, "vol1", true).Return(errors.New("remove error"))

	err := service.VolumeRemove(ctx, "vol1", true)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "docker remove volume")
}
