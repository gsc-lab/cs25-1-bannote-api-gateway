package client

import (
	"errors"

	"github.com/gin-gonic/gin"
)

const (
	UserServiceKey      = "userService"
	TokenServiceKey     = "tokenService"
	ScheduleServiceKey  = "scheduleService"
	StudyroomServiceKey = "studyroomService"
)

// Container holds all gRPC clients
type Container struct {
	UserService      *UserServiceClient
	TokenService     *TokenServiceClient
	ScheduleService  *ScheduleServiceClient
	StudyroomService *StudyroomServiceClient
}

// NewContainer creates a new client container
func NewContainer(userServiceAddr string, tokenServiceAddr string, scheduleServiceAddr string, studyroomServiceAddr string) (*Container, error) {
	userClient, err := NewUserServiceClient(userServiceAddr)
	tokenClient, err := NewTokenServiceClient(tokenServiceAddr)
	scheduleClient, err := NewScheduleServiceClient(scheduleServiceAddr)
	studyroomClient, err := NewStudyroomServiceClient(studyroomServiceAddr)
	if err != nil {
		return nil, err
	}

	return &Container{
		UserService:      userClient,
		TokenService:     tokenClient,
		ScheduleService:  scheduleClient,
		StudyroomService: studyroomClient,
	}, nil
}

// Close closes all gRPC connections
func (c *Container) Close() error {
	var errs []error

	if err := c.UserService.Close(); err != nil {
		errs = append(errs, err)
	}
	if err := c.TokenService.Close(); err != nil {
		errs = append(errs, err)
	}
	if err := c.ScheduleService.Close(); err != nil {
		errs = append(errs, err)
	}
	if err := c.StudyroomService.Close(); err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

// Middleware injects the container into gin context
func (c *Container) Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(UserServiceKey, c.UserService)
		ctx.Set(TokenServiceKey, c.TokenService)
		ctx.Set(ScheduleServiceKey, c.ScheduleService)
		ctx.Set(StudyroomServiceKey, c.StudyroomService)
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

func GetScheduleService(c *gin.Context) *ScheduleServiceClient {
	return c.MustGet(ScheduleServiceKey).(*ScheduleServiceClient)
}

func GetStudyroomService(c *gin.Context) *StudyroomServiceClient {
	return c.MustGet(StudyroomServiceKey).(*StudyroomServiceClient)
}
