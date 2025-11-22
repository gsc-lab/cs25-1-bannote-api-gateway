package studyroom_service

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	roompb "github.com/gsc-lab/cs25-1-bannote-api-gateway/gen/go/studyroom-service/room"
	"github.com/gsc-lab/cs25-1-bannote-api-gateway/grpc/client"
	"github.com/gsc-lab/cs25-1-bannote-api-gateway/utils"
)

type listRoomRequest struct {
	Page int32 `form:"page"`
	Size int32 `form:"size"`
}

func ListRoom(c *gin.Context) {
	studyroomClient := client.GetStudyroomService(c)

	var request *listRoomRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Required parameter", "details": err.Error()})
		return
	}

	ctx := utils.ContextWithMetadata(c)
	resp, err := studyroomClient.Room.ListRooms(ctx, &roompb.ListRoomsRequest{
		Page: request.Page,
		Size: request.Size,
	})

	if err != nil {
		utils.HandleGRPCError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  resp.Rooms,
		"total": resp.GetTotalCount(),
		"page":  resp.Page,
		"size":  resp.Size,
	})
}

type createRoomRequest struct {
	DepartmentCode string `json:"department_code"`
	Name           string `json:"name"`
	MaximumMember  int32  `json:"maximum_member"`
}

func CreateRoom(c *gin.Context) {
	studyroomClient := client.GetStudyroomService(c)

	var request *createRoomRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Required parameter", "details": err.Error()})
		return
	}

	ctx := utils.ContextWithMetadata(c)
	resp, err := studyroomClient.Room.CreateRoom(ctx, &roompb.CreateRoomRequest{
		Name:           request.Name,
		DepartmentCode: request.DepartmentCode,
		MaximumMember:  request.MaximumMember,
	})

	if err != nil {
		utils.HandleGRPCError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": resp.Room,
		"id":   resp.Room.Id,
	})
}

func DeleteRoom(c *gin.Context) {
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
	_, err = studyroomClient.Room.DeleteRoom(ctx, &roompb.DeleteRoomRequest{
		Id: id,
	})

	if err != nil {
		utils.HandleGRPCError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Room deleted successfully",
	})
}
