package usecases

import (
	"time"
)

type Member struct {
	Group           string
	ID              string
	LastUpdatedTime *time.Time
	CreatedAt       time.Time
	Partition
}

func (m Member) IsActive(duration time.Duration) bool {
	dateToUse := m.CreatedAt
	if m.LastUpdatedTime != nil {
		dateToUse = *m.LastUpdatedTime
	}

	return dateToUse.Add(duration).After(time.Now())
}

type Partition struct {
	Ordinal int
	Total   int
}

type MemberGroup struct {
	allMembers []Member
}

func NewMemberGroup() *MemberGroup {
	return &MemberGroup{
		allMembers: make([]Member, 0),
	}
}

func (m *MemberGroup) Group() string {
	return m.allMembers[0].Group
}

func (m *MemberGroup) Partition() string {
	return m.allMembers[0].Group
}

func (m *MemberGroup) Add(member Member) {
	m.allMembers = append(m.allMembers, member)
}

func (m *MemberGroup) UnstableMembers(lifeTime time.Duration) []Member {
	members := make([]Member, 0)
	for i := range m.allMembers {
		if !m.allMembers[i].IsActive(lifeTime) {
			members = append(members, m.allMembers[i])
		}
	}

	return members
}

func (m *MemberGroup) StableMembers(lifeTime time.Duration) []Member {
	members := make([]Member, 0)
	for i := range m.allMembers {
		if m.allMembers[i].IsActive(lifeTime) {
			members = append(members, m.allMembers[i])
		}
	}

	return members
}

func (m *MemberGroup) RearrangeOrders() {

}
