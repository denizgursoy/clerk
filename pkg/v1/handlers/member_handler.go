package handlers

import (
	"context"

	"github.com/denizgursoy/clerk/pkg/v1/usecases"
	"github.com/denizgursoy/clerk/proto"
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

func (m MemberGRPCHandler) AddMember(ctx context.Context, request *proto.MemberRequest) (*proto.MemberResponse, error) {
	group, err := m.memberUseCase.AddNewMemberToGroup(ctx, request.GetGroup())
	if err != nil {
		return nil, err
	}

	return toResponse(group), nil
}

func (m MemberGRPCHandler) Ping(ctx context.Context, request *proto.PingRequest) (*proto.PingResponse, error) {
	// TODO implement me
	panic("implement me")
}

func toResponse(m usecases.Member) *proto.MemberResponse {
	return &proto.MemberResponse{
		Id:    m.ID,
		Group: m.Group,
	}
}
