package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	departmentpb "github.com/gsc-lab/cs25-1-bannote-api-gateway/gen/go/user-service/department"
	"github.com/gsc-lab/cs25-1-bannote-api-gateway/grpc/client"
	"github.com/gsc-lab/cs25-1-bannote-api-gateway/utils"
)

// DepartmentResponse is a custom response structure with id instead of department_code
type DepartmentResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// convertDepartment converts protobuf Department to API response
func convertDepartment(dept *departmentpb.Department) DepartmentResponse {
	var createdAt time.Time
	if dept.CreatedAt != nil {
		createdAt = dept.CreatedAt.AsTime()
	}
	return DepartmentResponse{
		ID:        dept.DepartmentCode,
		Name:      dept.Name,
		CreatedAt: createdAt,
	}
}

// convertDepartments converts multiple protobuf Departments to API responses
func convertDepartments(depts []*departmentpb.Department) []DepartmentResponse {
	result := make([]DepartmentResponse, len(depts))
	for i, dept := range depts {
		result[i] = convertDepartment(dept)
	}
	return result
}

func GetDepartment(c *gin.Context) {
	userClient := client.GetUserService(c)
	code := c.Param("id")

	ctx := utils.ContextWithMetadata(c)
	resp, err := userClient.Department.GetDepartment(ctx, &departmentpb.GetDepartmentRequest{
		DepartmentCode: code,
	})

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get department", "details": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"department": convertDepartment(resp.Department),
	})
}

type CreateDepartmentRequest struct {
	DepartmentCode string `json:"id"`
	DepartmentName string `json:"name"`
}

func CreateDepartment(c *gin.Context) {
	var request CreateDepartmentRequest
	userClient := client.GetUserService(c)

	err := c.ShouldBind(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Required parameter", "details": err.Error()})
		return
	}

	ctx := utils.ContextWithMetadata(c)
	resp, err := userClient.Department.CreateDepartment(ctx, &departmentpb.CreateDepartmentRequest{
		DepartmentCode: request.DepartmentCode,
		DepartmentName: request.DepartmentName,
	})

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create departments", "details": err.Error()})
		return
	}

	c.JSON(
		http.StatusCreated,
		convertDepartment(resp.GetDepartment()),
	)
}

type ListDepartmentsRequest struct {
	Page int32 `form:"page"`
	Size int32 `form:"size"`
}

func ListDepartments(c *gin.Context) {
	var request ListDepartmentsRequest
	userClient := client.GetUserService(c)
	err := c.ShouldBind(&request)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Required parameter", "details": err.Error()})
		return
	}

	ctx := utils.ContextWithMetadata(c)
	resp, err := userClient.Department.ListDepartments(ctx, &departmentpb.ListDepartmentsRequest{
		Page: request.Page,
		Size: request.Size,
	})

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to list departments", "details": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"data":  convertDepartments(resp.Departments),
		"total": resp.TotalCount,
		"page":  resp.Page,
		"size":  resp.Size,
	})
}

type UpdateDepartmentRequest struct {
	DepartmentName string `json:"name"`
}

func UpdateDepartments(c *gin.Context) {
	userClient := client.GetUserService(c)
	code := c.Param("id")

	var request UpdateDepartmentRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Required parameter", "details": err.Error()})
		return
	}

	ctx := utils.ContextWithMetadata(c)
	resp, err := userClient.Department.UpdateDepartment(ctx, &departmentpb.UpdateDepartmentRequest{
		DepartmentCode: code,
		Name:           &request.DepartmentName,
	})

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update department", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, convertDepartment(resp.GetDepartment()))
}

func DeleteDepartments(c *gin.Context) {
	userClient := client.GetUserService(c)
	code := c.Param("id")

	ctx := utils.ContextWithMetadata(c)
	resp, err := userClient.Department.DeleteDepartment(ctx, &departmentpb.DeleteDepartmentRequest{
		DepartmentCode: code,
	})

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get department", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
