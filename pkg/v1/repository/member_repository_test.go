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

	testCreationTime = time.Date(2024, 2, 15, 9, 57, 45, 0, time.UTC)
	testUpdateTime   = time.Date(2024, 2, 15, 10, 5, 45, 0, time.UTC)
)

func (s *ETCDTestSuite) Test_SaveMember() {
	ctx := context.Background()
	err := s.r.SaveNewMember(ctx, fistTestMember)
	require.NoError(s.T(), err)

	member, err := s.r.FindMemberByID(ctx, fistTestMember.ID)
	require.NoError(s.T(), err)
	require.Equal(s.T(), fistTestMember, *member)

}

func (s *ETCDTestSuite) Test_GetAllMembers() {
	ctx := context.Background()
	err := s.r.SaveNewMember(ctx, fistTestMember)
	require.NoError(s.T(), err)

	err = s.r.SaveNewMember(ctx, secondTestMember)
	require.NoError(s.T(), err)

	allMembers, err := s.r.GetAllMembers(ctx)
	require.NoError(s.T(), err)
	require.Len(s.T(), allMembers, 2)
}

func (s *ETCDTestSuite) Test_FindMemberByID() {
	ctx := context.Background()
	err := s.r.SaveNewMember(ctx, fistTestMember)
	require.NoError(s.T(), err)

	actualMember, err := s.r.FindMemberByID(ctx, fistTestMember.ID)
	require.NoError(s.T(), err)
	require.NotZero(s.T(), actualMember)
}

func (s *ETCDTestSuite) Test_DeleteMemberFrom() {
	ctx := context.Background()
	err := s.r.SaveNewMember(ctx, fistTestMember)
	require.NoError(s.T(), err)

	err = s.r.DeleteMemberByID(ctx, fistTestMember.ID)
	require.NoError(s.T(), err)

	_, err = s.r.FindMemberByID(ctx, fistTestMember.ID)
	require.ErrorIs(s.T(), err, usecases.ErrMemberNotFound)
}

func (s *ETCDTestSuite) Test_SaveLastUpdatedTime() {
	ctx := context.Background()
	err := s.r.SaveNewMember(ctx, fistTestMember)
	require.NoError(s.T(), err)

	err = s.r.SaveLastUpdatedTimeByID(ctx, fistTestMember.ID, testUpdateTime)
	require.NoError(s.T(), err)

	member, err := s.r.FindMemberByID(ctx, fistTestMember.ID)
	require.NoError(s.T(), err)
	require.Equal(s.T(), testUpdateTime, *member.LastUpdatedTime)
}
