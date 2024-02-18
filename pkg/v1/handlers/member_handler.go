package handlers

import (
	"context"
	"math/rand"
	"time"

	"github.com/denizgursoy/clerk/pkg/v1/usecases"
	"github.com/denizgursoy/clerk/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type MemberGRPCHandler struct {
	memberUseCase usecases.MemberUseCase
	proto.UnimplementedMemberServiceServer
}

func NewMemberGRPCHandler(grpcServer *grpc.Server, memberUseCase usecases.MemberUseCase) *MemberGRPCHandler {
	handler := &MemberGRPCHandler{memberUseCase: memberUseCase}
	proto.RegisterMemberServiceServer(grpcServer, handler)

	return handler
}

func (m MemberGRPCHandler) AddMember(ctx context.Context, request *proto.MemberRequest) (*proto.Member, error) {
	group, err := m.memberUseCase.AddNewMemberToGroup(ctx, request.GetGroup())
	if err != nil {
		return new(proto.Member), err
	}

	return toResponse(group), nil
}

func (m MemberGRPCHandler) Ping(ctx context.Context, member *proto.Member) (*empty.Empty, error) {
	err := m.memberUseCase.GetHealthCheckFromMember(ctx, toMember(member))

	return new(empty.Empty), err
}

func (m MemberGRPCHandler) RemoveMember(ctx context.Context, member *proto.Member) (*empty.Empty, error) {
	err := m.memberUseCase.RemoveMember(ctx, toMember(member))

	return new(empty.Empty), err
}

func (m MemberGRPCHandler) Listen(_ *proto.Member, l proto.MemberService_ListenServer) error {
	for {
		select {
		case <-time.Tick(10 * time.Second):
			log.Info().Msg("sending new partition")
			err := l.Send(&proto.Partition{
				Ordinal: rand.Int63(),
				Total:   rand.Int63(),
			})
			if err != nil {
				log.Err(err).Msg("error")
			}
		}
	}
}

func toResponse(m usecases.Member) *proto.Member {
	return &proto.Member{
		Id:    m.ID,
		Group: m.Group,
	}
}

func toMember(p *proto.Member) usecases.Member {
	return usecases.Member{
		Group: p.Group,
		ID:    p.Id,
	}
}
