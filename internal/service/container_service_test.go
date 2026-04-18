package service

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	mock_docker "docker-panel/docker_cli_wrapper/mock"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestContainerService_ContainerList_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_docker.NewMockDockerClient(ctrl)
	service := NewContainerService(mockClient)
	ctx := context.Background()
	req := ContainerListRequest{All: true, Limit: 10}

	mockContainers := []container.Summary{
		{
			ID:    "container_id_1",
			Names: []string{"/test_container"},
			Image: "nginx:latest",
		},
	}

	mockClient.EXPECT().ContainerList(ctx, gomock.Any()).Return(mockContainers, nil)

	res, err := service.ContainerList(ctx, req)
	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, "container_id_1", res[0].ID)
}

func TestContainerService_ContainerList_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_docker.NewMockDockerClient(ctrl)
	service := NewContainerService(mockClient)
	ctx := context.Background()

	mockClient.EXPECT().ContainerList(ctx, gomock.Any()).Return([]container.Summary{}, nil)

	res, err := service.ContainerList(ctx, ContainerListRequest{})
	assert.NoError(t, err)
	assert.Len(t, res, 0)
}

func TestContainerService_ContainerList_DockerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_docker.NewMockDockerClient(ctrl)
	service := NewContainerService(mockClient)
	ctx := context.Background()

	mockClient.EXPECT().ContainerList(ctx, gomock.Any()).Return(([]container.Summary)(nil), errors.New("docker error"))

	res, err := service.ContainerList(ctx, ContainerListRequest{})
	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "docker list containers")
}

func TestContainerService_ContainerInspect_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_docker.NewMockDockerClient(ctrl)
	service := NewContainerService(mockClient)
	ctx := context.Background()

	containerID := "test_id"
	mockInspect := container.InspectResponse{
		ContainerJSONBase: &container.ContainerJSONBase{
			ID:    containerID,
			Name:  "/test_container",
			State: &container.State{Status: "running", Running: true},
		},
		Config: &container.Config{
			Image: "nginx:latest",
		},
	}

	mockClient.EXPECT().ContainerInspect(ctx, containerID).Return(mockInspect, nil)

	res, err := service.ContainerInspect(ctx, containerID)
	assert.NoError(t, err)
	assert.Equal(t, containerID, res.ID)
	assert.Equal(t, "running", res.State.Status)
}

func TestContainerService_ContainerInspect_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_docker.NewMockDockerClient(ctrl)
	service := NewContainerService(mockClient)
	ctx := context.Background()

	mockClient.EXPECT().ContainerInspect(ctx, "invalid_id").Return(container.InspectResponse{}, errors.New("No such container: invalid_id"))

	res, err := service.ContainerInspect(ctx, "invalid_id")
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestContainerService_ContainerCreate_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_docker.NewMockDockerClient(ctrl)
	service := NewContainerService(mockClient)
	ctx := context.Background()

	req := ContainerCreateRequest{
		Name:  "new_container",
		Image: "ubuntu:latest",
	}

	mockResponse := container.CreateResponse{ID: "new_container_id"}
	mockClient.EXPECT().ContainerCreate(ctx, gomock.Any(), gomock.Any(), gomock.Any(), req.Name).Return(mockResponse, nil)

	id, err := service.ContainerCreate(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, "new_container_id", id)
}

func TestContainerService_ContainerCreate_Conflict(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_docker.NewMockDockerClient(ctrl)
	service := NewContainerService(mockClient)
	ctx := context.Background()

	req := ContainerCreateRequest{
		Name:  "existing_container",
		Image: "ubuntu:latest",
	}

	mockClient.EXPECT().ContainerCreate(ctx, gomock.Any(), gomock.Any(), gomock.Any(), req.Name).Return(container.CreateResponse{}, errors.New("Conflict. The name is already in use"))

	id, err := service.ContainerCreate(ctx, req)
	assert.Error(t, err)
	assert.Empty(t, id)
}

func TestContainerService_ContainerStart_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_docker.NewMockDockerClient(ctrl)
	service := NewContainerService(mockClient)
	ctx := context.Background()
	containerID := "test_id"

	mockClient.EXPECT().ContainerStart(ctx, containerID, container.StartOptions{}).Return(nil)

	err := service.ContainerStart(ctx, containerID)
	assert.NoError(t, err)
}

