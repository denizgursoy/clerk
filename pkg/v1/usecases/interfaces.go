//go:generate mockgen -source=interfaces.go -destination=interfaces_mock.go -package=usecases
package usecases

import "context"

type Cache interface {
	AddNewInstance() (int, error)
}

type MemberUseCase interface {
	AddNewMemberToGroup(ctx context.Context, group string) (Member, error)
}

type MemberRepository interface {
	SaveNewMemberToGroup(ctx context.Context, group string) (string, error)
}
