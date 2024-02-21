//go:generate mockgen -source=interfaces.go -destination=interfaces_mock.go -package=usecases
package usecases

import (
	"context"
	"time"
)

type Cache interface {
	AddNewInstance() (int, error)
}

type MemberUseCase interface {
	AddNewMemberToGroup(ctx context.Context, group string) (Member, error)
	GetHealthCheckFromMember(ctx context.Context, member Member) error
	RemoveMember(ctx context.Context, member Member) error
	TriggerBalance()
	StopBalance()
}

type MemberRepository interface {
	SaveNewMemberToGroup(ctx context.Context, group string) (string, error)
	DeleteMemberFrom(ctx context.Context, member Member) error
	SaveLastUpdatedTime(ctx context.Context, member Member) error
	RemoveAllMemberNotAvailableDuringDuration(ctx context.Context, seconds time.Duration) error
}
