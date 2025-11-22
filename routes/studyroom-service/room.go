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
		rooms.POST("", middleware.GRPCMetadata(), handlers.CreateRoom)
		rooms.DELETE(":id", middleware.GRPCMetadata(), handlers.DeleteRoom)
	}
}
