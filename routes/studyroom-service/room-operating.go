package studyroom_service

import (
	"github.com/gin-gonic/gin"
	handlers "github.com/gsc-lab/cs25-1-bannote-api-gateway/handlers/studyroom-service"
	"github.com/gsc-lab/cs25-1-bannote-api-gateway/middleware"
)

func RegisterRoomOperatingRoutes(rg *gin.RouterGroup) {
	rooms := rg.Group("/room-operating")

	{
		rooms.GET(":id", middleware.GRPCMetadata(), handlers.GetRoomOperating)

		rooms.PATCH(":id", middleware.GRPCMetadata(), handlers.UpdateRoom)
		rooms.POST("", middleware.GRPCMetadata(), handlers.CreateRoom)
		rooms.DELETE(":id", middleware.GRPCMetadata(), handlers.DeleteRoom)
	}
}
