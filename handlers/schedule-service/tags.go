package schedule_service

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	tagpb "github.com/gsc-lab/cs25-1-bannote-api-gateway/gen/go/schedule-service/tag"
	"github.com/gsc-lab/cs25-1-bannote-api-gateway/grpc/client"
	"github.com/gsc-lab/cs25-1-bannote-api-gateway/utils"
)

func convertTag(tag *tagpb.Tag) map[string]interface{} {
	return map[string]interface{}{
		"id":         tag.GetTagId(),
		"name":       tag.GetName(),
		"created_by": tag.GetCreatedBy(),
		"created_at": tag.GetCreatedAt(),
	}
}

func convertTags(tags []*tagpb.Tag) []map[string]interface{} {
	result := make([]map[string]interface{}, len(tags))
	for i, tag := range tags {
		result[i] = convertTag(tag)
	}
	return result
}

type listTagsRequest struct {
	Page int32 `form:"page"`
	Size int32 `form:"size"`
}

func ListTags(c *gin.Context) {
	scheduleClient := client.GetScheduleService(c)

	var request *listTagsRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Required parameter", "details": err.Error()})
		return
	}

	ctx := utils.ContextWithMetadata(c)
	resp, err := scheduleClient.Tag.GetTagList(ctx, &tagpb.GetTagListRequest{
		Page:    request.Page,
		PerPage: request.Size,
	})

	if err != nil {
		utils.HandleGRPCError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  convertTags(resp.TagListResponse.Tags),
		"total": resp.TagListResponse.TotalCount,
		"page":  request.Page,
		"size":  request.Size,
	})
}

type createTagRequest struct {
	Name string `json:"name"`
}

func CreateTag(c *gin.Context) {
	scheduleClient := client.GetScheduleService(c)

	var request *createTagRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Required parameter", "details": err.Error()})
		return
	}

	ctx := utils.ContextWithMetadata(c)
	resp, err := scheduleClient.Tag.CreateTag(ctx, &tagpb.CreateTagRequest{
		Name: request.Name,
	})

	if err != nil {
		utils.HandleGRPCError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": convertTag(resp.Tag),
		"id":   resp.Tag.TagId,
	})
}

func DeleteTag(c *gin.Context) {
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
	_, err = scheduleClient.Tag.DeleteTag(ctx, &tagpb.DeleteTagRequest{
		TagId: id,
	})

	if err != nil {
		utils.HandleGRPCError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Room deleted successfully",
	})
}

func GetTag(c *gin.Context) {
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
	resp, err := scheduleClient.Tag.GetTag(ctx, &tagpb.GetTagRequest{
		TagId: id,
	})

	if err != nil {
		utils.HandleGRPCError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": convertTag(resp.Tag),
		"id":   resp.Tag.TagId,
	})
}

func GetManyTags(c *gin.Context) {
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
	resp, err := scheduleClient.Tag.GetManyTags(ctx, &tagpb.GetManyTagsRequest{
		TagIds: filter.ID,
	})

	if err != nil {
		utils.HandleGRPCError(c, err)
		return
	}

	c.JSON(http.StatusOK, convertTags(resp.Tags))
}
