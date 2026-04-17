package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResponse_Success(t *testing.T) {
	data := map[string]string{"key": "value"}
	resp := Success(data)

	assert.Equal(t, 0, resp.Code)
	assert.Equal(t, "操作成功", resp.Message)
	assert.Equal(t, data, resp.Data)
}

func TestResponse_Error(t *testing.T) {
	code := ErrCodeInvalidParam
	msg := "参数错误"
	resp := Error(code, msg)

	assert.Equal(t, code, resp.Code)
	assert.Equal(t, msg, resp.Message)
	assert.Nil(t, resp.Data)
}

func TestResponse_Constants(t *testing.T) {
	assert.Equal(t, 1000, ErrCodeSystem)
	assert.Equal(t, 1001, ErrCodeInvalidParam)
	assert.Equal(t, 1004, ErrCodeNotFound)
	assert.Equal(t, 1401, ErrCodeUnauthorized)
	assert.Equal(t, 1403, ErrCodeForbidden)

	assert.Equal(t, 2000, ErrCodeDockerClient)
	assert.Equal(t, 2001, ErrCodeDockerAPI)
	assert.Equal(t, 2002, ErrCodeDockerTimeout)

	assert.Equal(t, 3001, ErrCodeContainerNotFound)
	assert.Equal(t, 3002, ErrCodeContainerRunning)
	assert.Equal(t, 3003, ErrCodeContainerStopped)
	assert.Equal(t, 3004, ErrCodeContainerConflict)

	assert.Equal(t, 4001, ErrCodeImageNotFound)
	assert.Equal(t, 4002, ErrCodeImageInUse)

	assert.Equal(t, 5001, ErrCodeNetworkNotFound)
	assert.Equal(t, 5002, ErrCodeNetworkInUse)

	assert.Equal(t, 6001, ErrCodeVolumeNotFound)
	assert.Equal(t, 6002, ErrCodeVolumeInUse)
}
