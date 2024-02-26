package usecases

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMember_IsActive(t *testing.T) {
	t.Run("should return true if time after the duration ", func(t *testing.T) {
		member := Member{}
		member.CreatedAt = time.Now()

		require.True(t, member.IsActiveForTheLast(2*time.Second))
	})
	t.Run("should return false if time before the duration", func(t *testing.T) {
		member := Member{}
		member.CreatedAt = time.Now().Add(-3 * time.Second)

		require.False(t, member.IsActiveForTheLast(2*time.Second))
	})

	t.Run("should use last updated time if it is not empty", func(t *testing.T) {
		member := Member{}
		member.CreatedAt = time.Now().Add(-3 * time.Second)
		lastUpdateTime := time.Now().Add(-1 * time.Second)
		member.LastUpdatedTime = &lastUpdateTime

		require.True(t, member.IsActiveForTheLast(2*time.Second))

		member.CreatedAt = time.Now().Add(-1 * time.Second)
		lastUpdateTime = time.Now().Add(-5 * time.Second)
		member.LastUpdatedTime = &lastUpdateTime

		require.False(t, member.IsActiveForTheLast(2*time.Second))
	})
}
