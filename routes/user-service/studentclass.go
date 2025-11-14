package user_service

import (
	"github.com/gin-gonic/gin"
	handlers "github.com/gsc-lab/cs25-1-bannote-api-gateway/handlers/user-service"
	"github.com/gsc-lab/cs25-1-bannote-api-gateway/middleware"
)

func RegisterStudentClassRoute(rg *gin.RouterGroup) {
	studentClass := rg.Group("/studentclasses")

	{
		studentClass.GET(":id", handlers.GetStudentClass)
		studentClass.GET("", handlers.ListStudentClass)
		studentClass.GET("/many", handlers.GetManyStudentClasses)

		studentClass.POST("", middleware.GRPCMetadata(), handlers.CreateStudentClass)
		studentClass.PATCH(":id", middleware.GRPCMetadata(), handlers.UpdateStudentClass)
		studentClass.DELETE(":id", middleware.GRPCMetadata(), handlers.DeleteStudentClass)
	}
}
