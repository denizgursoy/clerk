package repository

import (
	"context"

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
