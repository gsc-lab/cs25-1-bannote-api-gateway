package utils

import (
	"context"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"

	"github.com/gsc-lab/cs25-1-bannote-api-gateway/constants"
	"github.com/gsc-lab/cs25-1-bannote-api-gateway/middleware"
)

// ContextWithMetadata gin.Context에서 user 정보를 추출해 gRPC context에 메타데이터로 추가
func ContextWithMetadata(c *gin.Context) context.Context {
	userCode := middleware.GetUserCode(c)
	userRole := middleware.GetUserRole(c)

	// 둘 다 없으면 그냥 기본 context 반환
	if userCode == "" && userRole == "" {
		return context.Background()
	}

	// 메타데이터 생성
	md := metadata.New(map[string]string{})

	if userCode != "" {
		md.Set(constants.MetadataUserCode, userCode)
	}

	if userRole != "" {
		md.Set(constants.MetadataUserRole, userRole)
	}

	return metadata.NewOutgoingContext(context.Background(), md)
}
