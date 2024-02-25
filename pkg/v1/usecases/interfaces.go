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
	SaveNewMember(ctx context.Context, member Member) error
	DeleteMemberByID(ctx context.Context, id string) error
	SaveLastUpdatedTimeByID(ctx context.Context, id string, updateTime time.Time) error
	GetPartitionOfTheMemberByID(ctx context.Context, id string) (Partition, error)
	GetAllMembers(ctx context.Context) ([]*Member, error)
	UpdatePartitions(ctx context.Context, idPartitionMap map[string]Partition) error
}
