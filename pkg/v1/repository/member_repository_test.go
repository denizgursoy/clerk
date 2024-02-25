package repository

import (
	"context"
	"time"

	"github.com/denizgursoy/clerk/pkg/v1/usecases"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

const TestGroup = "test-group"

var (
	fistTestMember = usecases.Member{
		Group:     TestGroup,
		ID:        usecases.GenerateMemberName(TestGroup, uuid.NewString()),
		CreatedAt: testCreationTime,
	}
	secondTestMember = usecases.Member{
		Group:     TestGroup,
		ID:        usecases.GenerateMemberName(TestGroup, uuid.NewString()),
		CreatedAt: testCreationTime,
	}
	ctx              = context.Background()
	testCreationTime = time.Date(2024, 2, 15, 9, 57, 45, 0, time.UTC)
	testUpdateTime   = time.Date(2024, 2, 15, 10, 5, 45, 0, time.UTC)
)

func (s *ETCDTestSuite) Test_SaveMember() {
	ctx := context.Background()
	err := s.r.SaveNewMember(ctx, fistTestMember)
	require.NoError(s.T(), err)

	member, err := s.r.FetchMemberByID(ctx, fistTestMember.ID)
	require.NoError(s.T(), err)
	require.Equal(s.T(), fistTestMember, *member)

}

func (s *ETCDTestSuite) Test_GetAllMembers() {
	err := s.r.SaveNewMember(ctx, fistTestMember)
	require.NoError(s.T(), err)

	err = s.r.SaveNewMember(ctx, secondTestMember)
	require.NoError(s.T(), err)

	allMembers, err := s.r.FetchAllMembers(ctx)
	require.NoError(s.T(), err)
	require.Len(s.T(), allMembers, 2)
}

func (s *ETCDTestSuite) Test_FindMemberByID() {
	err := s.r.SaveNewMember(ctx, fistTestMember)
	require.NoError(s.T(), err)

	actualMember, err := s.r.FetchMemberByID(ctx, fistTestMember.ID)
	require.NoError(s.T(), err)
	require.NotZero(s.T(), actualMember)
}

func (s *ETCDTestSuite) Test_DeleteMemberFrom() {
	err := s.r.SaveNewMember(ctx, fistTestMember)
	require.NoError(s.T(), err)

	err = s.r.DeleteMemberByID(ctx, fistTestMember.ID)
	require.NoError(s.T(), err)

	_, err = s.r.FetchMemberByID(ctx, fistTestMember.ID)
	require.ErrorIs(s.T(), err, usecases.ErrMemberNotFound)
}

func (s *ETCDTestSuite) Test_SaveLastUpdatedTime() {
	err := s.r.SaveNewMember(ctx, fistTestMember)
	require.NoError(s.T(), err)

	err = s.r.SaveLastUpdatedTimeByID(ctx, fistTestMember.ID, testUpdateTime)
	require.NoError(s.T(), err)

	member, err := s.r.FetchMemberByID(ctx, fistTestMember.ID)
	require.NoError(s.T(), err)
	require.Equal(s.T(), testUpdateTime, *member.LastUpdatedTime)
}

func (s *ETCDTestSuite) Test_UpdatePartitions() {
	err := s.r.SaveNewMember(ctx, fistTestMember)
	require.NoError(s.T(), err)

	err = s.r.SaveNewMember(ctx, secondTestMember)
	require.NoError(s.T(), err)
	maps := map[string]usecases.Partition{
		fistTestMember.ID: {
			Ordinal: 1,
			Total:   2,
		},
		secondTestMember.ID: {
			Ordinal: 2,
			Total:   2,
		},
	}
	err = s.r.UpdatePartitions(ctx, maps)
	require.NoError(s.T(), err)

	actualFirstMember, err := s.r.FetchMemberByID(ctx, fistTestMember.ID)
	require.NoError(s.T(), err)
	require.Equal(s.T(), maps[fistTestMember.ID], actualFirstMember.Partition)

	actualSecondMember, err := s.r.FetchMemberByID(ctx, secondTestMember.ID)
	require.NoError(s.T(), err)
	require.Equal(s.T(), maps[secondTestMember.ID], actualSecondMember.Partition)
}
