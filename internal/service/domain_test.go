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

func TestDomain(t *testing.T) {
	t.Parallel()

	t.Run("Gateway Success", func(t *testing.T) {
		t.Parallel()
		userID := uuid.NewV7()

		ctx := contextWithValue(t, userID, enum.UserRoleUser)

		id := seedDomain(t, userID, getCname(), getHostname())

		req := &gatewayProfilev1.DomainRequest{Id: id.String()}

		res, err := gatewayProfileServiceClient.Domain(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, res)
		assert.NotEmpty(t, res.GetUserId())
	})

	t.Run("Gateway Not Found", func(t *testing.T) {
		t.Parallel()
		userID := uuid.NewV7()

		ctx := contextWithValue(t, userID, enum.UserRoleUser)

		req := &gatewayProfilev1.DomainRequest{Id: uuid.NewV7().String()}

		res, err := gatewayProfileServiceClient.Domain(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrNotFound("Domain"), err)
	})

	t.Run("Root Success", func(t *testing.T) {
		t.Parallel()
		userID := uuid.NewV7()

		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)

		id := seedDomain(t, userID, getCname(), getHostname())

		req := &rootProfilev1.DomainRequest{Id: id.String(), UserId: userID.String()}

		res, err := rootProfileServiceClient.Domain(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, res)
		assert.NotEmpty(t, res.GetId())
	})

	t.Run("Success Not Found", func(t *testing.T) {
		t.Parallel()
		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)

		req := &rootProfilev1.DomainRequest{Id: uuid.NewV7().String(), UserId: uuid.NewV7().String()}

		res, err := rootProfileServiceClient.Domain(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)

		assert.Equal(t, errs.ErrNotFound("Domain"), err)
	})
}
