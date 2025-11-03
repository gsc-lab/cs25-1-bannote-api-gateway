package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gsc-lab/cs25-1-bannote-api-gateway/grpc/client"
	"github.com/gsc-lab/cs25-1-bannote-api-gateway/services"
)

func HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":    "UP",
		"timestamp": time.Now().Unix(),
	})
}

func ServiceHealthCheck(c *gin.Context) {
	container := &client.Container{
		UserService:  client.GetUserService(c),
		TokenService: client.GetTokenService(c),
	}

	result := services.CheckAllServices(container)

	statusCode := 200
	if result.OverallStatus == "DEGRADED" {
		statusCode = 503
	}

	c.JSON(statusCode, result)
}
