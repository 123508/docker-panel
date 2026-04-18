package handler

import (
	"bytes"
	"docker-panel/docker_cli_wrapper/mock"
	"docker-panel/internal/service"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/docker/docker/api/types/network"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

func setupNetworkTest(t *testing.T) (*gin.Engine, *mock.MockDockerClient) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	mockClient := mock.NewMockDockerClient(ctrl)
	svc := service.NewNetworkService(mockClient)
	h := NewNetworkHandler(svc)

	r := gin.New()
	v1 := r.Group("/api/v1")
	RegisterNetworkRoutes(v1, h)

	return r, mockClient
}

func TestNetworkHandler_List(t *testing.T) {
	r, mockClient := setupNetworkTest(t)

	t.Run("Success", func(t *testing.T) {
		mockClient.EXPECT().
			NetworkList(gomock.Any(), gomock.Any()).
			Return([]network.Inspect{{ID: "net1"}}, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/networks", nil)
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

func TestNetworkHandler_Inspect(t *testing.T) {
	r, mockClient := setupNetworkTest(t)

	t.Run("Success", func(t *testing.T) {
		mockClient.EXPECT().
			NetworkInspect(gomock.Any(), "net1", gomock.Any()).
			Return(network.Inspect{ID: "net1"}, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/networks/net1", nil)
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

func TestNetworkHandler_Create(t *testing.T) {
	r, mockClient := setupNetworkTest(t)

	t.Run("Success", func(t *testing.T) {
		mockClient.EXPECT().
			NetworkCreate(gomock.Any(), "new_net", gomock.Any()).
			Return(network.CreateResponse{ID: "new_net_id"}, nil)

		body := bytes.NewBufferString(`{"name": "new_net"}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/networks", body)
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
		req := httptest.NewRequest(http.MethodPost, "/api/v1/networks", body)
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

func TestNetworkHandler_Remove(t *testing.T) {
	r, mockClient := setupNetworkTest(t)

	t.Run("Success", func(t *testing.T) {
		mockClient.EXPECT().
			NetworkRemove(gomock.Any(), "net1").
			Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/networks/net1", nil)
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

func TestNetworkHandler_Connect(t *testing.T) {
	r, mockClient := setupNetworkTest(t)

	t.Run("Success", func(t *testing.T) {
		mockClient.EXPECT().
			NetworkConnect(gomock.Any(), "net1", "container1", gomock.Any()).
			Return(nil)

		body := bytes.NewBufferString(`{"container_id": "container1"}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/networks/net1/connect", body)
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
}

func TestNetworkHandler_Disconnect(t *testing.T) {
	r, mockClient := setupNetworkTest(t)

	t.Run("Success", func(t *testing.T) {
		mockClient.EXPECT().
			NetworkDisconnect(gomock.Any(), "net1", "container1", gomock.Any()).
			Return(nil)

		body := bytes.NewBufferString(`{"container_id": "container1"}`)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/networks/net1/disconnect", body)
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
}
