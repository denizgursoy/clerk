package client

import (
	"context"
	"fmt"

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

func (c *ClerkClient) AddMember(ctx context.Context, a string) (*proto.MemberResponse, error) {
	request := proto.MemberRequest{
		Group: a,
	}
	return c.grpcClient.AddMember(ctx, &request)

}
