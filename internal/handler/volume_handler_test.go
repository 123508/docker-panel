package handler

import (
	"bytes"
	"docker-panel/docker_cli_wrapper/mock"
	"docker-panel/internal/service"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/docker/docker/api/types/volume"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

func setupVolumeTest(t *testing.T) (*gin.Engine, *mock.MockDockerClient) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	mockClient := mock.NewMockDockerClient(ctrl)
	svc := service.NewVolumeService(mockClient)
	h := NewVolumeHandler(svc)

	r := gin.New()
	v1 := r.Group("/api/v1")
	RegisterVolumeRoutes(v1, h)

	return r, mockClient
}

func TestVolumeHandler_List(t *testing.T) {
	r, mockClient := setupVolumeTest(t)

	t.Run("Success", func(t *testing.T) {
		mockClient.EXPECT().
			VolumeList(gomock.Any(), gomock.Any()).
			Return(volume.ListResponse{Volumes: []*volume.Volume{{Name: "vol1"}}}, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/volumes", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}
		var resp service.Response
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}
		if resp.Code != 0 {
			t.Errorf("expected code 0, got %d", resp.Code)
		}
	})
}

func TestVolumeHandler_Inspect(t *testing.T) {
	r, mockClient := setupVolumeTest(t)

	t.Run("Success", func(t *testing.T) {
		mockClient.EXPECT().
			VolumeInspect(gomock.Any(), "vol1").
			Return(volume.Volume{Name: "vol1"}, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/volumes/vol1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}
		var resp service.Response
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}
		if resp.Code != 0 {
			t.Errorf("expected code 0, got %d", resp.Code)
		}
	})
}

func TestVolumeHandler_Create(t *testing.T) {
	r, mockClient := setupVolumeTest(t)

	t.Run("Success", func(t *testing.T) {
		mockClient.EXPECT().
			VolumeCreate(gomock.Any(), gomock.Any()).
			Return(volume.Volume{Name: "new_vol"}, nil)

		body := bytes.NewBufferString(`{"name": "new_vol"}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/volumes", body)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}
		var resp service.Response
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}
		if resp.Code != 0 {
			t.Errorf("expected code 0, got %d", resp.Code)
		}
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		body := bytes.NewBufferString(`{invalid json`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/volumes", body)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
		var resp service.Response
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}
		if resp.Code == 0 {
			t.Errorf("expected non-zero code, got %d", resp.Code)
		}
	})
}

func TestVolumeHandler_Remove(t *testing.T) {
	r, mockClient := setupVolumeTest(t)

	t.Run("Success", func(t *testing.T) {
		mockClient.EXPECT().
			VolumeRemove(gomock.Any(), "vol1", gomock.Any()).
			Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/volumes/vol1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}
		var resp service.Response
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}
		if resp.Code != 0 {
			t.Errorf("expected code 0, got %d", resp.Code)
		}
	})
}
