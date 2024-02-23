package handlers

import (
	"context"

	"github.com/denizgursoy/clerk/pkg/v1/usecases"
	"github.com/denizgursoy/clerk/proto"
	"github.com/golang/protobuf/ptypes/empty"
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
	if err != nil {
		return new(empty.Empty), err
	}

	return new(empty.Empty), nil
}

func (m MemberGRPCHandler) RemoveMember(ctx context.Context, member *proto.Member) (*empty.Empty, error) {
	err := m.memberUseCase.RemoveMember(ctx, toMember(member))
	if err != nil {
		return new(empty.Empty), err
	}

	return new(empty.Empty), nil
}

func (m MemberGRPCHandler) QueryPartition(ctx context.Context, member *proto.Member) (*proto.Partition, error) {
	partition, err := m.memberUseCase.GetPartitionOfTheMember(ctx, toMember(member))
	if err != nil {
		return new(proto.Partition), err
	}

	return toProtoPartition(partition), nil
}

func toProtoPartition(p usecases.Partition) *proto.Partition {
	return &proto.Partition{
		Ordinal: int32(p.Ordinal),
		Total:   int32(p.Total),
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
