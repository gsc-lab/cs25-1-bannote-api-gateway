package schedule_service

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	grouppb "github.com/gsc-lab/cs25-1-bannote-api-gateway/gen/go/schedule-service/group"
	"github.com/gsc-lab/cs25-1-bannote-api-gateway/grpc/client"
	"github.com/gsc-lab/cs25-1-bannote-api-gateway/utils"
)

// TODO: 백엔드 페이지네이션 구현 되면 수정 필요
type listGroupsRequest struct {
	Page   int32   `form:"page"`
	Size   int32   `form:"size"`
	TagIds []int64 `form:"tag_ids"`
}

func ListGroups(c *gin.Context) {
	scheduleClient := client.GetScheduleService(c)

	var request *listGroupsRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Required parameter", "details": err.Error()})
		return
	}

	ctx := utils.ContextWithMetadata(c)
	resp, err := scheduleClient.Group.GetGroupList(ctx, &grouppb.GetGroupListRequest{
		//Page: request.Page,
		//Size: request.Size,
		TagIds: request.TagIds,
	})

	if err != nil {
		utils.HandleGRPCError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  resp.GroupListResponse.Groups,
		"total": 10,
		"page":  request.Page,
		"size":  request.Size,
	})
}

type createGroupRequest struct {
	GroupName         string  `json:"group_name"`
	GroupDescription  string  `json:"group_description"`
	GroupPermissionId int64   `json:"group_permission_id"`
	GroupTypeId       int64   `json:"group_type_id"`
	ColorDefault      string  `json:"color_default"`
	ColorHighlight    string  `json:"color_highlight"`
	IsPublic          bool    `json:"is_public"`
	IsPublished       bool    `json:"is_published"`
	TagIds            []int64 `json:"tag_ids"`
}

func CreateGroup(c *gin.Context) {
	scheduleClient := client.GetScheduleService(c)

	var request *createGroupRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Required parameter", "details": err.Error()})
		return
	}

	ctx := utils.ContextWithMetadata(c)
	resp, err := scheduleClient.Group.CreateGroup(ctx, &grouppb.CreateGroupRequest{
		GroupName:         request.GroupName,
		GroupDescription:  request.GroupDescription,
		GroupPermissionId: request.GroupPermissionId,
		GroupTypeId:       request.GroupTypeId,
		ColorDefault:      request.ColorDefault,
		ColorHighlight:    request.ColorHighlight,
		IsPublic:          request.IsPublic,
		IsPublished:       request.IsPublished,
		TagIds:            request.TagIds,
	})

	if err != nil {
		utils.HandleGRPCError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": resp.Group,
		"id":   resp.Group.GroupId,
	})
}

func DeleteGroup(c *gin.Context) {
	scheduleClient := client.GetScheduleService(c)

	idStr := c.Param("id")

	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Required parameter"})
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format", "details": err.Error()})
		return
	}

	ctx := utils.ContextWithMetadata(c)
	_, err = scheduleClient.Group.DeleteGroup(ctx, &grouppb.DeleteGroupRequest{
		GroupId: id,
	})

	if err != nil {
		utils.HandleGRPCError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Room deleted successfully",
	})
}

func GetGroup(c *gin.Context) {
	scheduleClient := client.GetScheduleService(c)

	idStr := c.Param("id")

	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Required parameter"})
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format", "details": err.Error()})
		return
	}

	ctx := utils.ContextWithMetadata(c)
	resp, err := scheduleClient.Group.GetGroup(ctx, &grouppb.GetGroupRequest{
		GroupId: id,
	})

	if err != nil {
		utils.HandleGRPCError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": resp.Group,
		"id":   resp.Group.GroupId,
	})
}
