package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/gsc-lab/cs25-1-bannote-api-gateway/handlers"
)

func SetupRoutes(router *gin.Engine) {
	router.GET("/health", handlers.HealthCheck)

	auth := router.Group("/auth")
	auth.GET("/code/google", handlers.GoogleCallback)
}
