package utils

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

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

// HandleGRPCError gRPC 에러를 HTTP 응답으로 변환
func HandleGRPCError(c *gin.Context, err error) {
	st, ok := status.FromError(err)
	if !ok {
		// gRPC 에러가 아닌 경우
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal server error",
			"details": err.Error(),
		})
		return
	}

	// gRPC 상태 코드를 HTTP 상태 코드로 변환
	httpCode := grpcCodeToHTTP(st.Code())

	// 상세한 에러 응답 구성
	response := gin.H{
		"error":   st.Message(),
		"code":    st.Code().String(),
		"details": st.Details(),
	}

	c.JSON(httpCode, response)
}

// grpcCodeToHTTP gRPC 상태 코드를 HTTP 상태 코드로 변환
func grpcCodeToHTTP(code codes.Code) int {
	switch code {
	case codes.OK:
		return http.StatusOK
	case codes.Canceled:
		return http.StatusRequestTimeout
	case codes.Unknown:
		return http.StatusInternalServerError
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests
	case codes.FailedPrecondition:
		return http.StatusBadRequest
	case codes.Aborted:
		return http.StatusConflict
	case codes.OutOfRange:
		return http.StatusBadRequest
	case codes.Unimplemented:
		return http.StatusNotImplemented
	case codes.Internal:
		return http.StatusInternalServerError
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	case codes.DataLoss:
		return http.StatusInternalServerError
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}
