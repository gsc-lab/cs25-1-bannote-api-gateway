package studyroom_service

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	roomexceptionpb "github.com/gsc-lab/cs25-1-bannote-api-gateway/gen/go/studyroom-service/room_exception"
	"github.com/gsc-lab/cs25-1-bannote-api-gateway/grpc/client"
	"github.com/gsc-lab/cs25-1-bannote-api-gateway/utils"
)

func GetRoomException(c *gin.Context) {
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

	fromDate := c.Query("from_date")

	ctx := utils.ContextWithMetadata(c)
	resp, err := studyroomClient.RoomException.GetRoomExceptions(ctx, &roomexceptionpb.GetRoomExceptionsRequest{
		RoomId:   id,
		FromDate: fromDate,
	})

	if err != nil {
		utils.HandleGRPCError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": resp.RoomExceptions,
		"id":   id,
	})
}

func UpdateRoomException(c *gin.Context) {
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

	var body struct {
		Exceptions []*roomexceptionpb.RoomExceptionUpdateItem `json:"exceptions"`
		FromDate   string                                     `json:"from_data"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	ctx := utils.ContextWithMetadata(c)
	resp, err := studyroomClient.RoomException.UpdateRoomExceptions(ctx, &roomexceptionpb.UpdateRoomExceptionsRequest{
		RoomId:     id,
		Exceptions: body.Exceptions,
		//FromData:   body.FromDate,
	})

	if err != nil {
		utils.HandleGRPCError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Exception updated successfully",
		"data":    resp,
	})
}
