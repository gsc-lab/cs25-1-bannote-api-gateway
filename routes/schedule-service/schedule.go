package schedule_service

import (
	"github.com/gin-gonic/gin"
	handlers "github.com/gsc-lab/cs25-1-bannote-api-gateway/handlers/schedule-service"
	"github.com/gsc-lab/cs25-1-bannote-api-gateway/middleware"
)

func RegisterScheduleRoutes(rg *gin.RouterGroup) {
	groupsMember := rg.Group("/schedule")

	{
		groupsMember.GET("", middleware.GRPCMetadata(), handlers.GetScheduleList)
		groupsMember.GET(":schedule-id", middleware.GRPCMetadata(), handlers.GetSchedule)
		groupsMember.POST("", middleware.GRPCMetadata(), handlers.CreateSchedule)
		groupsMember.PATCH(":schedule-id", middleware.GRPCMetadata(), handlers.UpdateSchedule)
		groupsMember.DELETE(":schedule-id", middleware.GRPCMetadata(), handlers.DeleteSchedule)
	}
}
