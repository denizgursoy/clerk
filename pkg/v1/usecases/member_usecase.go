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
	if err := m.r.RemoveAllMemberNotAvailableDuringDuration(ctx, m.c.LifeSpanDurationInSeconds); err != nil {
		return fmt.Errorf("could not remove all dead members: %w", err)
	}

	return nil
}

func (m MemberUserCase) GetPartitionOfTheMember(ctx context.Context, member Member) (Partition, error) {
	return m.r.GetCurrentPartitionOfTheMember(ctx, member)
}
