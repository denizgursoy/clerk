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
)

type ClerkClient struct {
	config     ClerkServerConfig
	grpcClient proto.MemberServiceClient
	member     Member
	fn         NotifyFunction
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

	ctx, _ := context.WithCancel(parentContext)
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
func (c *ClerkClient) Remove(ctx context.Context) error {
	_, err := c.grpcClient.RemoveMember(ctx, toProto(c.member))

	return err
}

func (c *ClerkClient) statPinging(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.Tick(time.Duration(c.config.KeepAliveDurationInSeconds) * time.Second):
			_, _ = c.grpcClient.Ping(ctx, toProto(c.member))
		}
	}
}

func (c *ClerkClient) executeFunction(ctx context.Context) error {
	stream, err := c.grpcClient.Listen(ctx, toProto(c.member))
	if err != nil {
		return err
	}

	for {
		select {
		case <-stream.Context().Done():
			return nil
		default:
			recv, err := stream.Recv()
			log.Info().
				Int64("ordinal", recv.Ordinal).
				Int64("total", recv.Total).
				Msg("got new partitions")
			if err != nil {
				log.Err(err).Msg("could not get from stream")
			}
			if err = c.fn(ctx, recv.Ordinal, recv.Total); err != nil {
				log.Err(err).Msg("could not execute function")
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
