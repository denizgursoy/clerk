package handlers

import (
	"github.com/denizgursoy/clerk/pkg/v1/usecases"
	"github.com/denizgursoy/clerk/proto"
)

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
