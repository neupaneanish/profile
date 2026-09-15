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
	gatewayProfilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
)

func TestIcons(t *testing.T) {
	t.Parallel()

	t.Run("Social Icons Success", func(t *testing.T) {
		t.Parallel()
		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleUser)

		req := &gatewayProfilev1.SocialIconsRequest{}

		res, err := gatewayProfileServiceClient.SocialIcons(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, res)
		assert.GreaterOrEqual(t, len(res.GetIcons()), 0)
	})

	t.Run("Icons", func(t *testing.T) {
		t.Parallel()

		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)

		name := strings.ToLower(rand.Text()[:8])
		seedIcon(t, name, nil)
		seedIcon(t, name+name, nil)

		req := &rootProfilev1.IconsRequest{}

		res, err := rootProfileServiceClient.Icons(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, res)
		assert.GreaterOrEqual(t, len(res.GetIcons()), 2)
	})
}
