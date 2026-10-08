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
	externalProfilev1 "neupaneanish.com.np/profile/internal/protobuf/external/profile/v1"
	gatewayProfilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
)

func TestSocials(t *testing.T) {
	t.Parallel()

	t.Run("Gateway", func(t *testing.T) {
		t.Parallel()

		userID := uuid.NewV7()

		ctx := contextWithValue(t, userID, enum.UserRoleUser)

		req := &gatewayProfilev1.SocialsRequest{}

		res, err := gatewayProfileServiceClient.Socials(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, res)
		assert.Empty(t, res.GetSocials())
	})

	t.Run("Root", func(t *testing.T) {
		t.Parallel()

		userID := uuid.NewV7()
		username := strings.ToLower(rand.Text()[:8])

		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)
		seedSocial(t, userID, username)
		seedSocial(t, userID, username+"1")

		req := &rootProfilev1.SocialsRequest{UserId: userID.String()}

		res, err := rootProfileServiceClient.Socials(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, res)
		assert.Len(t, res.GetSocials(), 2)
	})

	t.Run("External", func(t *testing.T) {
		t.Parallel()
		name := rand.Text()[:8]
		hostname := name + ".com"
		userID := uuid.NewV7()
		seedSocial(t, userID, name)
		ctx := externalContextWithValue(t, userID, hostname)

		req := &externalProfilev1.SocialsRequest{}
		res, err := externalProfileServiceClient.Socials(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, res)
		assert.Len(t, res.GetSocials(), 1)
	})
}
