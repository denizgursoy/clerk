package client

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/denizgursoy/clerk/proto"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type ClerkClient struct {
	config     ClerkServerConfig
	grpcClient proto.MemberServiceClient
	member     Member
	fn         NotifyFunction
	partition  proto.Partition
	cancelFunc context.CancelFunc
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

func (c *ClerkClient) Start(parentContext context.Context, group string, fn NotifyFunction) error {
	if fn == nil {
		return ErrEmptyFunction
	}

	if len(strings.TrimSpace(group)) == 0 {
		return ErrEmptyGroup
	}

	ctx, cancel := context.WithCancel(parentContext)
	c.cancelFunc = cancel
	member, err := c.grpcClient.AddMember(ctx, &proto.MemberRequest{Group: group})
	if err != nil {
		return err
	}
	c.fn = fn
	c.member = convert(member)
	go c.executeFunction(ctx)
	c.statPinging(ctx)

	return nil
}
func (c *ClerkClient) Remove() error {
	_, err := c.grpcClient.RemoveMember(context.Background(), toProto(c.member))
	if err != nil {
		return err
	}
	if err = c.terminate(); err != nil {
		return err
	}

	return nil
}

func (c *ClerkClient) terminate() error {
	c.cancelFunc()

	return nil
}

func (c *ClerkClient) statPinging(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.Tick(time.Duration(c.config.KeepAliveDurationInSeconds) * time.Second):
			_, err := c.grpcClient.Ping(ctx, toProto(c.member))
			if err != nil {
				if status, ok := status.FromError(err); ok {
					// Check gRPC status code and message
					log.Printf("gRPC status code: %d, message: %s", status.Code(), status.Message())
				}
				break
			}
		}
	}
}

func (c *ClerkClient) executeFunction(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.Tick(time.Duration(c.config.KeepAliveDurationInSeconds) * time.Second):
			partition, err := c.grpcClient.Listen(ctx, toProto(c.member))
			if err != nil {
				if status, ok := status.FromError(err); ok {
					// Check gRPC status code and message
					log.Printf("gRPC status code: %d, message: %s", status.Code(), status.Message())
				}

				break
			}

			if partition.Ordinal != c.partition.Ordinal ||
				partition.Total != c.partition.Total {

				c.partition.Ordinal = partition.Ordinal
				c.partition.Total = partition.Total

				c.fn(ctx, c.partition.Ordinal, c.partition.Total)
			}
		}
	}

}

func convert(res *proto.Member) Member {
	return Member{
		Group: res.Group,
		ID:    res.Id,
	}
}

func toProto(m Member) *proto.Member {
	return &proto.Member{
		Group: m.Group,
		Id:    m.ID,
	}
}
