package repository

import (
	"context"
	"slices"
	"time"

	"github.com/denizgursoy/clerk/pkg/v1/usecases"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type MemberETCDRepository struct {
	members map[string][]usecases.Member
}

func NewMemberETCDRepository() *MemberETCDRepository {
	return &MemberETCDRepository{
		members: make(map[string][]usecases.Member),
	}
}

func (m MemberETCDRepository) SaveNewMemberToGroup(ctx context.Context, group string) (string, error) {
	uuid := uuid.New().String()
	if len(m.members[group]) == 0 {
		m.members[group] = make([]usecases.Member, 0)
	}
	member := usecases.Member{
		Group:           group,
		ID:              uuid,
		LastUpdatedTime: nil,
	}
	m.members[group] = append(m.members[group], member)

	return uuid, nil
}

func (m MemberETCDRepository) DeleteMemberFrom(ctx context.Context, member usecases.Member) error {
	members, ok := m.members[member.Group]
	if !ok {
		return usecases.ErrGroupNotFound
	}

	for i, m := range members {
		if m.ID == member.ID {
			members = slices.Delete(members, i, i+1)

			return nil
		}
	}

	return usecases.ErrMemberNotFound
}

func (m MemberETCDRepository) SaveLastUpdatedTime(ctx context.Context, member usecases.Member) error {
	members := m.members[member.Group]
	for i := range members {
		if members[i].ID == member.ID {
			now := time.Now()
			members[i].LastUpdatedTime = &now

			return nil
		}
	}

	return usecases.ErrMemberNotFound
}

func (m MemberETCDRepository) RemoveAllMemberNotAvailableDuringDuration(ctx context.Context,
	duration time.Duration) error {
	for group, members := range m.members {
		m.members[group] = slices.DeleteFunc(members, func(member usecases.Member) bool {
			if member.LastUpdatedTime != nil {
				if member.LastUpdatedTime.Add(duration).Before(time.Now()) {
					log.Info().Str("id", member.ID).Msg("deleting")

					return true
				}
			}

			return false
		})
	}

	return nil
}
