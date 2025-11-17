package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	alloweddomainpb "github.com/gsc-lab/cs25-1-bannote-api-gateway/gen/go/user-service/alloweddomain"
	"github.com/gsc-lab/cs25-1-bannote-api-gateway/grpc/client"
	"github.com/gsc-lab/cs25-1-bannote-api-gateway/utils"
)

// convertAllowedDomain converts a single AllowedDomain to React Admin compatible format
func convertAllowedDomain(domain *alloweddomainpb.AllowedDomain) map[string]interface{} {
	return map[string]interface{}{
		"id":         domain.GetDomain(), // Use domain as unique identifier
		"domain":     domain.GetDomain(),
		"created_by": domain.GetCreatedBy(),
		"created_at": domain.GetCreatedAt(),
	}
}

// convertAllowedDomains converts a list of AllowedDomains to React Admin compatible format
func convertAllowedDomains(domains []*alloweddomainpb.AllowedDomain) []map[string]interface{} {
	result := make([]map[string]interface{}, len(domains))
	for i, domain := range domains {
		result[i] = convertAllowedDomain(domain)
	}
	return result
}

type createAllowedDomainRequest struct {
	Domain string `json:"domain"`
}

func CreateAllowedDomain(c *gin.Context) {
	userClient := client.GetUserService(c)

	var request *createAllowedDomainRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Required parameter", "details": err.Error()})
		return
	}

	ctx := utils.ContextWithMetadata(c)
	resp, err := userClient.AllowedDomain.AddAllowedDomain(ctx, &alloweddomainpb.AddAllowedDomainRequest{
		Domain: request.Domain,
	})

	if err != nil {
		utils.HandleGRPCError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": convertAllowedDomain(resp.GetAllowedDomain()),
		"id":   resp.AllowedDomain.Domain,
	})
}

func DeleteAllowedDomain(c *gin.Context) {
	userClient := client.GetUserService(c)

	domain := c.Param("id")
	if domain == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Domain parameter required"})
		return
	}
	// Remove leading slash from wildcard parameter
	if len(domain) > 0 && domain[0] == '/' {
		domain = domain[1:]
	}

	ctx := utils.ContextWithMetadata(c)
	resp, err := userClient.AllowedDomain.RemoveAllowedDomain(ctx, &alloweddomainpb.RemoveAllowedDomainRequest{
		Domain: domain,
	})

	if err != nil {
		utils.HandleGRPCError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": convertAllowedDomain(resp.GetAllowedDomain()),
	})
}

type listAllowedDomainRequest struct {
	Page int32 `form:"page"`
	Size int32 `form:"size"`
}

func ListAllowedDomain(c *gin.Context) {
	userClient := client.GetUserService(c)

	var request *listAllowedDomainRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Required parameter", "details": err.Error()})
		return
	}

	ctx := utils.ContextWithMetadata(c)
	resp, err := userClient.AllowedDomain.ListAllowedDomain(ctx, &alloweddomainpb.ListAllowedDomainRequest{
		Page: request.Page,
		Size: request.Size,
	})

	if err != nil {
		utils.HandleGRPCError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  convertAllowedDomains(resp.GetAllowedDomain()),
		"total": resp.GetTotalCount(),
		"page":  resp.Page,
		"size":  resp.Size,
	})
}
