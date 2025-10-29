package client

import (
	"fmt"

	tokenpb "github.com/gsc-lab/cs25-1-bannote-api-gateway/gen/go/token-service/token"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TokenServiceClient struct {
	conn *grpc.ClientConn

	Token tokenpb.TokenServiceClient
}

func NewTokenServiceClient(address string) (*TokenServiceClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to token-service: %w", err)
	}

	client := &TokenServiceClient{
		conn:  conn,
		Token: tokenpb.NewTokenServiceClient(conn),
	}

	fmt.Printf("✅ Connected to token-service at %s (all domains)\n", address)
	return client, nil
}

func (c *TokenServiceClient) Close() error {
	return c.conn.Close()
}
