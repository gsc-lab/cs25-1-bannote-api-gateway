package handlers

import (
	"context"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gsc-lab/cs25-1-bannote-api-gateway/config"
	userpb "github.com/gsc-lab/cs25-1-bannote-api-gateway/gen/go/user-service/user"
	"github.com/gsc-lab/cs25-1-bannote-api-gateway/grpc/client"
	"google.golang.org/api/idtoken"

	tokenpb "github.com/gsc-lab/cs25-1-bannote-api-gateway/gen/go/token-service/token"
)

func GoogleCallback(c *gin.Context) {
	userClient := client.GetUserService(c)
	credential := c.Query("credential")

	if credential == "" {
		c.JSON(400, gin.H{"error": "Missing credential"})
		return
	}

	payload, err := idtoken.Validate(context.Background(), credential, config.AppConfig.GoogleClientID)
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
		c.JSON(200, gin.H{"data": resp})
		return
	}

	if !resp.CanLogin {
		c.JSON(401, gin.H{"error": "User is not allowed to login"})
		return
	}

	tokenClient := client.GetTokenService(c)

	// UserRole enum을 콤마로 구분된 문자열로 변환 (USER_ROLE_ 제거)
	roleStrings := make([]string, len(resp.User.GetUserRoles()))
	for i, role := range resp.User.GetUserRoles() {
		roleStrings[i] = strings.TrimPrefix(role.String(), "USER_ROLE_")
	}
	rolesStr := strings.Join(roleStrings, ",")

	token, err := tokenClient.Token.GenerateAccessToken(context.Background(), &tokenpb.GenerateAccessTokenRequest{
		UserId: resp.User.GetUserCode(),
		Roles:  rolesStr,
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

	// 토큰을 쿠키로 설정
	c.SetCookie(
		"access_token",                // 쿠키 이름
		token.GetAccessToken(),        // 토큰 값
		config.AppConfig.CookieMaxAge, // 만료 시간 (환경변수에서 설정)
		"/",                           // 경로
		"",                            // 도메인
		false,                         // Secure (HTTPS만 허용 여부, 개발 환경이므로 false)
		true,                          // HttpOnly (JavaScript 접근 차단)
	)

	c.JSON(200, gin.H{
		"message":   "User already exists",
		"exists":    resp.Exists,
		"can_login": resp.CanLogin,
		"user":      userResponse,
	})

}
