package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gsc-lab/cs25-1-bannote-api-gateway/config"
	common "github.com/gsc-lab/cs25-1-bannote-api-gateway/gen/go/user-service/common"
	userpb "github.com/gsc-lab/cs25-1-bannote-api-gateway/gen/go/user-service/user"
	"github.com/gsc-lab/cs25-1-bannote-api-gateway/grpc/client"
	"github.com/gsc-lab/cs25-1-bannote-api-gateway/utils"
	"google.golang.org/api/idtoken"
)

type UserResponse struct {
	ID               string    `json:"id"`
	UserEmail        string    `json:"user_email"`
	UserName         string    `json:"user_name"`
	FamilyName       string    `json:"family_name"`
	GivenName        string    `json:"given_name"`
	UserType         string    `json:"user_type"`
	UserStatus       string    `json:"user_status"`
	Bio              string    `json:"bio"`
	ProfileImageUrl  string    `json:"profile_image_url"`
	CreatedAt        time.Time `json:"created_at"`
	DeletedAt        time.Time `json:"deleted_at,omitempty"`
	UserRoles        []string  `json:"user_roles"`
	StudentClassCode *string   `json:"student_class_code,omitempty"`
	StudentClassName *string   `json:"student_class_name,omitempty"`
	DepartmentCode   *string   `json:"department_code,omitempty"`
	DepartmentName   *string   `json:"department_name,omitempty"`
}

func convertUser(u *userpb.UserDetail) UserResponse {
	var createdAt, deletedAt time.Time
	if u.CreatedAt != nil {
		createdAt = u.CreatedAt.AsTime()
	}
	if u.DeletedAt != nil {
		deletedAt = u.DeletedAt.AsTime()
	}

	roles := make([]string, len(u.UserRoles))
	for i, r := range u.UserRoles {
		roles[i] = utils.StringFromUserRole(r)
	}

	return UserResponse{
		ID:               u.UserCode,
		UserEmail:        u.UserEmail,
		UserName:         u.FamilyName + " " + u.GivenName,
		FamilyName:       u.FamilyName,
		GivenName:        u.GivenName,
		UserType:         utils.StringFromUserType(u.UserType),
		UserStatus:       utils.StringFromUserStatus(u.UserStatus),
		Bio:              u.Bio,
		ProfileImageUrl:  u.ProfileImageUrl,
		CreatedAt:        createdAt,
		DeletedAt:        deletedAt,
		UserRoles:        roles,
		StudentClassCode: u.StudentClassCode,
		StudentClassName: u.StudentClassName,
		DepartmentCode:   u.DepartmentCode,
		DepartmentName:   u.DepartmentName,
	}
}

func convertUsers(users []*userpb.UserDetail) []UserResponse {
	result := make([]UserResponse, len(users))
	for i, u := range users {
		result[i] = convertUser(u)
	}
	return result
}

type ListUsersRequest struct {
	Page             int32   `form:"page"`
	Size             int32   `form:"size"`
	UserType         *string `form:"user_type,omitempty"`
	UserStatus       *string `form:"user_status,omitempty"`
	DepartmentCode   *string `form:"department_code,omitempty"`
	StudentClassCode *string `form:"student_class_code,omitempty"`
}

func ListUsers(c *gin.Context) {
	var request ListUsersRequest
	userClient := client.GetUserService(c)

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	ctx := utils.ContextWithMetadata(c)

	var userStatus *common.UserStatus
	if request.UserStatus != nil {
		userStatus = utils.ParseUserStatus(*request.UserStatus)
		if userStatus == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_status value"})
			return
		}
	}

	var userType *common.UserType
	if request.UserType != nil {
		userType = utils.ParseUserType(*request.UserType)
		if userType == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_type value"})
			return
		}
	}

	resp, err := userClient.User.ListUsers(ctx, &userpb.ListUsersRequest{
		Page:             request.Page,
		Size:             request.Size,
		Status:           userStatus,
		Type:             userType,
		DepartmentCode:   request.DepartmentCode,
		StudentClassCode: request.StudentClassCode,
	})

	if err != nil {
		utils.HandleGRPCError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  convertUsers(resp.Users),
		"total": resp.TotalCount,
		"page":  resp.Page,
		"size":  resp.Size,
	})
}

