package usecases

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
)

type MemberUserCase struct {
	r MemberRepository
}

func NewMemberUserCase(repo MemberRepository) *MemberUserCase {
	return &MemberUserCase{r: repo}
}

func (m MemberUserCase) AddNewMemberToGroup(ctx context.Context, group string) (Member, error) {
	id, err := m.r.SaveNewMemberToGroup(ctx, group)
	if err != nil {
		return Member{}, fmt.Errorf("could not add new member: %w", err)
	}
	log.Info().Str("group", group).Str("id", id).Msg("created new")
	return Member{
		Group: group,
		ID:    id,
	}, nil
}

func (m MemberUserCase) GetHealthCheckFromMember(ctx context.Context, member Member) error {
	log.Info().Str("group", member.Group).Str("id", member.ID).Msg("got the ping")

	return nil
}

func (m MemberUserCase) RemoveMember(ctx context.Context, member Member) error {
	return m.r.DeleteMemberFrom(ctx, member)
}
