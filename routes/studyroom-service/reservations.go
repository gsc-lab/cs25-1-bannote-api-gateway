package studyroom_service

import (
	"github.com/gin-gonic/gin"
	handlers "github.com/gsc-lab/cs25-1-bannote-api-gateway/handlers/studyroom-service"
	"github.com/gsc-lab/cs25-1-bannote-api-gateway/middleware"
)

func RegisterReservationRoutes(rg *gin.RouterGroup) {
	rooms := rg.Group("/reservations")

	{
		rooms.GET("", middleware.GRPCMetadata(), handlers.GetManyReservations)
		rooms.PATCH(":code", middleware.GRPCMetadata(), handlers.UpdateReservation)
		rooms.POST("", middleware.GRPCMetadata(), handlers.CreateReservation)
		rooms.DELETE(":code", middleware.GRPCMetadata(), handlers.DeleteReservation)
	}
}
