package services

import (
	"context"
	"time"

	healthv1 "github.com/gsc-lab/cs25-1-bannote-api-gateway/gen/go/common-service/healthcheck"
	"github.com/gsc-lab/cs25-1-bannote-api-gateway/grpc/client"
)

type ServiceStatus struct {
	Service string `json:"service"`
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

type HealthCheckResult struct {
	OverallStatus string          `json:"overall_status"`
	Timestamp     int64           `json:"timestamp"`
	Services      []ServiceStatus `json:"services"`
}

func CheckAllServices(container *client.Container) *HealthCheckResult {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result := &HealthCheckResult{
		Timestamp: time.Now().Unix(),
		Services:  make([]ServiceStatus, 0),
	}

	allHealthy := true

	// Check UserService
	userStatus := checkService(ctx, container.UserService.Health, "user-service")
	result.Services = append(result.Services, userStatus)
	if userStatus.Status != "SERVING" {
		allHealthy = false
	}

	// Check TokenService
	tokenStatus := checkService(ctx, container.TokenService.Health, "token-service")
	result.Services = append(result.Services, tokenStatus)
	if tokenStatus.Status != "SERVING" {
		allHealthy = false
	}

	if allHealthy {
		result.OverallStatus = "UP"
	} else {
		result.OverallStatus = "DEGRADED"
	}

	return result
}

func checkService(ctx context.Context, healthClient healthv1.HealthClient, serviceName string) ServiceStatus {
	req := &healthv1.HealthCheckRequest{
		Service: "",
	}

	resp, err := healthClient.Check(ctx, req)
	if err != nil {
		return ServiceStatus{
			Service: serviceName,
			Status:  "DOWN",
			Message: err.Error(),
		}
	}

	return ServiceStatus{
		Service: serviceName,
		Status:  resp.Status.String(),
		Message: "",
	}
}
