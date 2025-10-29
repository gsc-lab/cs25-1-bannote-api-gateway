package client

import (
	"github.com/gin-gonic/gin"
)

const (
	UserServiceKey  = "userService"
	TokenServiceKey = "tokenService"
)

// Container holds all gRPC clients
type Container struct {
	UserService  *UserServiceClient
	TokenService *TokenServiceClient
}

// NewContainer creates a new client container
func NewContainer(userServiceAddr string, tokenServiceAddr string) (*Container, error) {
	userClient, err := NewUserServiceClient(userServiceAddr)
	tokenClient, err := NewTokenServiceClient(tokenServiceAddr)
	if err != nil {
		return nil, err
	}

	return &Container{
		UserService:  userClient,
		TokenService: tokenClient,
	}, nil
}

// Close closes all gRPC connections
func (c *Container) Close() error {
	return c.UserService.Close()
}

// Middleware injects the container into gin context
func (c *Container) Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(UserServiceKey, c.UserService)
		ctx.Set(TokenServiceKey, c.TokenService)
		ctx.Next()
	}
}

// GetUserService retrieves UserServiceClient from gin context
func GetUserService(c *gin.Context) *UserServiceClient {
	return c.MustGet(UserServiceKey).(*UserServiceClient)
}

func GetTokenService(c *gin.Context) *TokenServiceClient {
	return c.MustGet(TokenServiceKey).(*TokenServiceClient)
}
