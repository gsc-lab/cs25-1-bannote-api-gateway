package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/gsc-lab/cs25-1-bannote-api-gateway/config"
	"github.com/gsc-lab/cs25-1-bannote-api-gateway/constants"
)

// CustomClaims JWT claims 구조
type CustomClaims struct {
	UserID string `json:"user_id"`
	Roles  string `json:"roles"`
	jwt.RegisteredClaims
}

// GRPCMetadata gRPC 메타데이터 세팅 (JWT 로컬 검증)
func GRPCMetadata() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := extractToken(c)

		// 토큰이 없는 경우
		if tokenString == "" {
			c.JSON(401, gin.H{
				"error":   "Unauthorized",
				"message": "로그인이 필요합니다.",
			})
			c.Abort()
			return
		}

		// JWT 검증
		token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return config.AppConfig.PublicKey, nil
		})

		// 토큰이 있는데 유효하지 않으면 401 (변조, 만료 등)
		if err != nil || !token.Valid {
			c.JSON(401, gin.H{
				"error":   "Unauthorized",
				"message": "유효하지 않은 인증 토큰입니다",
			})
			c.Abort()
			return
		}

		// claims에서 user 정보 추출
		if claims, ok := token.Claims.(*CustomClaims); ok {
			if claims.UserID != "" {
				c.Set(constants.MetadataUserCode, claims.UserID)
			}
			if claims.Roles != "" {
				c.Set(constants.MetadataUserRole, claims.Roles)
			}
		}

		c.Next()
	}
}

// extractToken 쿠키 또는 Header의 Authorization에서 토큰을 가져옴
func extractToken(c *gin.Context) string {
	if token, err := c.Cookie("access_token"); err == nil && token != "" {
		return token
	}

	authHeader := c.GetHeader("Authorization")
	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}

	return ""
}

func GetUserCode(c *gin.Context) string {
	if userCode, exists := c.Get(constants.MetadataUserCode); exists {
		if code, ok := userCode.(string); ok {
			return code
		}
	}
	return ""
}

func GetUserRole(c *gin.Context) string {
	if userRole, exists := c.Get(constants.MetadataUserRole); exists {
		if role, ok := userRole.(string); ok {
			return role
		}
	}
	return ""
}
