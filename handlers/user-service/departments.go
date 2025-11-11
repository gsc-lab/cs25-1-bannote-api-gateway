package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	departmentpb "github.com/gsc-lab/cs25-1-bannote-api-gateway/gen/go/user-service/department"
	"github.com/gsc-lab/cs25-1-bannote-api-gateway/grpc/client"
	"github.com/gsc-lab/cs25-1-bannote-api-gateway/utils"
)

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
		"department": resp.Department,
	})
}

type CreateDepartmentRequest struct {
	DepartmentCode string `json:"department_code"`
	DepartmentName string `json:"department_name"`
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

	c.JSON(http.StatusCreated, resp)
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
		"departments": resp.Departments,
		"total_count": resp.TotalCount,
		"page":        resp.Page,
		"size":        resp.Size,
	})
}

type UpdateDepartmentRequest struct {
	DepartmentName string `json:"department_name"`
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

	c.JSON(http.StatusOK, resp)
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
