package handler

import (
	"bytes"
	"docker-panel/docker_cli_wrapper/mock"
	"docker-panel/internal/service"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/docker/docker/api/types/image"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

func setupImageTest(t *testing.T) (*gin.Engine, *mock.MockDockerClient) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	mockClient := mock.NewMockDockerClient(ctrl)
	svc := service.NewImageService(mockClient)
	h := NewImageHandler(svc)

	r := gin.New()
	v1 := r.Group("/api/v1")
	RegisterImageRoutes(v1, h)

	return r, mockClient
}

func TestImageHandler_List(t *testing.T) {
	r, mockClient := setupImageTest(t)

	t.Run("Success", func(t *testing.T) {
		mockClient.EXPECT().
			ImageList(gomock.Any(), gomock.Any()).
			Return([]image.Summary{{ID: "img1"}}, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/images", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}
		var resp service.Response
		json.Unmarshal(w.Body.Bytes(), &resp)
		if resp.Code != 0 {
			t.Errorf("expected code 0, got %d", resp.Code)
		}
	})

	t.Run("InvalidParams", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/images?all=invalid", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
		var resp service.Response
		json.Unmarshal(w.Body.Bytes(), &resp)
		if resp.Code == 0 {
			t.Errorf("expected non-zero code, got %d", resp.Code)
		}
	})
}

func TestImageHandler_Inspect(t *testing.T) {
	r, mockClient := setupImageTest(t)

	t.Run("Success", func(t *testing.T) {
		mockClient.EXPECT().
			ImageInspect(gomock.Any(), "img1").
			Return(image.InspectResponse{ID: "img1"}, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/images/img1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}
		var resp service.Response
		json.Unmarshal(w.Body.Bytes(), &resp)
		if resp.Code != 0 {
			t.Errorf("expected code 0, got %d", resp.Code)
		}
	})
}

func TestImageHandler_Remove(t *testing.T) {
	r, mockClient := setupImageTest(t)

	t.Run("Success", func(t *testing.T) {
		mockClient.EXPECT().
			ImageRemove(gomock.Any(), "img1", gomock.Any()).
			Return([]image.DeleteResponse{{Deleted: "img1"}}, nil)

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/images/img1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}
		var resp service.Response
		json.Unmarshal(w.Body.Bytes(), &resp)
		if resp.Code != 0 {
			t.Errorf("expected code 0, got %d", resp.Code)
		}
	})
}

func TestImageHandler_Pull(t *testing.T) {
	r, mockClient := setupImageTest(t)

	t.Run("Success", func(t *testing.T) {
		mockReadCloser := io.NopCloser(strings.NewReader("pulling..."))
		mockClient.EXPECT().
			ImagePull(gomock.Any(), "nginx:latest", gomock.Any()).
			Return(mockReadCloser, nil)

		body := bytes.NewBufferString(`{"image": "nginx:latest"}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/images/pull", body)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}
		if w.Body.String() != "pulling..." {
			t.Errorf("expected body 'pulling...', got %s", w.Body.String())
		}
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		body := bytes.NewBufferString(`{invalid json`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/images/pull", body)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
		var resp service.Response
		json.Unmarshal(w.Body.Bytes(), &resp)
		if resp.Code == 0 {
			t.Errorf("expected non-zero code, got %d", resp.Code)
		}
	})
}

func TestImageHandler_Tag(t *testing.T) {
	r, mockClient := setupImageTest(t)

	t.Run("Success", func(t *testing.T) {
		mockClient.EXPECT().
			ImageTag(gomock.Any(), "img1", "myrepo/img1:v1").
			Return(nil)

		body := bytes.NewBufferString(`{"source": "img1", "target": "myrepo/img1:v1"}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/images/tag", body)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}
		var resp service.Response
		json.Unmarshal(w.Body.Bytes(), &resp)
		if resp.Code != 0 {
			t.Errorf("expected code 0, got %d", resp.Code)
		}
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		body := bytes.NewBufferString(`{invalid json`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/images/tag", body)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
		var resp service.Response
		json.Unmarshal(w.Body.Bytes(), &resp)
		if resp.Code == 0 {
			t.Errorf("expected non-zero code, got %d", resp.Code)
		}
	})
}
