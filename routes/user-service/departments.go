package user_service

import (
	"github.com/gin-gonic/gin"
	handlers "github.com/gsc-lab/cs25-1-bannote-api-gateway/handlers/user-service"
	"github.com/gsc-lab/cs25-1-bannote-api-gateway/middleware"
)

func RegisterDepartmentRoutes(rg *gin.RouterGroup) {
	departments := rg.Group("/departments")
	departments.Use(middleware.GRPCMetadata())
	{
		departments.GET("", handlers.ListDepartments)
		departments.GET(":id", handlers.GetDepartment)
	}
}
