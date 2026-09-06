//go:build integration

package service_test

import (
	"testing"
	"uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"neupaneanish.com.np/profile/internal/enum"
	"neupaneanish.com.np/profile/internal/errs"
	profilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
)

func TestPlatform(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)
		id := seedPlatform(t)

		req := &profilev1.PlatformRequest{Id: id.String()}

		res, err := rootProfileServiceClient.Platform(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, id.String(), res.GetId())
	})

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)

		req := &profilev1.PlatformRequest{Id: uuid.NewV7().String()}

		res, err := rootProfileServiceClient.Platform(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrNotFound, err)
	})
}