func TestContainerService_ContainerStart_AlreadyRunning(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_docker.NewMockDockerClient(ctrl)
	service := NewContainerService(mockClient)
	ctx := context.Background()
	containerID := "test_id"

	mockClient.EXPECT().ContainerStart(ctx, containerID, container.StartOptions{}).Return(errors.New("Container is already running"))

	err := service.ContainerStart(ctx, containerID)
	assert.Error(t, err)
}

func TestContainerService_ContainerStop_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_docker.NewMockDockerClient(ctrl)
	service := NewContainerService(mockClient)
	ctx := context.Background()
	containerID := "test_id"
	timeout := 10

	mockClient.EXPECT().ContainerStop(ctx, containerID, gomock.Any()).Return(nil)

	err := service.ContainerStop(ctx, containerID, timeout)
	assert.NoError(t, err)
}

func TestContainerService_ContainerRestart_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_docker.NewMockDockerClient(ctrl)
	service := NewContainerService(mockClient)
	ctx := context.Background()
	containerID := "test_id"

	mockClient.EXPECT().ContainerRestart(ctx, containerID, gomock.Any()).Return(nil)

	err := service.ContainerRestart(ctx, containerID, 10)
	assert.NoError(t, err)
}

func TestContainerService_ContainerRemove_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_docker.NewMockDockerClient(ctrl)
	service := NewContainerService(mockClient)
	ctx := context.Background()
	containerID := "test_id"

	mockClient.EXPECT().ContainerRemove(ctx, containerID, container.RemoveOptions{Force: false, RemoveVolumes: false}).Return(nil)

	err := service.ContainerRemove(ctx, containerID, false, false)
	assert.NoError(t, err)
}

func TestContainerService_ContainerRemove_Force(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_docker.NewMockDockerClient(ctrl)
	service := NewContainerService(mockClient)
	ctx := context.Background()
	containerID := "test_id"

	mockClient.EXPECT().ContainerRemove(ctx, containerID, container.RemoveOptions{Force: true, RemoveVolumes: false}).Return(nil)

	err := service.ContainerRemove(ctx, containerID, true, false)
	assert.NoError(t, err)
}

func TestContainerService_ContainerRemove_RunningWithoutForce(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_docker.NewMockDockerClient(ctrl)
	service := NewContainerService(mockClient)
	ctx := context.Background()
	containerID := "test_id"

	mockClient.EXPECT().ContainerRemove(ctx, containerID, container.RemoveOptions{Force: false, RemoveVolumes: false}).Return(errors.New("cannot remove a running container"))

	err := service.ContainerRemove(ctx, containerID, false, false)
	assert.Error(t, err)
}

func TestContainerService_ContainerPause_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_docker.NewMockDockerClient(ctrl)
	service := NewContainerService(mockClient)
	ctx := context.Background()
	containerID := "test_id"

	mockClient.EXPECT().ContainerPause(ctx, containerID).Return(nil)

	err := service.ContainerPause(ctx, containerID)
	assert.NoError(t, err)
}

func TestContainerService_ContainerUnpause_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_docker.NewMockDockerClient(ctrl)
	service := NewContainerService(mockClient)
	ctx := context.Background()
	containerID := "test_id"

	mockClient.EXPECT().ContainerUnpause(ctx, containerID).Return(nil)

	err := service.ContainerUnpause(ctx, containerID)
	assert.NoError(t, err)
}

func TestContainerService_ContainerLogs_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_docker.NewMockDockerClient(ctrl)
	service := NewContainerService(mockClient)
	ctx := context.Background()
	containerID := "test_id"
	opts := ContainerLogsOptions{ShowStdout: true}

	mockReader := io.NopCloser(strings.NewReader("log line"))
	mockClient.EXPECT().ContainerLogs(ctx, containerID, gomock.Any()).Return(mockReader, nil)

	reader, err := service.ContainerLogs(ctx, containerID, opts)
	assert.NoError(t, err)

	bytes, _ := io.ReadAll(reader)
	assert.Equal(t, "log line", string(bytes))
}