type createUserRequest struct {
	UserCode         string  `json:"user_code"`
	UserEmail        string  `json:"user_email"`
	FamilyName       string  `json:"family_name"`
	GivenName        string  `json:"given_name"`
	UserType         string  `json:"user_type"`
	ProfileImageUrl  string  `json:"profile_image_url"`
	StudentClassCode *string `json:"student_class_code"`
	DepartmentCode   *string `json:"department_code"`
	Credential       string  `json:"credential"`
}

func CreateUser(c *gin.Context) {
	var request createUserRequest
	userClient := client.GetUserService(c)

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	if request.Credential == "" {
		c.JSON(400, gin.H{"error": "Missing credential"})
		return
	}

	payload, err := idtoken.Validate(context.Background(), request.Credential, config.AppConfig.GoogleClientID)
	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid token", "details": err.Error()})
		return
	}

	email, _ := payload.Claims["email"].(string)

	if email != request.UserEmail {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Credential", "details": "Invalid Email"})
		return
	}

	var userType *common.UserType
	userType = utils.ParseUserType(request.UserType)
	if userType == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_type value"})
		return
	}

	resp, err := userClient.User.CreateUser(c, &userpb.CreateUserRequest{
		UserCode:         request.UserCode,
		UserEmail:        request.UserEmail,
		FamilyName:       request.FamilyName,
		GivenName:        request.GivenName,
		UserType:         *userType,
		ProfileImageUrl:  request.ProfileImageUrl,
		StudentClassCode: request.StudentClassCode,
		DepartmentCode:   request.DepartmentCode,
	})

	if err != nil {
		utils.HandleGRPCError(c, err)
		return
	}

	if !resp.Success {
		c.JSON(http.StatusBadRequest, gin.H{"error": resp.Reason})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success":   resp.Success,
		"can_login": resp.CanLogin,
		"data":      convertUser(resp.User),
	})
}

func UpdateUser(c *gin.Context) {
	userClient := client.GetUserService(c)
	userCode := c.Param("user_code")

	if userCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Required parameter"})
		return
	}

	var body struct {
		FamilyName      *string `json:"family_name"`
		GivenName       *string `json:"given_name"`
		ProfileImageUrl *string `json:"profile_image_url"`
		Bio             *string `json:"bio"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	ctx := utils.ContextWithMetadata(c)

	resp, err := userClient.User.UpdateUser(ctx, &userpb.UpdateUserRequest{
		UserCode:        userCode,
		FamilyName:      body.FamilyName,
		GivenName:       body.GivenName,
		ProfileImageUrl: body.ProfileImageUrl,
		Bio:             body.Bio,
	})

	if err != nil {
		utils.HandleGRPCError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": resp.User,
	})

}

func SearchUser(c *gin.Context) {
	userClient := client.GetUserService(c)

	var request struct {
		Name       string  `form:"q"`
		Page       int32   `form:"page"`
		Size       int32   `form:"size"`
		UserType   *string `form:"user_type"`
		UserStatus *string `form:"user_status"`
	}

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Required parameter", "details": err.Error()})
		return
	}

	var userStatus *common.UserStatus
	if request.UserStatus != nil {
		userStatus = utils.ParseUserStatus(*request.UserStatus)
		if userStatus == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_status value"})
			return
		}
	}

	var userType *common.UserType
	if request.UserType != nil {
		userType = utils.ParseUserType(*request.UserType)
		if userType == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_type value"})
			return
		}
	}

	ctx := utils.ContextWithMetadata(c)
	resp, err := userClient.User.SearchUsersByName(ctx, &userpb.SearchUsersByNameRequest{
		Name:   request.Name,
		Page:   request.Page,
		Size:   request.Size,
		Type:   userType,
		Status: userStatus,
	})

	if err != nil {
		utils.HandleGRPCError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  resp.Users,
		"total": resp.TotalCount,
	})
}

func GetUser(c *gin.Context) {
	userClient := client.GetUserService(c)
	userCode := c.Param("user_code")

	if userCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Required parameter"})
		return
	}

	ctx := utils.ContextWithMetadata(c)

	resp, err := userClient.User.GetUser(ctx, &userpb.GetUserRequest{
		UserCode: userCode,
	})

	if err != nil {
		utils.HandleGRPCError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": convertUser(resp.User),
	})

}
