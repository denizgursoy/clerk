package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	testConfig = ClerkServerConfig{
		Address:                    "localhost:8080",
		KeepAliveDurationInSeconds: 1,
	}
	testContext = context.Background()
)

func TestNewClerkClient(t *testing.T) {
	t.Run("should ", func(t *testing.T) {
		client, err := NewClerkClient(testConfig)
		require.NoError(t, err)

		member, err := client.AddMember(testContext, "sds", A)
		require.NoError(t, err)
		require.NotZero(t, member)
	})
}

func A(ctx context.Context, ordinal, total int64) {

}
