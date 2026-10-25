//go:build integration

package service_test

import (
	"testing"
	"uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"neupaneanish.com.np/profile/internal/enum"
	gatewayProfilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
)

func TestCounts(t *testing.T) {
	t.Parallel()

	t.Run("Gateway Counts", func(t *testing.T) {
		t.Parallel()
		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleUser)

		req := &gatewayProfilev1.CountsRequest{}

		res, err := gatewayProfileServiceClient.Counts(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, res)
		assert.GreaterOrEqual(t, res.Socials, int64(0))
	})

	t.Run("Root Counts", func(t *testing.T) {
		t.Parallel()

		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)

		req := &rootProfilev1.CountsRequest{}

		res, err := rootProfileServiceClient.Counts(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, res)
		assert.GreaterOrEqual(t, res.Profiles, int64(0))
	})
}
