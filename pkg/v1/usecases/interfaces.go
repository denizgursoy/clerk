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
	GetPartitionOfTheMember(ctx context.Context, member Member) (Partition, error)
}

type MemberRepository interface {
	SaveNewMemberToGroup(ctx context.Context, group string) (Member, error)
	DeleteMemberFrom(ctx context.Context, member Member) error
	SaveLastUpdatedTime(ctx context.Context, member Member) error
	RemoveAllMemberNotAvailableDuringDuration(ctx context.Context, seconds time.Duration) error
	GetCurrentPartitionOfTheMember(ctx context.Context, member Member) (Partition, error)
	SetPartitionOfTheMember(ctx context.Context, member Member, p Partition) error
}
