package studyroom_service

import (
	"github.com/gin-gonic/gin"
	handlers "github.com/gsc-lab/cs25-1-bannote-api-gateway/handlers/studyroom-service"
	"github.com/gsc-lab/cs25-1-bannote-api-gateway/middleware"
)

func RegisterRoomRoutes(rg *gin.RouterGroup) {
	rooms := rg.Group("/studyrooms")

	{
		rooms.GET("", middleware.GRPCMetadata(), handlers.ListRoom)
		rooms.GET(":id", middleware.GRPCMetadata(), handlers.GetRoom)

		rooms.PATCH(":id", middleware.GRPCMetadata(), handlers.UpdateRoom)
		rooms.POST("", middleware.GRPCMetadata(), handlers.CreateRoom)
		rooms.DELETE(":id", middleware.GRPCMetadata(), handlers.DeleteRoom)
	}
}
