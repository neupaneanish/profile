//go:build integration

package service_test

import (
	"crypto/rand"
	"testing"
	"uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"neupaneanish.com.np/profile/internal/enum"
	gatewayProfilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
)

func TestTemplates(t *testing.T) {
	t.Parallel()

	t.Run("Gateway Success", func(t *testing.T) {
		t.Parallel()
		seedTemplate(t)
		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleUser)
		req := &gatewayProfilev1.TemplatesRequest{}

		res, err := gatewayProfileServiceClient.Templates(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, res)
		assert.GreaterOrEqual(t, len(res.GetTemplates()), 1)
	})

	t.Run("Root Success", func(t *testing.T) {
		t.Parallel()
		seedTemplate(t)
		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)
		req := &rootProfilev1.TemplatesRequest{}

		res, err := rootProfileServiceClient.Templates(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, res)
		assert.GreaterOrEqual(t, len(res.GetTemplates()), 1)
	})

	t.Run("Template Icons Success", func(t *testing.T) {
		t.Parallel()

		seedIcon(t, rand.Text(), nil)

		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)
		req := &rootProfilev1.TemplateIconsRequest{}
		res, err := rootProfileServiceClient.TemplateIcons(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, res)
		assert.GreaterOrEqual(t, len(res.GetIcons()), 1)
	})
}
