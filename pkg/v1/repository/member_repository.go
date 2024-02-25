package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/denizgursoy/clerk/pkg/v1/usecases"
	"go.etcd.io/etcd/client/v3"
)

type MemberETCDRepository struct {
	e *clientv3.Client
}

func NewMemberETCDRepository(e *clientv3.Client) *MemberETCDRepository {
	return &MemberETCDRepository{
		e: e,
	}
}

func (m *MemberETCDRepository) SaveNewMember(ctx context.Context, member usecases.Member) error {
	marshal, err := json.Marshal(member)
	if err != nil {
		return fmt.Errorf("could not marshal member: %w", err)
	}
	_, err = m.e.Put(ctx, member.ID, string(marshal))
	if err != nil {
		return fmt.Errorf("could not save member to etcd: %w", err)
	}

	return nil
}

func (m *MemberETCDRepository) DeleteMemberByID(ctx context.Context, id string) error {
	response, err := m.e.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("could not delete member: %w", err)
	}
	if response.Deleted == 0 {
		return usecases.ErrMemberNotFound
	}

	return nil
}

func (m *MemberETCDRepository) SaveLastUpdatedTimeByID(ctx context.Context, id string, updateTime time.Time) error {
	member, err := m.FindMemberByID(ctx, id)
	if err != nil {
		return fmt.Errorf("could not find any member to update time: %w", err)
	}
	member.LastUpdatedTime = &updateTime

	if err = m.SaveNewMember(ctx, *member); err != nil {
		return fmt.Errorf("could not update last updated time: %w", err)
	}

	return nil
}

func (m *MemberETCDRepository) GetPartitionOfTheMemberByID(ctx context.Context,
	id string) (usecases.Partition, error) {
	member, err := m.FindMemberByID(ctx, id)
	if err != nil {
		return usecases.Partition{}, err
	}

	return member.Partition, nil
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
	allClerkRecords, err := m.e.Get(ctx, usecases.IDPrefix, clientv3.WithPrefix())
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
	// for i := range members {
	// 	partition, ok := idPartitionMap[members[i].ID]
	// 	if ok {
	// 		members[i].Partition = partition
	// 	}
	// }
	//
	return nil
}
