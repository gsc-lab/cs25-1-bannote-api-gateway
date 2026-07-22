package studyroom_service

import (
	"github.com/gin-gonic/gin"
	handlers "github.com/gsc-lab/cs25-1-bannote-api-gateway/handlers/studyroom-service"
	"github.com/gsc-lab/cs25-1-bannote-api-gateway/middleware"
)

func RegisterRoomExceptionRoutes(rg *gin.RouterGroup) {
	rooms := rg.Group("/room-exception")

	{
		rooms.GET(":id", middleware.GRPCMetadata(), handlers.GetRoomException)
		rooms.PATCH(":id", middleware.GRPCMetadata(), handlers.UpdateRoomException)
	}
}
