//go:build integration

package service_test

import (
	"crypto/rand"
	"strings"
	"testing"
	"uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	timestimestamppb "google.golang.org/protobuf/types/known/timestamppb"
	"neupaneanish.com.np/profile/internal/enum"
	"neupaneanish.com.np/profile/internal/errs"
	profilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
	"neupaneanish.com.np/profile/internal/repository"
)

func TestDeleteSocial(t *testing.T) {
	t.Parallel()

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		userID := uuid.NewV7()
		ctx := contextWithValue(t, userID, enum.UserRoleUser)

		req := &profilev1.DeleteSocialRequest{
			Id:        uuid.NewV7().String(),
			UpdatedAt: timestimestamppb.Now(),
		}

		res, err := gatewayProfileServiceClient.DeleteSocial(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrConflict, err)
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		userID := uuid.NewV7()

		social := getSocial(t, userID)
		require.NotNil(t, social)

		ctx := contextWithValue(t, social.UserID, enum.UserRoleUser)
		req := &profilev1.DeleteSocialRequest{
			Id:        social.ID.String(),
			UpdatedAt: timestimestamppb.New(social.UpdatedAt),
		}

		res, err := gatewayProfileServiceClient.DeleteSocial(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, res)
	})

	t.Run("Root Error", func(t *testing.T) {
		t.Parallel()

		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)

		req := &rootProfilev1.DeleteSocialRequest{
			Id:        uuid.NewV7().String(),
			UserId:    uuid.NewV7().String(),
			UpdatedAt: timestimestamppb.Now(),
		}
		res, err := rootProfileServiceClient.DeleteSocial(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrConflict, err)
	})

	t.Run("Root Success", func(t *testing.T) {
		t.Parallel()

		social := getSocial(t, uuid.NewV7())
		require.NotNil(t, social)

		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)
		req := &rootProfilev1.DeleteSocialRequest{
			Id:        social.ID.String(),
			UserId:    social.UserID.String(),
			UpdatedAt: timestimestamppb.New(social.UpdatedAt),
		}

		res, err := rootProfileServiceClient.DeleteSocial(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, res)
	})
}

func getSocial(t *testing.T, userID uuid.UUID) *repository.SocialsRow {
	t.Helper()

	username := strings.ToLower(rand.Text()[:8])

	_, id := seedSocial(t, userID, username)

	params := &repository.SocialsParams{UserID: userID}

	rows, err := cfg.Repository.Socials(t.Context(), params)
	require.NoError(t, err)

	for _, s := range rows {
		if s.ID == id {
			return s
		}
	}
	return nil
}
