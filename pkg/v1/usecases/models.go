package usecases

import (
	"time"
)

var DefaultPartition = Partition{
	Ordinal: 0,
	Total:   0,
}

type Member struct {
	Group           string
	ID              string
	LastUpdatedTime *time.Time
	CreatedAt       time.Time
	Partition
}

func (m Member) IsActiveForTheLast(duration time.Duration) bool {
	return m.LastActiveDate().Add(duration).After(time.Now())
}
func (m Member) LastActiveDate() time.Time {
	if m.LastUpdatedTime != nil {
		return *m.LastUpdatedTime
	}
	return m.CreatedAt
}

type Partition struct {
	Ordinal int
	Total   int
}

type MemberGroup struct {
	allMembers []*Member
}

func NewMemberGroup() *MemberGroup {
	return &MemberGroup{
		allMembers: make([]*Member, 0),
	}
}

func (m *MemberGroup) Group() string {
	return m.allMembers[0].Group
}

func (m *MemberGroup) Partition() string {
	return m.allMembers[0].Group
}

func (m *MemberGroup) Add(member *Member) {
	m.allMembers = append(m.allMembers, member)
}

func (m *MemberGroup) StableAndUnstableMembers(d time.Duration) ([]*Member, []*Member) {
	stableMembers := make([]*Member, 0)
	unstableMembers := make([]*Member, 0)
	for i := range m.allMembers {
		if m.allMembers[i].IsActiveForTheLast(d) {
			stableMembers = append(stableMembers, m.allMembers[i])
		} else {
			unstableMembers = append(unstableMembers, m.allMembers[i])
		}
	}

	return stableMembers, unstableMembers
}

func (m *MemberGroup) RearrangeOrders() {

}
