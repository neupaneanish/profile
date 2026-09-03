//go:build integration

package service_test

import (
	"testing"
	"uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"neupaneanish.com.np/profile/internal/enum"
	"neupaneanish.com.np/profile/internal/errs"
	gatewayProfilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
)

func TestGatewayProfile(t *testing.T) {
	t.Parallel()

	t.Run("Not Found", func(t *testing.T) {
		t.Parallel()

		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleUser)
		req := &gatewayProfilev1.ProfileRequest{}
		res, err := gatewayProfileServiceClient.Profile(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrProfileNotFound, err)
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		userID := uuid.NewV7()
		seedProfile(t, userID)

		ctx := contextWithValue(t, userID, enum.UserRoleRoot)
		req := &gatewayProfilev1.ProfileRequest{}
		res, err := gatewayProfileServiceClient.Profile(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, userID.String(), res.GetProfile().GetUserId())
	})
}

func TestRootProfile(t *testing.T) {
	t.Parallel()

	t.Run("Not Found", func(t *testing.T) {
		t.Parallel()

		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)
		req := &rootProfilev1.ProfileRequest{UserId: uuid.NewV7().String()}
		res, err := rootProfileServiceClient.Profile(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrProfileNotFound, err)
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		userID := uuid.NewV7()
		seedProfile(t, userID)

		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)
		req := &rootProfilev1.ProfileRequest{UserId: userID.String()}
		res, err := rootProfileServiceClient.Profile(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, userID.String(), res.GetProfile().GetUserId())
	})
}
