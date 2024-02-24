package repository

import (
	"context"
	"encoding/json"

	"github.com/denizgursoy/clerk/pkg/v1/usecases"
	"github.com/stretchr/testify/require"
)

func (s *ETCDTestSuite) TestSaveMember() {
	ctx := context.Background()
	member, err := s.r.SaveNewMemberToGroup(ctx, "test-group")
	require.NoError(s.T(), err)
	require.NotZero(s.T(), member)

	get, err := s.etcdClient.Get(ctx, member.ID)
	require.NoError(s.T(), err)

	parsedMember := usecases.Member{}
	err = json.Unmarshal(get.Kvs[0].Value, &parsedMember)
	require.NoError(s.T(), err)
	require.NotZero(s.T(), parsedMember)

}

func (s *ETCDTestSuite) TestGetAllMembers() {
	ctx := context.Background()
	_, err := s.r.SaveNewMemberToGroup(ctx, "test-group")
	require.NoError(s.T(), err)

	_, err = s.r.SaveNewMemberToGroup(ctx, "test-group")
	require.NoError(s.T(), err)

	allMembers, err := s.r.GetAllMembers(ctx)
	require.NoError(s.T(), err)
	require.Len(s.T(), allMembers, 2)
}
