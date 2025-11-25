package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/gsc-lab/cs25-1-bannote-api-gateway/handlers"
	schedule_service "github.com/gsc-lab/cs25-1-bannote-api-gateway/routes/schedule-service"
	studyroom_service "github.com/gsc-lab/cs25-1-bannote-api-gateway/routes/studyroom-service"
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
	user_service.RegisterAllowedDomainRoutes(router)

	studyroom_service.RegisterRoomRoutes(router)
	studyroom_service.RegisterRoomOperatingRoutes(router)
	studyroom_service.RegisterRoomExceptionRoutes(router)

	schedule_service.RegisterTagsRoutes(router)
	schedule_service.RegisterGroupRoutes(router)
}
