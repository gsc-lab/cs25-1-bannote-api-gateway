package client

import (
	"fmt"

	healthv1 "github.com/gsc-lab/cs25-1-bannote-api-gateway/gen/go/common-service/healthcheck"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ScheduleServiceClient struct {
	conn *grpc.ClientConn

	Health healthv1.HealthClient
}

func NewScheduleServiceClient(address string) (*ScheduleServiceClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to schedule-service: %w", err)
	}

	client := &ScheduleServiceClient{
		conn:   conn,
		Health: healthv1.NewHealthClient(conn),
	}

	fmt.Printf("✅ Connected to schedule-service at %s (all domains)\n", address)
	return client, nil
}

func (c *ScheduleServiceClient) Close() error {
	return c.conn.Close()
}
