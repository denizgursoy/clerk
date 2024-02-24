package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/denizgursoy/clerk/pkg/v1/usecases"
	"github.com/google/uuid"
	"go.etcd.io/etcd/client/v3"
)

var members = make([]*usecases.Member, 0)

const (
	ETCDRecordPrefix = "clerk"
	Delimiter        = "-"
)

type MemberETCDRepository struct {
	e *clientv3.Client
}

func NewMemberETCDRepository(c *clientv3.Client) *MemberETCDRepository {
	return &MemberETCDRepository{
		e: c,
	}
}

func (m *MemberETCDRepository) SaveNewMemberToGroup(ctx context.Context, group string) (usecases.Member, error) {
	member := usecases.Member{
		Group:           group,
		ID:              getETCDPrefix(group) + uuid.New().String(),
		LastUpdatedTime: nil,
		CreatedAt:       time.Now(),
	}
	marshal, err := json.Marshal(member)
	if err != nil {
		return usecases.Member{}, fmt.Errorf("could not marshal member: %w", err)
	}
	_, err = m.e.Put(ctx, member.ID, string(marshal))
	if err != nil {
		return usecases.Member{}, fmt.Errorf("could not save member to etcd: %w", err)
	}

	return member, nil
}

func (m *MemberETCDRepository) DeleteMemberFrom(ctx context.Context, member usecases.Member) error {
	response, err := m.e.Delete(ctx, member.ID)
	if err != nil {
		return fmt.Errorf("could not delete member: %w", err)
	}
	if response.Deleted == 0 {
		return usecases.ErrMemberNotFound
	}

	return nil
}

func (m *MemberETCDRepository) SaveLastUpdatedTime(ctx context.Context, member usecases.Member) error {
	findMember, err := m.FindMemberByID(ctx, member.ID)
	if err != nil {
		return err
	}

	now := time.Now()
	findMember.LastUpdatedTime = &now

	return nil
}

func (m *MemberETCDRepository) GetCurrentPartitionOfTheMember(ctx context.Context,
	member usecases.Member) (usecases.Partition, error) {
	findMember, err := m.FindMemberByID(ctx, member.ID)
	if err != nil {
		return usecases.Partition{}, err
	}

	return findMember.Partition, nil
}

func (m *MemberETCDRepository) FindMemberByID(ctx context.Context, id string) (*usecases.Member, error) {
	response, err := m.e.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("could not fetch the member: %w", err)
	}
	if response.Count != 1 {
		return nil, usecases.ErrMemberNotFound
	}
	member := &usecases.Member{}
	if err = json.Unmarshal(response.Kvs[0].Value, member); err != nil {
		return nil, fmt.Errorf("could not unmarshall member: %w", err)
	}

	return member, nil
}

func (m *MemberETCDRepository) GetAllMembers(ctx context.Context) ([]*usecases.Member, error) {
	allClerkRecords, err := m.e.Get(ctx, ETCDRecordPrefix, clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("could not fetch all members: %w", err)
	}
	members := make([]*usecases.Member, allClerkRecords.Count)
	for i := range allClerkRecords.Kvs {
		member := &usecases.Member{}
		if err := json.Unmarshal(allClerkRecords.Kvs[i].Value, member); err != nil {
			return nil, fmt.Errorf("could not unmarshall: %w", err)
		}
		members[i] = member
	}

	return members, nil
}

func (m *MemberETCDRepository) UpdatePartitions(ctx context.Context, idPartitionMap map[string]usecases.Partition) error {
	for i := range members {
		partition, ok := idPartitionMap[members[i].ID]
		if ok {
			members[i].Partition = partition
		}
	}

	return nil
}

func getETCDPrefix(group string) string {
	return fmt.Sprintf("%s-%s-", ETCDRecordPrefix, group)
}