func TestContainerService_ContainerExec_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_docker.NewMockDockerClient(ctrl)
	service := NewContainerService(mockClient)
	ctx := context.Background()
	containerID := "test_id"
	req := ContainerExecRequest{Cmd: []string{"sh"}}

	execResp := container.ExecCreateResponse{ID: "exec_id"}
	hijackResp := types.HijackedResponse{}

	mockClient.EXPECT().ContainerExecCreate(ctx, containerID, gomock.Any()).Return(execResp, nil)
	mockClient.EXPECT().ContainerExecAttach(ctx, execResp.ID, gomock.Any()).Return(hijackResp, nil)

	res, err := service.ContainerExec(ctx, containerID, req)
	assert.NoError(t, err)
	assert.Equal(t, hijackResp, res)
}

func TestContainerService_ContainerExec_EmptyCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_docker.NewMockDockerClient(ctrl)
	service := NewContainerService(mockClient)
	ctx := context.Background()
	containerID := "test_id"
	req := ContainerExecRequest{Cmd: []string{}}

	mockClient.EXPECT().ContainerExecCreate(ctx, containerID, gomock.Any()).Return(container.ExecCreateResponse{}, errors.New("empty command"))

	_, err := service.ContainerExec(ctx, containerID, req)
	assert.Error(t, err)
}

func TestContainerService_ContainerRename_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_docker.NewMockDockerClient(ctrl)
	service := NewContainerService(mockClient)
	ctx := context.Background()
	containerID := "test_id"
	newName := "new_name"

	mockClient.EXPECT().ContainerRename(ctx, containerID, newName).Return(nil)

	err := service.ContainerRename(ctx, containerID, newName)
	assert.NoError(t, err)
}

func TestContainerService_ContainerRename_EmptyName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_docker.NewMockDockerClient(ctrl)
	service := NewContainerService(mockClient)
	ctx := context.Background()
	containerID := "test_id"

	mockClient.EXPECT().ContainerRename(ctx, containerID, "").Return(errors.New("invalid name"))

	err := service.ContainerRename(ctx, containerID, "")
	assert.Error(t, err)
}

func TestContainerService_ContainerRename_Conflict(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_docker.NewMockDockerClient(ctrl)
	service := NewContainerService(mockClient)
	ctx := context.Background()
	containerID := "test_id"
	newName := "existing_name"

	mockClient.EXPECT().ContainerRename(ctx, containerID, newName).Return(errors.New("name is already in use"))

	err := service.ContainerRename(ctx, containerID, newName)
	assert.Error(t, err)
}

func TestContainerService_ContainerTop_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_docker.NewMockDockerClient(ctrl)
	service := NewContainerService(mockClient)
	ctx := context.Background()
	containerID := "test_id"
	psArgs := "aux"

	topBody := container.TopResponse{
		Titles:    []string{"PID", "USER"},
		Processes: [][]string{{"1", "root"}},
	}

	mockClient.EXPECT().ContainerTop(ctx, containerID, []string{psArgs}).Return(topBody, nil)

	res, err := service.ContainerTop(ctx, containerID, psArgs)
	assert.NoError(t, err)
	assert.Equal(t, topBody.Titles, res.Titles)
	assert.Equal(t, topBody.Processes, res.Processes)
}

func TestContainerService_ContainerTop_NotRunning(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_docker.NewMockDockerClient(ctrl)
	service := NewContainerService(mockClient)
	ctx := context.Background()
	containerID := "test_id"
	psArgs := ""

	mockClient.EXPECT().ContainerTop(ctx, containerID, []string{}).Return(container.TopResponse{}, errors.New("Container is not running"))

	res, err := service.ContainerTop(ctx, containerID, psArgs)
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestContainerService_ContainerExport_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_docker.NewMockDockerClient(ctrl)
	service := NewContainerService(mockClient)
	ctx := context.Background()
	containerID := "test_id"

	mockReader := io.NopCloser(strings.NewReader("exported data"))
	mockClient.EXPECT().ContainerExport(ctx, containerID).Return(mockReader, nil)

	reader, err := service.ContainerExport(ctx, containerID)
	assert.NoError(t, err)

	bytes, _ := io.ReadAll(reader)
	assert.Equal(t, "exported data", string(bytes))
}
