//go:build integration

package service_test

import (
	"crypto/rand"
	"testing"
	"uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"neupaneanish.com.np/profile/internal/enum"
	externalProfilev1 "neupaneanish.com.np/profile/internal/protobuf/external/profile/v1"
	gatewayProfilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
)

func TestExists(t *testing.T) {
	t.Parallel()
	t.Run("Gateway Exists", func(t *testing.T) {
		t.Parallel()
		userID := uuid.NewV7()
		seedProfile(t, userID)

		ctx := contextWithValue(t, userID, enum.UserRoleUser)
		req := &gatewayProfilev1.ExistsRequest{}
		res, err := gatewayProfileServiceClient.Exists(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, res)
		assert.True(t, res.GetExists().GetProfile())
	})

	t.Run("External Exists", func(t *testing.T) {
		t.Parallel()
		userID := uuid.NewV7()

		ctx := externalContextWithValue(t, userID, rand.Text()[:8]+".com")
		req := &externalProfilev1.ExistsRequest{}
		res, err := externalProfileServiceClient.Exists(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, res)
		assert.False(t, res.GetExists().GetSocials())
	})
}
