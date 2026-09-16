//go:build integration

package service_test

import (
	"crypto/rand"
	"testing"
	"uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
	"neupaneanish.com.np/profile/internal/enum"
	"neupaneanish.com.np/profile/internal/errs"
	externalProfilev1 "neupaneanish.com.np/profile/internal/protobuf/external/profile/v1"
	gatewayProfilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
)

func TestAbout(t *testing.T) {
	t.Parallel()

	t.Run("Gateway Success", func(t *testing.T) {
		t.Parallel()

		about := seedAbout(t)

		ctx := contextWithValue(t, about.UserID, enum.UserRoleUser)
		req := &gatewayProfilev1.AboutRequest{}

		res, err := gatewayProfileServiceClient.About(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, about.UserID.String(), res.GetAbout().GetUserId())
	})

	t.Run("Gateway No Metadata", func(t *testing.T) {
		t.Parallel()

		req := &gatewayProfilev1.AboutRequest{}

		res, err := gatewayProfileServiceClient.About(t.Context(), req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrUnauthenticated, err)
	})

	t.Run("Gateway Invalid UserID", func(t *testing.T) {
		t.Parallel()

		md := metadata.Pairs(
			"x-user-id", rand.Text(),
			"x-role", "role",
			"x-jti", uuid.NewV7().String(),
			"x-username", rand.Text()[:8],
		)

		ctx := metadata.NewOutgoingContext(t.Context(), md)
		req := &gatewayProfilev1.AboutRequest{}

		res, err := gatewayProfileServiceClient.About(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrUnauthenticated, err)
	})

	t.Run("Gateway Invalid Role", func(t *testing.T) {
		t.Parallel()

		ctx := contextWithValue(t, uuid.NewV7(), "Role")
		req := &gatewayProfilev1.AboutRequest{}

		res, err := gatewayProfileServiceClient.About(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrUnauthenticated, err)
	})

	t.Run("Gateway Error", func(t *testing.T) {
		t.Parallel()

		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleUser)
		req := &gatewayProfilev1.AboutRequest{}

		res, err := gatewayProfileServiceClient.About(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrNotFound("About"), err)
	})

	t.Run("Gateway Error Root", func(t *testing.T) {
		t.Parallel()

		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)
		req := &gatewayProfilev1.AboutRequest{}

		res, err := gatewayProfileServiceClient.About(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrPermissionDenied, err)
	})

	t.Run("Root Success", func(t *testing.T) {
		t.Parallel()

		about := seedAbout(t)

		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)
		req := &rootProfilev1.AboutRequest{UserId: about.UserID.String()}

		res, err := rootProfileServiceClient.About(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, about.UserID.String(), res.GetAbout().GetUserId())
	})

	t.Run("Root Error", func(t *testing.T) {
		t.Parallel()

		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)
		req := &rootProfilev1.AboutRequest{UserId: uuid.NewV7().String()}

		res, err := rootProfileServiceClient.About(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrNotFound("About"), err)
	})

	t.Run("Root Error Gateway", func(t *testing.T) {
		t.Parallel()

		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleUser)
		req := &rootProfilev1.AboutRequest{UserId: uuid.NewV7().String()}

		res, err := rootProfileServiceClient.About(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrPermissionDenied, err)
	})

	t.Run("External Success", func(t *testing.T) {
		t.Parallel()

		about := seedAbout(t)

		hostname := rand.Text()[:8] + ".com"

		ctx := externalContextWithValue(t, about.UserID, hostname)
		req := &externalProfilev1.AboutRequest{}

		res, err := externalProfileServiceClient.About(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, res)
		assert.NotEmpty(t, res.GetAbout())
	})

	t.Run("External Error", func(t *testing.T) {
		t.Parallel()

		hostname := rand.Text()[:8] + ".com"

		ctx := externalContextWithValue(t, uuid.NewV7(), hostname)
		req := &externalProfilev1.AboutRequest{}

		res, err := externalProfileServiceClient.About(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrNotFound("About"), err)
	})

	t.Run("External Error no metadata", func(t *testing.T) {
		t.Parallel()

		req := &externalProfilev1.AboutRequest{}

		res, err := externalProfileServiceClient.About(t.Context(), req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrNotFound("Host"), err)
	})

	t.Run("External Error Invalid UserID", func(t *testing.T) {
		t.Parallel()

		hostname := rand.Text()[:8] + ".com"

		md := metadata.Pairs("x-hostname", hostname, "x-user-id", rand.Text())
		ctx := metadata.NewOutgoingContext(t.Context(), md)

		req := &externalProfilev1.AboutRequest{}

		res, err := externalProfileServiceClient.About(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrNotFound("Host"), err)
	})
}
