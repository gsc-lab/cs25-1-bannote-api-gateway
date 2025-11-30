package schedule_service

import (
	"github.com/gin-gonic/gin"
	handlers "github.com/gsc-lab/cs25-1-bannote-api-gateway/handlers/schedule-service"
	"github.com/gsc-lab/cs25-1-bannote-api-gateway/middleware"
)

func RegisterTagsRoutes(rg *gin.RouterGroup) {
	tags := rg.Group("/schedule-tags")

	{
		tags.GET("", middleware.GRPCMetadata(), handlers.ListTags)
		tags.GET(":id", middleware.GRPCMetadata(), handlers.GetTag)

		tags.POST("", middleware.GRPCMetadata(), handlers.CreateTag)
		tags.DELETE(":id", middleware.GRPCMetadata(), handlers.DeleteTag)
	}
}
