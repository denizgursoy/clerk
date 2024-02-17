package usecases

type MemberUserCase struct {
	repo MemberRepository
}

func NewMemberUserCase(repo MemberRepository) *MemberUserCase {
	return &MemberUserCase{repo: repo}
}
