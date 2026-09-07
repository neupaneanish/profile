//go:build integration

package service_test

import (
	"crypto/rand"
	"strings"
	"testing"
	"uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"neupaneanish.com.np/profile/internal/enum"
	"neupaneanish.com.np/profile/internal/errs"
	profilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
	"neupaneanish.com.np/profile/internal/repository"
)

func TestCreateSocial(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleUser)

		pID := seedPlatform(t)
		username := strings.ToLower(rand.Text()[:8])

		req := &profilev1.CreateSocialRequest{
			PlatformId: pID.String(),
			Username:   username,
		}

		res, err := gatewayProfileServiceClient.CreateSocial(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, res)
		assert.NotEmpty(t, res.GetId())
	})

	t.Run("UniqueViolation", func(t *testing.T) {
		t.Parallel()

		userID := uuid.NewV7()
		username := strings.ToLower(rand.Text()[:8])
		pID, _ := seedSocial(t, userID, username)

		ctx := contextWithValue(t, userID, enum.UserRoleUser)

		req := &profilev1.CreateSocialRequest{
			PlatformId: pID.String(),
			Username:   username,
		}

		res, err := gatewayProfileServiceClient.CreateSocial(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrUniqueViolation("Social"), err)
	})

	t.Run("ForeignKeyViolation", func(t *testing.T) {
		t.Parallel()

		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleUser)
		username := strings.ToLower(rand.Text()[:8])

		req := &profilev1.CreateSocialRequest{
			PlatformId: uuid.NewV7().String(),
			Username:   username,
		}
		res, err := gatewayProfileServiceClient.CreateSocial(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrForeignKeyViolation("Platform"), err)
	})
}

func seedSocial(t *testing.T, userID uuid.UUID, username string) (uuid.UUID, uuid.UUID) {
	t.Helper()

	pID := seedPlatform(t)

	params := &repository.CreateSocialParams{
		UserID:     userID,
		PlatformID: pID,
		Username:   username,
		CreatedBy:  userID,
		UpdatedBy:  userID,
	}
	id, err := cfg.Repository.CreateSocial(t.Context(), params)
	require.NoError(t, err)
	assert.NotEmpty(t, id)

	return pID, id
}
