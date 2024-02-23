package repository

import (
	"context"
	"slices"
	"time"

	"github.com/denizgursoy/clerk/pkg/v1/usecases"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

var members = make([]usecases.Member, 0)

type MemberETCDRepository struct {
}

func NewMemberETCDRepository() *MemberETCDRepository {
	return &MemberETCDRepository{}
}

func (m *MemberETCDRepository) SaveNewMemberToGroup(ctx context.Context, group string) (usecases.Member, error) {
	uuid := uuid.New().String()

	member := usecases.Member{
		Group:           group,
		ID:              uuid,
		LastUpdatedTime: nil,
		CreatedAt:       time.Now(),
	}
	members = append(members, member)

	return member, nil
}

func (m *MemberETCDRepository) DeleteMemberFrom(ctx context.Context, member usecases.Member) error {

	for i, mem := range members {
		if mem.ID == member.ID {
			members = slices.Delete(members, i, i+1)

			return nil
		}
	}

	return usecases.ErrMemberNotFound
}

func (m *MemberETCDRepository) SaveLastUpdatedTime(ctx context.Context, member usecases.Member) error {
	findMember, err := m.findMember(member)
	if err != nil {
		return err
	}

	now := time.Now()
	findMember.LastUpdatedTime = &now

	return nil
}

func (m *MemberETCDRepository) RemoveAllMemberNotAvailableDuringDuration(ctx context.Context,
	duration time.Duration) error {

	members = slices.DeleteFunc(members, func(member usecases.Member) bool {
		if member.LastUpdatedTime != nil {
			if member.LastUpdatedTime.Add(duration).Before(time.Now()) {
				log.Info().Str("id", member.ID).Msg("deleting")

				return true
			}
		}

		return false
	})

	return nil
}

func (m *MemberETCDRepository) GetCurrentPartitionOfTheMember(ctx context.Context,
	member usecases.Member) (usecases.Partition, error) {
	findMember, err := m.findMember(member)
	if err != nil {
		return usecases.Partition{}, err
	}

	return findMember.Partition, nil
}

// func (m MemberETCDRepository) SetPartitionOfTheMember(ctx context.Context,
// 	member usecases.Member, p usecases.Partition) error {
// 	findMember, err := m.findMember(member)
// 	if err != nil {
// 		return err
// 	}
// 	findMember.Partition = p
//
// 	return nil
// }

func (m *MemberETCDRepository) findMember(member usecases.Member) (usecases.Member, error) {
	for i := range members {
		if members[i].ID == member.ID {
			return members[i], nil
		}
	}

	return usecases.Member{}, usecases.ErrMemberNotFound
}

func (m *MemberETCDRepository) GetAllMembers(ctx context.Context) ([]usecases.Member, error) {
	return members, nil
}

func (m *MemberETCDRepository) DeleteMembers(ctx context.Context, members []usecases.Member) error {
	// TODO implement me
	// panic("implement me")
	return nil
}

func (m *MemberETCDRepository) UpdatePartitions(ctx context.Context, idPartitionMap map[string]usecases.Partition) error {
	for i := range members {
		partition, ok := idPartitionMap[members[i].ID]
		if ok {
			members[i].Partition = partition
		}
	}

	return nil
}
