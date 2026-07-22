package schedule_service

import (
	"github.com/gin-gonic/gin"
	handlers "github.com/gsc-lab/cs25-1-bannote-api-gateway/handlers/schedule-service"
	"github.com/gsc-lab/cs25-1-bannote-api-gateway/middleware"
)

func RegisterGroupRoutes(rg *gin.RouterGroup) {
	groups := rg.Group("/schedule-groups")

	{
		groups.GET("", middleware.GRPCMetadata(), handlers.ListGroups)
		groups.GET(":id", middleware.GRPCMetadata(), handlers.GetGroup)
		groups.GET("/many", middleware.GRPCMetadata(), handlers.GetManyGroups)

		groups.POST("", middleware.GRPCMetadata(), handlers.CreateGroup)
		groups.PATCH(":id", middleware.GRPCMetadata(), handlers.UpdateGroup)
		groups.DELETE(":id", middleware.GRPCMetadata(), handlers.DeleteGroup)
	}
}
