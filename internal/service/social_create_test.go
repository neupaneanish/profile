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
)

func TestCreateSocial(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleUser)

		name := strings.ToLower(rand.Text()[:8])
		icon := getIcon(t, name, nil)
		username := strings.ToLower(rand.Text()[:8])

		req := &profilev1.CreateSocialRequest{
			IconId:   icon.ID.String(),
			Username: username,
		}

		res, err := gatewayProfileServiceClient.CreateSocial(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, res)
	})

	t.Run("UniqueViolation", func(t *testing.T) {
		t.Parallel()

		userID := uuid.NewV7()
		username := strings.ToLower(rand.Text()[:8])
		iconID, _ := seedSocial(t, userID, username)

		ctx := contextWithValue(t, userID, enum.UserRoleUser)

		req := &profilev1.CreateSocialRequest{
			IconId:   iconID.String(),
			Username: username,
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
			IconId:   uuid.NewV7().String(),
			Username: username,
		}
		res, err := gatewayProfileServiceClient.CreateSocial(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrForeignKeyViolation("Icon"), err)
	})
}
