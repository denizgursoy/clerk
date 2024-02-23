package client

import (
	"context"
	"time"

	"github.com/denizgursoy/clerk/proto"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/status"
)

type (
	NotifyFunction func(ctx context.Context, ordinal, total int) error
	MemberConfig   struct {
		KeepAliveDuration time.Duration
	}
	ClerkServerConfig struct {
		Address string
	}

	Member struct {
		group      string
		id         string
		fn         NotifyFunction
		partition  proto.Partition
		cancelFunc context.CancelFunc
		config     MemberConfig
		grpcClient proto.MemberServiceClient
	}
)

var (
	defaultMemberConfig = MemberConfig{
		KeepAliveDuration: 2 * time.Second,
	}
)

func newMember(grpcClient proto.MemberServiceClient, member *proto.Member, c MemberConfig) *Member {
	return &Member{
		group:      member.Group,
		id:         member.Id,
		grpcClient: grpcClient,
		config:     c,
	}
}

// Start function initializes the pinging.
// It is a blocking function
func (m *Member) Start(c context.Context, fn NotifyFunction) error {
	if fn == nil {
		return ErrEmptyFunction
	}
	ctx, cancelFunc := context.WithCancel(c)
	m.cancelFunc = cancelFunc
	m.fn = fn

	go m.executeFunction(ctx)
	m.statPinging(ctx)

	return nil
}

func (m *Member) statPinging(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.Tick(m.config.KeepAliveDuration):
			_, err := m.grpcClient.Ping(ctx, toProto(m))
			if err != nil {
				if errStatus, ok := status.FromError(err); ok {
					// Check gRPC status code and message
					log.Printf("gRPC status code: %d, message: %s", errStatus.Code(), errStatus.Message())
				}
				break
			}
		}
	}
}

func (m *Member) executeFunction(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.Tick(m.config.KeepAliveDuration):
			partition, err := m.grpcClient.QueryPartition(ctx, toProto(m))
			if err != nil {
				if errorStatus, ok := status.FromError(err); ok {
					// Check gRPC status code and message
					log.Printf("gRPC status code: %d, message: %s", errorStatus.Code(), errorStatus.Message())
				}

				break
			}

			if partition.Ordinal != m.partition.Ordinal ||
				partition.Total != m.partition.Total {

				m.partition.Ordinal = partition.Ordinal
				m.partition.Total = partition.Total

				m.fn(ctx, int(m.partition.Ordinal), int(m.partition.Total))
			}
		}
	}

}

func (m *Member) terminate() error {
	m.cancelFunc()
	return nil
}

func (m *Member) Remove() error {
	_, err := m.grpcClient.RemoveMember(context.Background(), toProto(m))
	if err != nil {
		return err
	}
	if err = m.terminate(); err != nil {
		return err
	}

	return nil
}

func toProto(m *Member) *proto.Member {
	return &proto.Member{
		Group: m.group,
		Id:    m.id,
	}
}
