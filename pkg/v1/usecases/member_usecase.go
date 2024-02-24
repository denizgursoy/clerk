package usecases

import (
	"context"
	"fmt"
	"slices"
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
		stableMembers, unstableMembers := groups[i].StableAndUnstableMembers(m.c.LifeSpanDuration)
		log.Info().
			Str("group", groups[i].Group()).
			Int("unstable member count", len(unstableMembers)).
			Int("stable member count", len(stableMembers)).
			Msg("setting ordinal again")
		m.ClearOrdinals(ctx, unstableMembers)
		m.setNewOrdinals(ctx, stableMembers)
	}

	return nil
}

func (m MemberUserCase) setNewOrdinals(ctx context.Context, stableMembers []*Member) {
	allCount := len(stableMembers)
	usedOrdinals := getCurrentUsedOrdinals(stableMembers)
	idPartitionMap := make(map[string]Partition)
	for i := range stableMembers {
		if !isValidOrdinal(stableMembers[i].Ordinal, allCount) {
			ordinal := getFirstUnusedOrdinal(usedOrdinals, allCount)
			stableMembers[i].Ordinal = ordinal
			usedOrdinals = append(usedOrdinals, ordinal)
		}
		stableMembers[i].Total = allCount
		idPartitionMap[stableMembers[i].ID] = stableMembers[i].Partition
	}

	if err := m.r.UpdatePartitions(ctx, idPartitionMap); err != nil {
		log.Err(err).Msg("could not update the partitions")
	}
	log.Info().Msg("update partition successfully")
}

func isValidOrdinal(ordinal, allCount int) bool {
	if ordinal == 0 || ordinal > allCount {
		return false
	}

	return true
}

func getFirstUnusedOrdinal(ordinals []int, count int) int {
	for i := 1; i <= count; i++ {
		if !slices.Contains(ordinals, i) {
			return i
		}
	}

	return 0
}

func getCurrentUsedOrdinals(stableMembers []*Member) []int {
	usedOrdinals := make([]int, 0)
	for i := range stableMembers {
		ordinal := stableMembers[i].Ordinal
		if ordinal != 0 {
			usedOrdinals = append(usedOrdinals, ordinal)
		}
	}

	return usedOrdinals
}

func (m MemberUserCase) GetPartitionOfTheMember(ctx context.Context, member Member) (Partition, error) {
	return m.r.GetCurrentPartitionOfTheMember(ctx, member)
}

func (m MemberUserCase) ClearOrdinals(ctx context.Context, unStableMembers []*Member) error {
	idPartitionMap := make(map[string]Partition)
	for i := range unStableMembers {
		idPartitionMap[unStableMembers[i].ID] = DefaultPartition
	}
	if err := m.r.UpdatePartitions(ctx, idPartitionMap); err != nil {
		log.Err(err).Msg("could not update unstable partitions")

		return fmt.Errorf("could not update partiion of unstable: %w", err)
	}

	return nil
}

func ConvertToMembersToGroups(members []*Member) []*MemberGroup {
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
