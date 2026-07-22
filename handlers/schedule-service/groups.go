package schedule_service

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	grouppb "github.com/gsc-lab/cs25-1-bannote-api-gateway/gen/go/schedule-service/group"
	"github.com/gsc-lab/cs25-1-bannote-api-gateway/grpc/client"
	"github.com/gsc-lab/cs25-1-bannote-api-gateway/utils"
)

func convertGroup(group *grouppb.Group) map[string]interface{} {
	result := map[string]interface{}{
		"id":                group.GetGroupId(),
		"group_code":        group.GetGroupCode(),
		"group_type_id":     group.GetGroupTypeId(),
		"group_name":        group.GetGroupName(),
		"group_description": group.GetGroupDescription(),
		"is_public":         group.GetIsPublic(),
		"is_published":      group.GetIsPublished(),
		"color_default":     group.GetColorDefault(),
		"color_highlight":   group.GetColorHighlight(),
		"created_at":        group.GetCreatedAt(),
		"updated_at":        group.GetUpdatedAt(),
		"created_by":        group.GetCreatedBy(),
		"updated_by":        group.GetUpdatedBy(),
		"bookmark":          group.GetBookmark(),
	}

	// Tags가 있으면 변환하여 추가
	if tags := group.GetTags(); len(tags) > 0 {
		convertedTags := make([]map[string]interface{}, len(tags))
		for i, tag := range tags {
			convertedTags[i] = map[string]interface{}{
				"id":         tag.GetTagId(),
				"name":       tag.GetName(),
				"created_by": tag.GetCreatedBy(),
				"created_at": tag.GetCreatedAt(),
			}
		}
		result["tags"] = convertedTags
	}

	return result
}

func convertGroups(groups []*grouppb.Group) []map[string]interface{} {
	result := make([]map[string]interface{}, len(groups))
	for i, group := range groups {
		result[i] = convertGroup(group)
	}
	return result
}

// TODO: 백엔드 페이지네이션 구현 되면 수정 필요
type listGroupsRequest struct {
	Page     int32    `form:"page"`
	Size     int32    `form:"size"`
	TagIds   []int64  `form:"tag_ids"`
	TagNames []string `form:"tag_names"`
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
		Page:     request.Page,
		PerPage:  request.Size,
		TagNames: request.TagNames,
	})

	if err != nil {
		utils.HandleGRPCError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  convertGroups(resp.GroupListResponse.Groups),
		"total": resp.GroupListResponse.TotalCount,
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
		GroupPermissionId: 1,
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

	convertedGroup := convertGroup(resp.Group)
	c.JSON(http.StatusOK, gin.H{
		"data": convertedGroup,
		"id":   convertedGroup["id"],
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

	convertedGroup := convertGroup(resp.Group)
	c.JSON(http.StatusOK, gin.H{
		"data": convertedGroup,
		"id":   convertedGroup["id"],
	})
}

func UpdateGroup(c *gin.Context) {
	scheduleClient := client.GetScheduleService(c)

	idStr := c.Param("id")

	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Required parameter"})
		return
	}

	id, err := strconv.ParseInt(idStr, 0, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Required parameter"})
		return
	}

	var request struct {
		GroupName        *string `json:"group_name"`
		GroupDescription *string `json:"group_description"`
		ColorDefault     *string `json:"color_default"`
		ColorHighlight   *string `json:"color_highlight"`
		IsPublic         *bool   `json:"is_public"`
		IsPublished      *bool   `json:"is_published"`
		TagIds           []int64 `json:"tag_ids"`
	}

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Required parameter", "details": err.Error()})
		return
	}

	ctx := utils.ContextWithMetadata(c)
	resp, err := scheduleClient.Group.UpdateGroup(ctx, &grouppb.UpdateGroupRequest{
		GroupId:          id,
		GroupName:        request.GroupName,
		GroupDescription: request.GroupDescription,
		ColorDefault:     request.ColorDefault,
		ColorHighlight:   request.ColorHighlight,
		IsPublic:         request.IsPublic,
		IsPublished:      request.IsPublished,
		TagIds:           request.TagIds,
	})

	if err != nil {
		utils.HandleGRPCError(c, err)
		return
	}

	convertedGroup := convertGroup(resp.Group)
	c.JSON(http.StatusOK, gin.H{
		"data": convertedGroup,
		"id":   convertedGroup["id"],
	})
}

func GetManyGroups(c *gin.Context) {
	scheduleClient := client.GetScheduleService(c)

	filterStr := c.Query("filter")

	if filterStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "filter parameter is required"})
		return
	}

	var filter struct {
		ID []int64 `json:"id"`
	}

	if err := json.Unmarshal([]byte(filterStr), &filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid filter format"})
		return
	}

	ctx := utils.ContextWithMetadata(c)
	resp, err := scheduleClient.Group.GetManyGroups(ctx, &grouppb.GetManyGroupsRequest{
		GroupIds: filter.ID,
	})

	if err != nil {
		utils.HandleGRPCError(c, err)
		return
	}

	c.JSON(http.StatusOK, convertGroups(resp.Groups))
}
