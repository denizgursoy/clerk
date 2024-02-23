package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/denizgursoy/clerk/pkg/v1/config"
	"github.com/rs/zerolog/log"
)

type MemberUserCase struct {
	r MemberRepository
	c config.Config
	t *time.Ticker
}

func NewMemberUserCase(r MemberRepository, c config.Config) *MemberUserCase {
	return &MemberUserCase{r: r, c: c, t: time.NewTicker(c.CheckDuration)}
}

func (m MemberUserCase) AddNewMemberToGroup(ctx context.Context, group string) (Member, error) {
	member, err := m.r.SaveNewMemberToGroup(ctx, group)
	if err != nil {
		return Member{}, fmt.Errorf("could not add new member: %w", err)
	}
	log.Info().Str("group", group).Str("id", member.ID).Msg("created new")

	return member, nil
}

func (m MemberUserCase) GetHealthCheckFromMember(ctx context.Context, member Member) error {
	log.Info().Str("group", member.Group).Str("id", member.ID).Msg("got the ping")

	return m.r.SaveLastUpdatedTime(ctx, member)
}

func (m MemberUserCase) RemoveMember(ctx context.Context, member Member) error {
	return m.r.DeleteMemberFrom(ctx, member)
}

func (m MemberUserCase) TriggerBalance() {
	defer m.t.Stop()

	for range m.t.C {
		if err := m.balance(context.Background()); err != nil {
			log.Err(err).Msg("balance error")
		}
	}

}

func (m MemberUserCase) StopBalance() {
	m.t.Stop()
}

func (m MemberUserCase) balance(ctx context.Context) error {
	log.Info().Msg("checking balance")

	members, err := m.r.GetAllMembers(ctx)
	if err != nil {
		return fmt.Errorf("could not fetch all allMembers: %w", err)
	}

	groups := ConvertToMembersToGroups(members)
	for i := range groups {
		if unstableMembers := groups[i].UnstableMembers(m.c.LifeSpanDurationInSeconds); len(unstableMembers) > 0 {
			log.Info().Str("group", groups[i].Group()).Msg("is not stable")
			if err := m.r.DeleteMembers(ctx, unstableMembers); err != nil {
				return fmt.Errorf("could not remove all dead allMembers: %w", err)
			}
			groups[i].RearrangeOrders()
		}
	}

	return nil
}

func (m MemberUserCase) GetPartitionOfTheMember(ctx context.Context, member Member) (Partition, error) {
	return m.r.GetCurrentPartitionOfTheMember(ctx, member)
}

func ConvertToMembersToGroups(members []Member) []*MemberGroup {
	allGroups := make(map[string]*MemberGroup)
	for i := range members {
		groupName := members[i].Group
		group, ok := allGroups[groupName]
		if !ok {
			memberGroup := NewMemberGroup()
			allGroups[groupName] = memberGroup
			group = memberGroup
		}
		group.Add(members[i])
	}

	groups := make([]*MemberGroup, 0)
	for _, group := range allGroups {
		groups = append(groups, group)
	}

	return groups
}
