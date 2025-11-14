package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	common "github.com/gsc-lab/cs25-1-bannote-api-gateway/gen/go/user-service/common"
	studentclasspb "github.com/gsc-lab/cs25-1-bannote-api-gateway/gen/go/user-service/studentclass"
	"github.com/gsc-lab/cs25-1-bannote-api-gateway/grpc/client"
	"github.com/gsc-lab/cs25-1-bannote-api-gateway/utils"
)

// StudentClassResponse is a custom response structure with id instead of student_class_code
type StudentClassResponse struct {
	ID             string    `json:"id"`
	DepartmentCode string    `json:"department_code"`
	Name           string    `json:"name"`
	AdmissionYear  int32     `json:"admission_year"`
	GraduationYear int32     `json:"graduation_year"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	DeletedAt      time.Time `json:"deleted_at,omitempty"`
	CreatedBy      string    `json:"created_by"`
}

// convertStudentClass converts protobuf StudentClass to API response
func convertStudentClass(sc *studentclasspb.StudentClass) StudentClassResponse {
	var createdAt, deletedAt time.Time
	if sc.CreatedAt != nil {
		createdAt = sc.CreatedAt.AsTime()
	}
	if sc.DeletedAt != nil {
		deletedAt = sc.DeletedAt.AsTime()
	}

	status := utils.StringFromStudentClassStatus(sc.Status)

	return StudentClassResponse{
		ID:             sc.StudentClassCode,
		DepartmentCode: sc.DepartmentCode,
		Name:           sc.Name,
		AdmissionYear:  sc.AdmissionYear,
		GraduationYear: sc.GraduationYear,
		Status:         status,
		CreatedAt:      createdAt,
		DeletedAt:      deletedAt,
		CreatedBy:      sc.CreatedBy,
	}
}

// convertStudentClasses converts multiple protobuf StudentClasses to API responses
func convertStudentClasses(classes []*studentclasspb.StudentClass) []StudentClassResponse {
	result := make([]StudentClassResponse, len(classes))
	for i, sc := range classes {
		result[i] = convertStudentClass(sc)
	}
	return result
}

func GetStudentClass(c *gin.Context) {
	userClient := client.GetUserService(c)
	code := c.Param("id")

	ctx := utils.ContextWithMetadata(c)
	resp, err := userClient.StudentClass.GetStudentClass(ctx, &studentclasspb.GetStudentClassRequest{
		StudentClassCode: code,
	})

	if err != nil {
		utils.HandleGRPCError(c, err)
		return
	}

	c.JSON(http.StatusOK, convertStudentClass(resp.StudentClass))
}

func GetManyStudentClasses(c *gin.Context) {
	userClient := client.GetUserService(c)
	filterStr := c.Query("filter")

	if filterStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "filter parameter is required"})
		return
	}

	var filter struct {
		ID []string `json:"id"` // ID 배열로 파싱
	}

	if err := json.Unmarshal([]byte(filterStr), &filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid filter format"})
		return
	}

	resp, err := userClient.StudentClass.GetManyStudentClasses(c, &studentclasspb.GetManyStudentClassesRequest{
		StudentClassesCode: filter.ID,
	})

	if err != nil {
		utils.HandleGRPCError(c, err)
		return
	}

	c.JSON(http.StatusOK, convertStudentClasses(resp.StudentClasses))
}

type CreateStudentClassRequest struct {
	DepartmentCode   string `json:"department_code"`
	StudentClassCode string `json:"id"`
	StudentClassName string `json:"name"`
	AdmissionYear    int32  `json:"admission_year"`
	GraduationYear   int32  `json:"graduation_year"`
}

func CreateStudentClass(c *gin.Context) {
	var request CreateStudentClassRequest
	userClient := client.GetUserService(c)

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Required parameter", "details": err.Error()})
		return
	}

	ctx := utils.ContextWithMetadata(c)
	resp, err := userClient.StudentClass.CreateStudentClass(ctx, &studentclasspb.CreateStudentClassRequest{
		DepartmentCode:   request.DepartmentCode,
		StudentClassCode: request.StudentClassCode,
		StudentClassName: request.StudentClassName,
		AdmissionYear:    request.AdmissionYear,
		GraduationYear:   request.GraduationYear,
	})

	if err != nil {
		utils.HandleGRPCError(c, err)
		return
	}

	c.JSON(http.StatusCreated, convertStudentClass(resp.StudentClass))
}

type UpdateStudentClassRequest struct {
	DepartmentCode   *string `json:"department_code,omitempty"`
	StudentClassName *string `json:"name,omitempty"`
	AdmissionYear    *int32  `json:"admission_year,omitempty"`
	GraduationYear   *int32  `json:"graduation_year,omitempty"`
	Status           *string `json:"status,omitempty"`
}

func UpdateStudentClass(c *gin.Context) {
	var request UpdateStudentClassRequest
	userClient := client.GetUserService(c)
	code := c.Param("id")

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	ctx := utils.ContextWithMetadata(c)

	// 문자열 status를 enum으로 변환
	var status *common.StudentClassStatus
	if request.Status != nil {
		status = utils.ParseStudentClassStatus(*request.Status)
		if status == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status value. Use 'active' or 'graduated'"})
			return
		}
	}

	resp, err := userClient.StudentClass.UpdateStudentClass(ctx, &studentclasspb.UpdateStudentClassRequest{
		StudentClassCode: code,
		DepartmentCode:   request.DepartmentCode,
		StudentClassName: request.StudentClassName,
		AdmissionYear:    request.AdmissionYear,
		GraduationYear:   request.GraduationYear,
		Status:           status,
	})

	if err != nil {
		utils.HandleGRPCError(c, err)
		return
	}

	c.JSON(http.StatusOK, convertStudentClass(resp.StudentClass))
}

func DeleteStudentClass(c *gin.Context) {
	userClient := client.GetUserService(c)
	code := c.Param("id")

	ctx := utils.ContextWithMetadata(c)
	resp, err := userClient.StudentClass.DeleteStudentClass(ctx, &studentclasspb.DeleteStudentClassRequest{
		StudentClassCode: code,
	})

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete Student Class", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

type ListStudentClassRequest struct {
	Page               int32   `form:"page"`
	Size               int32   `form:"size"`
	DepartmentCode     *string `form:"department_code,omitempty"`
	StudentClassStatus *string `form:"status,omitempty"`
}

func ListStudentClass(c *gin.Context) {
	var request ListStudentClassRequest
	userClient := client.GetUserService(c)

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	ctx := utils.ContextWithMetadata(c)

	// 문자열 status를 enum으로 변환
	var status *common.StudentClassStatus
	if request.StudentClassStatus != nil {
		status = utils.ParseStudentClassStatus(*request.StudentClassStatus)
		if status == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status value. Use 'active' or 'graduated'"})
			return
		}
	}

	fmt.Println(status)

	resp, err := userClient.StudentClass.ListStudentClasses(ctx, &studentclasspb.ListStudentClassesRequest{
		Page:           request.Page,
		Size:           request.Size,
		DepartmentCode: request.DepartmentCode,
		Status:         status,
	})

	if err != nil {
		utils.HandleGRPCError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  convertStudentClasses(resp.StudentClasses),
		"total": resp.TotalCount,
		"page":  resp.Page,
		"size":  resp.Size,
	})
}
