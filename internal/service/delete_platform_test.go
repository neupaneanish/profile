//go:build integration

package service_test

import (
	"testing"
	"uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"neupaneanish.com.np/profile/internal/enum"
	"neupaneanish.com.np/profile/internal/errs"
	profilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
)

func TestDeletePlatform(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)

		platform := getPlatform(t)

		req := &profilev1.DeletePlatformRequest{
			Id:        platform.ID.String(),
			UpdatedAt: timestamppb.New(platform.UpdatedAt),
		}

		res, err := rootProfileServiceClient.DeletePlatform(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, res)
	})

	t.Run("Error", func(t *testing.T) {
		t.Parallel()
		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)

		req := &profilev1.DeletePlatformRequest{
			Id:        uuid.NewV7().String(),
			UpdatedAt: timestamppb.Now(),
		}

		res, err := rootProfileServiceClient.DeletePlatform(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrConflict, err)
	})
}
