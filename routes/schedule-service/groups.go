package schedule_service

import (
	"github.com/gin-gonic/gin"
	handlers "github.com/gsc-lab/cs25-1-bannote-api-gateway/handlers/schedule-service"
	"github.com/gsc-lab/cs25-1-bannote-api-gateway/middleware"
)

func RegisterGroupRoutes(rg *gin.RouterGroup) {
	rooms := rg.Group("/schedule-groups")

	{
		rooms.GET("", middleware.GRPCMetadata(), handlers.ListTags)
		rooms.GET(":id", middleware.GRPCMetadata(), handlers.GetGroup)

		rooms.POST("", middleware.GRPCMetadata(), handlers.CreateGroup)
		rooms.DELETE(":id", middleware.GRPCMetadata(), handlers.DeleteGroup)
	}
}
