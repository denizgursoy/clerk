package repository

import (
	"context"
	"slices"

	"github.com/denizgursoy/clerk/pkg/v1/usecases"
	"github.com/google/uuid"
)

type MemberETCDRepository struct {
	members map[string][]string
}

func NewMemberETCDRepository() *MemberETCDRepository {
	return &MemberETCDRepository{
		members: make(map[string][]string),
	}
}

func (m MemberETCDRepository) SaveNewMemberToGroup(ctx context.Context, group string) (string, error) {
	uuid := uuid.New().String()
	if len(m.members[group]) == 0 {
		m.members[group] = make([]string, 0)
	}

	m.members[group] = append(m.members[group], uuid)

	return uuid, nil
}

func (m MemberETCDRepository) DeleteMemberFrom(ctx context.Context, member usecases.Member) error {
	members, ok := m.members[member.Group]
	if !ok {
		return usecases.ErrGroupNotFound
	}

	if slices.Contains(members, member.ID) {
		m.members[member.Group] = slices.DeleteFunc(members, func(s string) bool {
			return s == member.ID
		})

		return nil
	}

	return usecases.ErrMemberNotFound
}
