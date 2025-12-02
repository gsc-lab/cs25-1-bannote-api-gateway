package schedule_service

import (
	"github.com/gin-gonic/gin"
	handlers "github.com/gsc-lab/cs25-1-bannote-api-gateway/handlers/schedule-service"
	"github.com/gsc-lab/cs25-1-bannote-api-gateway/middleware"
)

func RegisterGroupMemberRoutes(rg *gin.RouterGroup) {
	groupsMember := rg.Group("/group_member")

	{
		groupsMember.GET("", middleware.GRPCMetadata(), handlers.GetGroupsOfUser)
		groupsMember.GET(":group-id", middleware.GRPCMetadata(), handlers.GetUserInGroup)
		groupsMember.POST(":group-id", middleware.GRPCMetadata(), handlers.AddUserToGroup)
		groupsMember.DELETE(":group-id", middleware.GRPCMetadata(), handlers.RemoveUserFromGroup)
	}
}
