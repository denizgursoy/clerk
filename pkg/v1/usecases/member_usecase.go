package usecases

import (
	"context"
	"fmt"
)

type MemberUserCase struct {
	repo MemberRepository
}

func NewMemberUserCase(repo MemberRepository) *MemberUserCase {
	return &MemberUserCase{repo: repo}
}

func (m MemberUserCase) AddNewMemberToGroup(ctx context.Context, group string) (Member, error) {
	id, err := m.repo.SaveNewMemberToGroup(ctx, group)
	if err != nil {
		return Member{}, fmt.Errorf("could not add new member: %w", err)
	}

	return Member{
		Group: group,
		ID:    id,
	}, nil
}
