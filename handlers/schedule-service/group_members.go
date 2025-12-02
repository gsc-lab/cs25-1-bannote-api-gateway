package schedule_service

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	usergrouppb "github.com/gsc-lab/cs25-1-bannote-api-gateway/gen/go/schedule-service/user_group"
	"github.com/gsc-lab/cs25-1-bannote-api-gateway/grpc/client"
	"github.com/gsc-lab/cs25-1-bannote-api-gateway/utils"
)

func AddUserToGroup(c *gin.Context) {
	scheduleClient := client.GetScheduleService(c)

	idStr := c.Param("group-id")

	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Required parameter"})
		return
	}

	groupId, err := strconv.ParseInt(idStr, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group_id format", "details": err.Error()})
		return
	}

	var request struct {
		UserId string `json:"user_id"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Required parameter", "details": err.Error()})
		return
	}

	fmt.Println("debug - groupId:", groupId, "userId:", request.UserId)

	ctx := utils.ContextWithMetadata(c)
	resp, err := scheduleClient.UserGroup.AddUserToGroup(ctx, &usergrouppb.AddUserToGroupRequest{
		GroupId: groupId,
		UserId:  request.UserId,
	})

	if err != nil {
		utils.HandleGRPCError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": resp,
	})
}

func RemoveUserFromGroup(c *gin.Context) {
	scheduleClient := client.GetScheduleService(c)

	idStr := c.Param("group-id")

	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Required parameter"})
		return
	}

	groupId, err := strconv.ParseInt(idStr, 0, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Required parameter"})
		return
	}

	var request struct {
		UserId string `json:"user_id"`
	}

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Required parameter", "details": err.Error()})
		return
	}

	ctx := utils.ContextWithMetadata(c)
	resp, err := scheduleClient.UserGroup.RemoveUserFromGroup(ctx, &usergrouppb.RemoveUserFromGroupRequest{
		GroupId: groupId,
		UserId:  request.UserId,
	})

	if err != nil {
		utils.HandleGRPCError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func GetUserInGroup(c *gin.Context) {
	scheduleClient := client.GetScheduleService(c)

	idStr := c.Param("group-id")

	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Required parameter"})
		return
	}

	groupId, err := strconv.ParseInt(idStr, 0, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Required parameter"})
		return
	}

	ctx := utils.ContextWithMetadata(c)
	resp, err := scheduleClient.UserGroup.GetUsersInGroup(ctx, &usergrouppb.GetUsersInGroupRequest{
		GroupId: groupId,
	})

	if err != nil {
		utils.HandleGRPCError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": resp.Users,
	})
}

func GetGroupsOfUser(c *gin.Context) {

	scheduleClient := client.GetScheduleService(c)

	userId := c.Query("user_id")

	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Required parameter"})
		return
	}

	ctx := utils.ContextWithMetadata(c)
	resp, err := scheduleClient.UserGroup.GetGroupsOfUser(ctx, &usergrouppb.GetGroupsOfUserRequest{
		UserId: userId,
	})

	if err != nil {
		utils.HandleGRPCError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": resp.Groups,
	})
}
