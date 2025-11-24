package user_service

import (
	"github.com/gin-gonic/gin"
	handlers "github.com/gsc-lab/cs25-1-bannote-api-gateway/handlers/user-service"
	"github.com/gsc-lab/cs25-1-bannote-api-gateway/middleware"
)

func RegisterUserRoute(rg *gin.RouterGroup) {
	user := rg.Group("/users")

	{
		user.GET("", middleware.GRPCMetadata(), handlers.ListUsers)
		user.POST("register", handlers.CreateUser)
		user.PATCH(":user_code", middleware.GRPCMetadata(), handlers.UpdateUser)
	}
}
