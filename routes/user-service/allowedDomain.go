package user_service

import (
	"github.com/gin-gonic/gin"
	handlers "github.com/gsc-lab/cs25-1-bannote-api-gateway/handlers/user-service"
	"github.com/gsc-lab/cs25-1-bannote-api-gateway/middleware"
)

func RegisterAllowedDomainRoutes(rg *gin.RouterGroup) {
	alloweddomains := rg.Group("/alloweddomains")

	{
		alloweddomains.GET("/check", handlers.CheckAllowedDomain)

		alloweddomains.GET("", middleware.GRPCMetadata(), handlers.ListAllowedDomain)
		alloweddomains.POST("", middleware.GRPCMetadata(), handlers.CreateAllowedDomain)
		alloweddomains.DELETE("/*id", middleware.GRPCMetadata(), handlers.DeleteAllowedDomain)
	}
}
