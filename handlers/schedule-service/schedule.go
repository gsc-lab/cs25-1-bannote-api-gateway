package schedule_service

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	schedulepb "github.com/gsc-lab/cs25-1-bannote-api-gateway/gen/go/schedule-service/schedule"
	"github.com/gsc-lab/cs25-1-bannote-api-gateway/grpc/client"
	"github.com/gsc-lab/cs25-1-bannote-api-gateway/utils"
)

func GetScheduleList(c *gin.Context) {
	scheduleClient := client.GetScheduleService(c)

	var request struct {
		EndDate   *string `form:"end_date"`
		StartDate *string `form:"start_date"`
		GroupIds  string  `form:"group_ids"`
		Page      int32   `form:"page"`
		Size      int32   `form:"size"`
	}

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Required parameter", "details": err.Error()})
		return
	}

	// Convert GroupIds from comma-separated string to []int64
	var groupIds []int64
	if request.GroupIds != "" {
		groupIdStrs := strings.Split(request.GroupIds, ",")
		groupIds = make([]int64, 0, len(groupIdStrs))
		for _, idStr := range groupIdStrs {
			idStr = strings.TrimSpace(idStr)
			if idStr == "" {
				continue
			}
			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group_id format", "details": err.Error()})
				return
			}
			groupIds = append(groupIds, id)
		}
	}

	ctx := utils.ContextWithMetadata(c)
	resp, err := scheduleClient.Schedule.GetScheduleList(ctx, &schedulepb.GetScheduleListRequest{
		StartAt:  request.StartDate,
		EndAt:    request.EndDate,
		GroupIds: groupIds,
		Page:     request.Page,
		PerPage:  request.Size,
	})

	if err != nil {
		utils.HandleGRPCError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": resp.ScheduleListResponse.Schedules,
	})
}

func GetSchedule(c *gin.Context) {
	scheduleClient := client.GetScheduleService(c)

	strId := c.Param("schedule-id")

	if strId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "filter parameter is required"})
		return
	}

	scheduleId, err := strconv.ParseInt(strId, 0, 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid schedule_id format", "details": err.Error()})
		return
	}

	ctx := utils.ContextWithMetadata(c)
	resp, err := scheduleClient.Schedule.GetSchedule(ctx, &schedulepb.GetScheduleRequest{
		ScheduleId: scheduleId,
	})

	if err != nil {
		utils.HandleGRPCError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": resp.Schedule,
	})
}

func CreateSchedule(c *gin.Context) {
	scheduleClient := client.GetScheduleService(c)

	var body struct {
		Title         string  `json:"title"`
		Description   string  `json:"description"`
		StartAt       string  `json:"start_at"`
		EndAt         string  `json:"end_at"`
		IsAllDay      bool    `json:"is_all_day"`
		PlaceId       *int64  `json:"place_id"`
		PlaceText     *string `json:"place_text"`
		GroupId       int64   `json:"group_id"`
		Comment       *string `json:"comment"`
		IsHighlighted bool    `json:"is_highlighted"`
	}

	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Required parameter", "details": err.Error()})
		return
	}

	linkData := &schedulepb.ScheduleLinkData{
		Title:       body.Title,
		Description: body.Description,
		StartAt:     body.StartAt,
		EndAt:       body.EndAt,
		PlaceId:     body.PlaceId,
		PlaceText:   body.PlaceText,
		IsAllday:    body.IsAllDay,
	}

	fmt.Println(body.GroupId)

	ctx := utils.ContextWithMetadata(c)
	resp, err := scheduleClient.Schedule.CreateSchedule(ctx, &schedulepb.CreateScheduleRequest{
		GroupId:       body.GroupId,
		Comment:       body.Comment,
		IsHighlighted: body.IsHighlighted,
		Link:          linkData,
	})

	if err != nil {
		utils.HandleGRPCError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": resp.Schedule,
	})

}

func UpdateSchedule(c *gin.Context) {
	scheduleClient := client.GetScheduleService(c)

	strId := c.Param("schedule-id")

	if strId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "filter parameter is required"})
		return
	}

	scheduleId, err := strconv.ParseInt(strId, 0, 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid schedule_id format", "details": err.Error()})
		return
	}

	var body struct {
		Title         string  `json:"title"`
		Description   string  `json:"description"`
		StartAt       string  `json:"start_at"`
		EndAt         string  `json:"end_at"`
		IsAllDay      bool    `json:"is_all_day"`
		PlaceId       *int64  `json:"place_id"`
		PlaceText     *string `json:"place_text"`
		Comment       *string `json:"comment"`
		IsHighlighted *bool   `json:"is_highlighted"`
	}

	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Required parameter", "details": err.Error()})
		return
	}

	linkData := &schedulepb.ScheduleLinkData{
		Title:       body.Title,
		Description: body.Description,
		StartAt:     body.StartAt,
		EndAt:       body.EndAt,
		PlaceId:     body.PlaceId,
		PlaceText:   body.PlaceText,
		IsAllday:    body.IsAllDay,
	}

	fmt.Println(linkData.StartAt)

	ctx := utils.ContextWithMetadata(c)
	resp, err := scheduleClient.Schedule.UpdateSchedule(ctx, &schedulepb.UpdateScheduleRequest{
		ScheduleId:    scheduleId,
		Comment:       body.Comment,
		IsHighlighted: body.IsHighlighted,
		Link:          linkData,
	})

	if err != nil {
		utils.HandleGRPCError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": resp.Schedule,
	})
}

func DeleteSchedule(c *gin.Context) {
	scheduleClient := client.GetScheduleService(c)

	strId := c.Param("schedule-id")

	if strId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "filter parameter is required"})
		return
	}

	scheduleId, err := strconv.ParseInt(strId, 0, 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid schedule_id format", "details": err.Error()})
		return
	}

	ctx := utils.ContextWithMetadata(c)
	resp, err := scheduleClient.Schedule.DeleteSchedule(ctx, &schedulepb.DeleteScheduleRequest{
		ScheduleId: scheduleId,
	})

	if err != nil {
		utils.HandleGRPCError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": resp.Success,
	})

}
