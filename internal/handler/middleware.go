package handler

import (
	"docker-panel/internal/service"
	"docker-panel/internal/utils"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func RecoveryMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, err interface{}) {
		c.JSON(http.StatusInternalServerError, service.Error(service.ErrCodeSystem, "internal server error"))
	})
}

func LoggerMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return param.TimeStamp.Format(time.RFC3339) + " | " +
			param.Method + " " + param.Path + " | " +
			strings.TrimSpace(param.StatusCodeColor()+strings.TrimSpace(strconv.Itoa(param.StatusCode))) + "\n"
	})
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, service.Error(service.ErrCodeUnauthorized, "missing authorization header"))
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, service.Error(service.ErrCodeUnauthorized, "invalid authorization format"))
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, service.Error(service.ErrCodeUnauthorized, "invalid token"))
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Next()
	}
}

func respondJSON(c *gin.Context, resp service.Response) {
	var status int
	switch resp.Code {
	case service.ErrCodeUnauthorized:
		status = http.StatusUnauthorized
	case service.ErrCodeForbidden:
		status = http.StatusForbidden
	case service.ErrCodeNotFound,
		service.ErrCodeContainerNotFound,
		service.ErrCodeImageNotFound,
		service.ErrCodeNetworkNotFound,
		service.ErrCodeVolumeNotFound:
		status = http.StatusNotFound
	case 0:
		status = http.StatusOK
	default:
		status = http.StatusBadRequest
	}
	c.JSON(status, resp)
}
