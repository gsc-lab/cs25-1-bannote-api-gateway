package schedule_service

import (
	"github.com/gin-gonic/gin"
	handlers "github.com/gsc-lab/cs25-1-bannote-api-gateway/handlers/schedule-service"
	"github.com/gsc-lab/cs25-1-bannote-api-gateway/middleware"
)

func RegisterTagsRoutes(rg *gin.RouterGroup) {
	rooms := rg.Group("/schedule-tags")

	{
		rooms.GET("", middleware.GRPCMetadata(), handlers.ListTags)
		rooms.GET(":id", middleware.GRPCMetadata(), handlers.GetTag)

		rooms.POST("", middleware.GRPCMetadata(), handlers.CreateTag)
		rooms.DELETE(":id", middleware.GRPCMetadata(), handlers.DeleteTag)
	}
}
