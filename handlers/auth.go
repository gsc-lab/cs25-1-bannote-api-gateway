package handlers

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	userpb "github.com/gsc-lab/cs25-1-bannote-api-gateway/gen/go/user-service/user"
	"github.com/gsc-lab/cs25-1-bannote-api-gateway/grpc/client"
	"google.golang.org/api/idtoken"

	tokenpb "github.com/gsc-lab/cs25-1-bannote-api-gateway/gen/go/token-service/token"
)

const googleClientID = "616869349585-5kkpb0gvqlaq3kl4l67n6tq1nt6fee32.apps.googleusercontent.com"
const frontURL = "http://localhost:5173"

func GoogleCallback(c *gin.Context) {
	userClient := client.GetUserService(c)
	credential := c.Query("credential")

	if credential == "" {
		c.JSON(400, gin.H{"error": "Missing credential"})
		return
	}

	payload, err := idtoken.Validate(context.Background(), credential, googleClientID)
	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid token", "details": err.Error()})
		return
	}

	email, _ := payload.Claims["email"].(string)

	// 이메일로 가입된 유저인지 확인
	resp, err := userClient.User.UserLogin(context.Background(), &userpb.UserLoginRequest{
		Email: email,
	})

	if err != nil {
		fmt.Printf("❌ gRPC Error: %v\n", err)
		c.JSON(500, gin.H{"error": "Failed to check registration", "details": err.Error()})
		return
	}

	if !resp.Exists {
		c.Redirect(302, frontURL+"/register")
		return
	}

	if resp.CanLogin == false {
		c.JSON(401, gin.H{"error": "Invalid token", "details": err.Error()})
		return
	}

	tokenClient := client.GetTokenService(c)

	token, err := tokenClient.Token.GenerateAccessToken(context.Background(), &tokenpb.GenerateAccessTokenRequest{
		UserId: resp.User.GetUserCode(),
		Roles:  resp.User.GetGivenName(),
	})

	if err != nil {
		c.JSON(401, gin.H{"error": "Token error", "details": err.Error()})
		return
	}

	var userResponse gin.H
	if resp.User != nil {
		userResponse = gin.H{
			"user_code":         resp.User.UserCode,
			"user_email":        resp.User.UserEmail,
			"family_name":       resp.User.FamilyName,
			"given_name":        resp.User.GivenName,
			"user_type":         resp.User.UserType.String(),
			"user_status":       resp.User.UserStatus.String(),
			"bio":               resp.User.Bio,
			"profile_image_url": resp.User.ProfileImageUrl,
			"created_at":        resp.User.CreatedAt,
		}
	}

	c.JSON(200, gin.H{
		"message":   "User already exists",
		"exists":    resp.Exists,
		"can_login": resp.CanLogin,
		"user":      userResponse,
		"token":     token.GetAccessToken(),
	})

}
