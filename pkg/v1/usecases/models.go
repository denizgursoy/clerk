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
	Ordinal int64
	Total   int64
}

type MemberGroup struct {
	allMembers      []Member
	unstableMembers []Member
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
	for i := range m.allMembers {
		if !m.allMembers[i].IsActive(lifeTime) {
			m.unstableMembers = append(m.unstableMembers, m.allMembers[i])
		}
	}

	return m.unstableMembers
}

func (m *MemberGroup) RearrangeOrders() {

}
