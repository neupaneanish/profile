//go:build integration

package service_test

import (
	"crypto/rand"
	"testing"
	"uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"neupaneanish.com.np/profile/internal/enum"
	"neupaneanish.com.np/profile/internal/errs"
	gatewayProfilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
)

func TestUpdateAbout(t *testing.T) {
	t.Parallel()

	t.Run("Gateway Success", func(t *testing.T) {
		t.Parallel()

		about := seedAbout(t)
		ctx := contextWithValue(t, about.UserID, enum.UserRoleUser)

		req := &gatewayProfilev1.UpdateAboutRequest{
			About:     about.About + rand.Text(),
			UpdatedAt: timestamppb.New(about.UpdatedAt),
		}

		res, err := gatewayProfileServiceClient.UpdateAbout(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, res)
	})

	t.Run("Gateway Error", func(t *testing.T) {
		t.Parallel()

		about := seedAbout(t)
		ctx := contextWithValue(t, about.UserID, enum.UserRoleUser)

		req := &gatewayProfilev1.UpdateAboutRequest{
			About:     about.About,
			UpdatedAt: timestamppb.New(about.UpdatedAt),
		}

		res, err := gatewayProfileServiceClient.UpdateAbout(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrConflict, err)
	})

	t.Run("Root Success", func(t *testing.T) {
		t.Parallel()

		about := seedAbout(t)
		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)

		req := &rootProfilev1.UpdateAboutRequest{
			UserId:    about.UserID.String(),
			About:     about.About + rand.Text(),
			UpdatedAt: timestamppb.New(about.UpdatedAt),
		}

		res, err := rootProfileServiceClient.UpdateAbout(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, res)
	})

	t.Run("Root Error", func(t *testing.T) {
		t.Parallel()
		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)

		req := &rootProfilev1.UpdateAboutRequest{
			UserId:    uuid.NewV7().String(),
			About:     rand.Text(),
			UpdatedAt: timestamppb.Now(),
		}

		res, err := rootProfileServiceClient.UpdateAbout(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrConflict, err)
	})
}
