//go:build integration

package service_test

import (
	"crypto/rand"
	"strings"
	"testing"
	"uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"neupaneanish.com.np/profile/internal/enum"
	"neupaneanish.com.np/profile/internal/errs"
	profilev1 "neupaneanish.com.np/profile/internal/protobuf/common/profile/v1"
	gatewayProfilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
)

func TestUpdateSocial(t *testing.T) {
	t.Parallel()

	t.Run("Gateway Success", func(t *testing.T) {
		t.Parallel()
		userID := uuid.NewV7()
		username := strings.ToLower(rand.Text()[:8])

		ctx := contextWithValue(t, userID, enum.UserRoleUser)
		social := getSocial(t, userID)

		req := &gatewayProfilev1.UpdateSocialRequest{Social: &profilev1.UpdateSocial{
			Id:        social.ID.String(),
			Username:  username + "1",
			UpdatedAt: timestamppb.New(social.UpdatedAt),
		}}

		res, err := gatewayProfileServiceClient.UpdateSocial(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, social.ID.String(), res.GetId())
	})

	t.Run("Gateway Error", func(t *testing.T) {
		t.Parallel()
		userID := uuid.NewV7()

		ctx := contextWithValue(t, userID, enum.UserRoleUser)
		social := getSocial(t, userID)

		req := &gatewayProfilev1.UpdateSocialRequest{Social: &profilev1.UpdateSocial{
			Id:        social.ID.String(),
			Username:  social.Username,
			UpdatedAt: timestamppb.New(social.UpdatedAt),
		}}

		res, err := gatewayProfileServiceClient.UpdateSocial(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrConflict, err)
	})

	t.Run("Root Success", func(t *testing.T) {
		t.Parallel()
		userID := uuid.NewV7()
		username := strings.ToLower(rand.Text()[:8])

		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)
		social := getSocial(t, userID)

		req := &rootProfilev1.UpdateSocialRequest{
			Social: &profilev1.UpdateSocial{
				Id:        social.ID.String(),
				Username:  username + "1",
				UpdatedAt: timestamppb.New(social.UpdatedAt),
			},
			UserId: userID.String(),
		}

		res, err := rootProfileServiceClient.UpdateSocial(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, social.ID.String(), res.GetId())
	})

	t.Run("Root Error", func(t *testing.T) {
		t.Parallel()
		userID := uuid.NewV7()
		username := strings.ToLower(rand.Text()[:8])

		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)
		social := getSocial(t, userID)

		req := &rootProfilev1.UpdateSocialRequest{
			Social: &profilev1.UpdateSocial{
				Id:        social.ID.String(),
				Username:  username,
				UpdatedAt: timestamppb.Now(),
			},
			UserId: userID.String(),
		}

		res, err := rootProfileServiceClient.UpdateSocial(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrConflict, err)
	})
}
