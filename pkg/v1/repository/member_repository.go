package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
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
	etcdClient *clientv3.Client
}

func NewMemberETCDRepository(c *clientv3.Client) *MemberETCDRepository {
	return &MemberETCDRepository{
		etcdClient: c,
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
	_, err = m.etcdClient.Put(ctx, member.ID, string(marshal))
	if err != nil {
		return usecases.Member{}, fmt.Errorf("could not save member to etcd: %w", err)
	}

	return member, nil
}

func (m *MemberETCDRepository) DeleteMemberFrom(ctx context.Context, member usecases.Member) error {
	for i, mem := range members {
		if mem.ID == member.ID {
			members = slices.Delete(members, i, i+1)

			return nil
		}
	}

	return usecases.ErrMemberNotFound
}

func (m *MemberETCDRepository) SaveLastUpdatedTime(ctx context.Context, member usecases.Member) error {
	findMember, err := m.findMember(member)
	if err != nil {
		return err
	}

	now := time.Now()
	findMember.LastUpdatedTime = &now

	return nil
}

func (m *MemberETCDRepository) GetCurrentPartitionOfTheMember(ctx context.Context,
	member usecases.Member) (usecases.Partition, error) {
	findMember, err := m.findMember(member)
	if err != nil {
		return usecases.Partition{}, err
	}

	return findMember.Partition, nil
}

func (m *MemberETCDRepository) findMember(member usecases.Member) (*usecases.Member, error) {
	for i := range members {
		if members[i].ID == member.ID {
			return members[i], nil
		}
	}

	return nil, usecases.ErrMemberNotFound
}

func (m *MemberETCDRepository) GetAllMembers(ctx context.Context) ([]*usecases.Member, error) {
	allClerkRecords, err := m.etcdClient.Get(ctx, ETCDRecordPrefix, clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("could not fet all members: %w", err)
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
