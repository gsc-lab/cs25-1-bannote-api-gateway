package client

import (
	"fmt"

	healthv1 "github.com/gsc-lab/cs25-1-bannote-api-gateway/gen/go/common-service/healthcheck"
	grouppb "github.com/gsc-lab/cs25-1-bannote-api-gateway/gen/go/schedule-service/group"
	grouptagpb "github.com/gsc-lab/cs25-1-bannote-api-gateway/gen/go/schedule-service/group_tag"
	schedulepb "github.com/gsc-lab/cs25-1-bannote-api-gateway/gen/go/schedule-service/schedule"
	schedulelinkpb "github.com/gsc-lab/cs25-1-bannote-api-gateway/gen/go/schedule-service/schedule_link"
	tagpb "github.com/gsc-lab/cs25-1-bannote-api-gateway/gen/go/schedule-service/tag"
	usergrouppb "github.com/gsc-lab/cs25-1-bannote-api-gateway/gen/go/schedule-service/user_group"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ScheduleServiceClient struct {
	conn         *grpc.ClientConn
	Tag          tagpb.TagServiceClient
	Group        grouppb.GroupServiceClient
	GroupTag     grouptagpb.GroupTagServiceClient
	UserGroup    usergrouppb.UserGroupServiceClient
	Schedule     schedulepb.ScheduleServiceClient
	ScheduleLink schedulelinkpb.ScheduleLinkServiceClient
	Health       healthv1.HealthClient
}

func NewScheduleServiceClient(address string) (*ScheduleServiceClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to schedule-service: %w", err)
	}

	client := &ScheduleServiceClient{
		conn:         conn,
		Tag:          tagpb.NewTagServiceClient(conn),
		Group:        grouppb.NewGroupServiceClient(conn),
		GroupTag:     grouptagpb.NewGroupTagServiceClient(conn),
		UserGroup:    usergrouppb.NewUserGroupServiceClient(conn),
		Schedule:     schedulepb.NewScheduleServiceClient(conn),
		ScheduleLink: schedulelinkpb.NewScheduleLinkServiceClient(conn),
		Health:       healthv1.NewHealthClient(conn),
	}

	fmt.Printf("✅ Connected to schedule-service at %s (all domains)\n", address)
	return client, nil
}

func (c *ScheduleServiceClient) Close() error {
	return c.conn.Close()
}
