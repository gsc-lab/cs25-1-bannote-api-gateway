package client

import (
	"fmt"

	healthv1 "github.com/gsc-lab/cs25-1-bannote-api-gateway/gen/go/common-service/healthcheck"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type StudyroomServiceClient struct {
	conn *grpc.ClientConn

	Health healthv1.HealthClient
}

func NewStudyroomServiceClient(address string) (*StudyroomServiceClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to studyroom-service: %w", err)
	}

	client := &StudyroomServiceClient{
		conn:   conn,
		Health: healthv1.NewHealthClient(conn),
	}

	fmt.Printf("✅ Connected to studyroom-service at %s (all domains)\n", address)
	return client, nil
}

func (c *StudyroomServiceClient) Close() error {
	return c.conn.Close()
}
