package client

import (
	"context"
	"fmt"
	"time"

	"github.com/denizgursoy/clerk/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ClerkClient struct {
	config     ClerkServerConfig
	grpcClient proto.MemberServiceClient
}

func NewClerkClient(config ClerkServerConfig) (*ClerkClient, error) {
	conn, err := grpc.Dial(config.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("could not start grpcClient: %w", err)
	}
	c := &ClerkClient{config: config}
	c.grpcClient = proto.NewMemberServiceClient(conn)

	return c, nil
}

func (c *ClerkClient) AddMember(ctx context.Context, a string, fn Cre) (*proto.MemberResponse, error) {
	request := proto.MemberRequest{
		Group: a,
	}
	member, err := c.grpcClient.AddMember(ctx, &request)
	if err != nil {
		return nil, err
	}

	m := convert(member)
	for {
		select {
		case <-ctx.Done():
			return nil, err
		case <-time.Tick(time.Duration(c.config.KeepAliveDurationInSeconds) * time.Second):
			c.keepAlive(ctx, m)
		}
	}
}

func (c *ClerkClient) keepAlive(ctx context.Context, member Member) error {

	pingRequest := &proto.PingRequest{
		Group: member.Group,
		Id:    member.ID,
	}
	_, err := c.grpcClient.Ping(ctx, pingRequest)

	return err
}

func convert(res *proto.MemberResponse) Member {
	return Member{
		Group: res.Group,
		ID:    res.Id,
	}
}
