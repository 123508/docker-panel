package service

import (
	"context"
	mock_docker "docker-panel/docker_cli_wrapper/mock"
	"errors"
	"testing"

	"github.com/docker/docker/api/types/network"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNetworkService_NetworkList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_docker.NewMockDockerClient(ctrl)
	svc := NewNetworkService(mockClient)

	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockItems := []network.Inspect{
			{
				ID:     "net1",
				Name:   "bridge",
				Driver: "bridge",
				Scope:  "local",
				IPAM: network.IPAM{
					Driver: "default",
					Config: []network.IPAMConfig{
						{Subnet: "172.17.0.0/16", Gateway: "172.17.0.1"},
					},
				},
			},
		}
		mockClient.EXPECT().NetworkList(gomock.Any(), gomock.Any()).Return(mockItems, nil)

		items, err := svc.NetworkList(ctx)
		require.NoError(t, err)
		assert.Len(t, items, 1)
		assert.Equal(t, "net1", items[0].ID)
		assert.Equal(t, "bridge", items[0].Name)
		assert.Equal(t, "172.17.0.0/16", items[0].IPAM.Config[0].Subnet)
	})

	t.Run("error", func(t *testing.T) {
		mockClient.EXPECT().NetworkList(ctx, network.ListOptions{}).Return(nil, errors.New("list error"))
		items, err := svc.NetworkList(ctx)
		require.Error(t, err)
		assert.Nil(t, items)
		assert.Contains(t, err.Error(), "list error")
	})
}

func TestNetworkService_NetworkInspect(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_docker.NewMockDockerClient(ctrl)
	svc := NewNetworkService(mockClient)

	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockClient.EXPECT().NetworkInspect(ctx, "net1", network.InspectOptions{}).Return(network.Inspect{
			ID:     "net1",
			Name:   "bridge",
			Driver: "bridge",
		}, nil)

		item, err := svc.NetworkInspect(ctx, "net1")
		require.NoError(t, err)
		assert.NotNil(t, item)
		assert.Equal(t, "net1", item.ID)
		assert.Equal(t, "bridge", item.Name)
	})

	t.Run("error", func(t *testing.T) {
		mockClient.EXPECT().NetworkInspect(ctx, "net1", network.InspectOptions{}).Return(network.Inspect{}, errors.New("inspect error"))
		item, err := svc.NetworkInspect(ctx, "net1")
		require.Error(t, err)
		assert.Nil(t, item)
	})
}

func TestNetworkService_NetworkCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_docker.NewMockDockerClient(ctrl)
	svc := NewNetworkService(mockClient)

	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		req := NetworkCreateRequest{
			Name:   "new-net",
			Driver: "bridge",
			IPAM: &NetworkIPAM{
				Driver: "default",
				Config: []IPAMConfig{
					{Subnet: "10.0.0.0/24", Gateway: "10.0.0.1"},
				},
			},
		}

		expectedOpts := network.CreateOptions{
			Driver: "bridge",
			IPAM: &network.IPAM{
				Driver: "default",
				Config: []network.IPAMConfig{
					{Subnet: "10.0.0.0/24", Gateway: "10.0.0.1"},
				},
			},
		}

		mockClient.EXPECT().NetworkCreate(ctx, "new-net", expectedOpts).Return(network.CreateResponse{ID: "net-123"}, nil)

		id, err := svc.NetworkCreate(ctx, req)
		require.NoError(t, err)
		assert.Equal(t, "net-123", id)
	})
}

func TestNetworkService_NetworkRemove(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_docker.NewMockDockerClient(ctrl)
	svc := NewNetworkService(mockClient)
	ctx := context.Background()

	mockClient.EXPECT().NetworkRemove(ctx, "net1").Return(nil)
	err := svc.NetworkRemove(ctx, "net1")
	require.NoError(t, err)
}

func TestNetworkService_NetworkConnect(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_docker.NewMockDockerClient(ctrl)
	svc := NewNetworkService(mockClient)
	ctx := context.Background()

	mockClient.EXPECT().NetworkConnect(ctx, "net1", "container1", (*network.EndpointSettings)(nil)).Return(nil)
	err := svc.NetworkConnect(ctx, "net1", "container1")
	require.NoError(t, err)
}

func TestNetworkService_NetworkDisconnect(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_docker.NewMockDockerClient(ctrl)
	svc := NewNetworkService(mockClient)
	ctx := context.Background()

	mockClient.EXPECT().NetworkDisconnect(ctx, "net1", "container1", true).Return(nil)
	err := svc.NetworkDisconnect(ctx, "net1", "container1", true)
	require.NoError(t, err)
}
