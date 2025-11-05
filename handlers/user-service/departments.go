package handlers

import (
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

func ListDepartments(c *gin.Context) {
	userClient := client.GetUserService(c)

	ctx := utils.ContextWithMetadata(c)
	resp, err := userClient.Department.ListDepartments(ctx, &departmentpb.ListDepartmentsRequest{
		Page: 1,
		Size: 2,
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
