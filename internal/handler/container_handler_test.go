package handler

import (
	"bytes"
	mock_docker "docker-panel/docker_cli_wrapper/mock"
	"docker-panel/internal/service"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/docker/docker/api/types/container"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestRouter(svc *service.ContainerService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewContainerHandler(svc)
	RegisterContainerRoutes(r.Group("/api/v1"), h)
	return r
}

func TestContainerHandler_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_docker.NewMockDockerClient(ctrl)
	svc := service.NewContainerService(mockClient)
	r := setupTestRouter(svc)

	t.Run("success", func(t *testing.T) {
		mockClient.EXPECT().ContainerList(gomock.Any(), gomock.Any()).Return([]container.Summary{
			{ID: "c1", Names: []string{"/test-c1"}},
		}, nil)

		req := httptest.NewRequest("GET", "/api/v1/containers", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp service.Response
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, 0, resp.Code)
	})

	t.Run("error", func(t *testing.T) {
		mockClient.EXPECT().ContainerList(gomock.Any(), gomock.Any()).Return(nil, errors.New("list error"))

		req := httptest.NewRequest("GET", "/api/v1/containers", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var resp service.Response
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, service.ErrCodeDockerAPI, resp.Code)
	})

	t.Run("invalid_param", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/containers?all=invalid", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var resp service.Response
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, service.ErrCodeInvalidParam, resp.Code)
	})
}

func TestContainerHandler_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_docker.NewMockDockerClient(ctrl)
	svc := service.NewContainerService(mockClient)
	r := setupTestRouter(svc)

	t.Run("success", func(t *testing.T) {
		createReq := service.ContainerCreateRequest{
			Name:  "test-container",
			Image: "nginx:latest",
		}
		body, err := json.Marshal(createReq)
		require.NoError(t, err)

		mockClient.EXPECT().ContainerCreate(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), "test-container").
			Return(container.CreateResponse{ID: "new-c1"}, nil)

		req := httptest.NewRequest("POST", "/api/v1/containers", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp service.Response
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, 0, resp.Code)
	})

	t.Run("invalid json", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/containers", bytes.NewReader([]byte("{invalid")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var resp service.Response
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, service.ErrCodeInvalidParam, resp.Code)
	})
}

func TestContainerHandler_Start(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_docker.NewMockDockerClient(ctrl)
	svc := service.NewContainerService(mockClient)
	r := setupTestRouter(svc)

	t.Run("success", func(t *testing.T) {
		mockClient.EXPECT().ContainerStart(gomock.Any(), "c1", gomock.Any()).Return(nil)

		req := httptest.NewRequest("POST", "/api/v1/containers/c1/start", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestContainerHandler_Stop(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_docker.NewMockDockerClient(ctrl)
	svc := service.NewContainerService(mockClient)
	r := setupTestRouter(svc)

	t.Run("success", func(t *testing.T) {
		mockClient.EXPECT().ContainerStop(gomock.Any(), "c1", gomock.Any()).Return(nil)

		req := httptest.NewRequest("POST", "/api/v1/containers/c1/stop", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestContainerHandler_Remove(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_docker.NewMockDockerClient(ctrl)
	svc := service.NewContainerService(mockClient)
	r := setupTestRouter(svc)

	t.Run("success", func(t *testing.T) {
		mockClient.EXPECT().ContainerRemove(gomock.Any(), "c1", gomock.Any()).Return(nil)

		req := httptest.NewRequest("DELETE", "/api/v1/containers/c1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
