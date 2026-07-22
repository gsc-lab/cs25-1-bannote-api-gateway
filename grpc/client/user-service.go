package client

import (
	"fmt"

	healthv1 "github.com/gsc-lab/cs25-1-bannote-api-gateway/gen/go/common-service/healthcheck"
	alloweddomainpb "github.com/gsc-lab/cs25-1-bannote-api-gateway/gen/go/user-service/alloweddomain"
	departmentpb "github.com/gsc-lab/cs25-1-bannote-api-gateway/gen/go/user-service/department"
	rolepb "github.com/gsc-lab/cs25-1-bannote-api-gateway/gen/go/user-service/role"
	studentclasspb "github.com/gsc-lab/cs25-1-bannote-api-gateway/gen/go/user-service/studentclass"
	userpb "github.com/gsc-lab/cs25-1-bannote-api-gateway/gen/go/user-service/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserServiceClient struct {
	conn *grpc.ClientConn

	User          userpb.UserServiceClient
	Department    departmentpb.DepartmentServiceClient
	Role          rolepb.RoleServiceClient
	AllowedDomain alloweddomainpb.AllowedDomainServiceClient
	StudentClass  studentclasspb.StudentClassServiceClient
	Health        healthv1.HealthClient
}

func NewUserServiceClient(address string) (*UserServiceClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user-service: %w", err)
	}

	client := &UserServiceClient{
		conn:          conn,
		User:          userpb.NewUserServiceClient(conn),
		Department:    departmentpb.NewDepartmentServiceClient(conn),
		Role:          rolepb.NewRoleServiceClient(conn),
		AllowedDomain: alloweddomainpb.NewAllowedDomainServiceClient(conn),
		StudentClass:  studentclasspb.NewStudentClassServiceClient(conn),
		Health:        healthv1.NewHealthClient(conn),
	}

	fmt.Printf("✅ Connected to user-service at %s (all domains)\n", address)
	return client, nil
}

func (c *UserServiceClient) Close() error {
	return c.conn.Close()
}
