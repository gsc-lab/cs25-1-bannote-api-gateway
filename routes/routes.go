package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/gsc-lab/cs25-1-bannote-api-gateway/handlers"
	user_service "github.com/gsc-lab/cs25-1-bannote-api-gateway/routes/user-service"
)

func SetupRoutes(router *gin.RouterGroup) {
	router.GET("/health", handlers.HealthCheck)
	router.GET("/health/services", handlers.ServiceHealthCheck)

	auth := router.Group("/auth")
	auth.GET("/code/google", handlers.GoogleCallback)

	user_service.RegisterDepartmentRoutes(router)
	user_service.RegisterStudentClassRoute(router)
	user_service.RegisterUserRoute(router)
}
