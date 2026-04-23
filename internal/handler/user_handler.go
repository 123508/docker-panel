package handler

import (
	"net/http"

	"docker-panel/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, service.Error(service.ErrCodeInvalidParam, err.Error()))
		return
	}

	token, err := h.userService.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, service.Error(service.ErrCodeUnauthorized, "invalid username or password"))
		return
	}

	c.JSON(http.StatusOK, service.Success(gin.H{
		"token": token,
	}))
}
