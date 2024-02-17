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
	// TODO implement me
	panic("implement me")
}

func (m MemberGRPCHandler) Ping(ctx context.Context, request *proto.PingRequest) (*proto.PingResponse, error) {
	// TODO implement me
	panic("implement me")
}
