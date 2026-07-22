package studyroom_service

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	reservationpb "github.com/gsc-lab/cs25-1-bannote-api-gateway/gen/go/studyroom-service/reservation"
	"github.com/gsc-lab/cs25-1-bannote-api-gateway/grpc/client"
	"github.com/gsc-lab/cs25-1-bannote-api-gateway/utils"
)

func GetManyReservations(c *gin.Context) {
	studyroomClient := client.GetStudyroomService(c)

	type queries struct {
		StartTime string `form:"start_time"`
		EndTime   string `form:"end_time"`
		RoomIds   string `form:"room_ids"` // 콤마로 구분된 문자열 (예: "2,1")
	}

	var request queries
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Required parameter", "details": err.Error()})
		return
	}

	// room_ids 문자열을 []int64로 변환
	var roomIds []int64
	if request.RoomIds != "" {
		roomIdStrs := strings.Split(request.RoomIds, ",")
		for _, idStr := range roomIdStrs {
			idStr = strings.TrimSpace(idStr)
			if idStr == "" {
				continue
			}
			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room_id format", "details": err.Error()})
				return
			}
			roomIds = append(roomIds, id)
		}
	}

	ctx := utils.ContextWithMetadata(c)
	resp, err := studyroomClient.Reservation.ListReservations(ctx, &reservationpb.ListReservationsRequest{
		RoomIds:        roomIds,
		StartTimeAfter: request.StartTime,
		EndTimeBefore:  request.EndTime,
	})

	if err != nil {
		utils.HandleGRPCError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  resp.Reservations,
		"total": len(resp.Reservations),
	})
}

func CreateReservation(c *gin.Context) {
	studyroomClient := client.GetStudyroomService(c)

	var body struct {
		RoomId    int64    `json:"room_id"`
		StartTime string   `json:"start_time"`
		EndTime   string   `json:"end_time"`
		Purpose   string   `json:"purpose"`
		UserCodes []string `json:"user_codes"`
	}

	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Required parameter", "details": err.Error()})
		return
	}

	ctx := utils.ContextWithMetadata(c)
	resp, err := studyroomClient.Reservation.CreateReservation(ctx, &reservationpb.CreateReservationRequest{
		RoomId:    body.RoomId,
		StartTime: body.StartTime,
		EndTime:   body.EndTime,
		Purpose:   body.Purpose,
		Priority:  1,             // TODO: 우선 순위 설정 필요
		UserCodes: []int64{1234}, // TODO: 프론트에서 받는 입력으로 변경
	})

	if err != nil {
		utils.HandleGRPCError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": resp.Reservation,
		"id":   resp.Reservation.RoomId,
	})
}

func UpdateReservation(c *gin.Context) {
	studyroomClient := client.GetStudyroomService(c)

	code := c.Param("code")

	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Required parameter"})
		return
	}

	var body struct {
		RoomId    int64  `json:"room_id"`
		StartTime string `json:"start_time"`
		EndTime   string `json:"end_time"`
		Purpose   string `json:"purpose"`
	}

	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Required parameter", "details": err.Error()})
		return
	}

	ctx := utils.ContextWithMetadata(c)

	resp, err := studyroomClient.Reservation.UpdateReservation(ctx, &reservationpb.UpdateReservationRequest{
		Code:      code,
		RoomId:    body.RoomId,
		StartTime: body.StartTime,
		EndTime:   body.EndTime,
		Purpose:   body.Purpose,
	})

	if err != nil {
		utils.HandleGRPCError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": resp.Reservation,
		"id":   resp.Reservation.Code,
	})
}

func DeleteReservation(c *gin.Context) {
	studyroomClient := client.GetStudyroomService(c)

	code := c.Param("code")

	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Required parameter"})
		return
	}

	ctx := utils.ContextWithMetadata(c)

	_, err := studyroomClient.Reservation.DeleteReservation(ctx, &reservationpb.DeleteReservationRequest{
		Code: code,
	})

	if err != nil {
		utils.HandleGRPCError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
