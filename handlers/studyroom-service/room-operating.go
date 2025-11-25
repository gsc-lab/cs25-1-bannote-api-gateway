package studyroom_service

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	roomoperatingpb "github.com/gsc-lab/cs25-1-bannote-api-gateway/gen/go/studyroom-service/room_operating_hour"
	"github.com/gsc-lab/cs25-1-bannote-api-gateway/grpc/client"
	"github.com/gsc-lab/cs25-1-bannote-api-gateway/utils"
)

func GetRoomOperating(c *gin.Context) {
	studyroomClient := client.GetStudyroomService(c)

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
	resp, err := studyroomClient.RoomOperating.GetRoomOperatingHours(ctx, &roomoperatingpb.GetRoomOperatingHoursRequest{
		RoomId: id,
	})

	if err != nil {
		utils.HandleGRPCError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": resp.RoomOperatingHours,
		"id":   id,
	})
}
