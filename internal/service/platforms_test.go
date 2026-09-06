//go:build integration

package service_test

import (
	"testing"
	"uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"neupaneanish.com.np/profile/internal/enum"
	gatewayProfilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
)

func TestGatewayPlatform(t *testing.T) {
	t.Parallel()

	ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleUser)

	req := &gatewayProfilev1.PlatformsRequest{}

	res, err := gatewayProfileServiceClient.Platforms(ctx, req)

	require.NoError(t, err)
	assert.NotNil(t, res)
	assert.GreaterOrEqual(t, len(res.GetPlatforms()), 0)
}

func TestRootPlatform(t *testing.T) {
	t.Parallel()

	ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)

	seedPlatform(t)
	seedPlatform(t)

	req := &rootProfilev1.PlatformsRequest{}

	res, err := rootProfileServiceClient.Platforms(ctx, req)

	require.NoError(t, err)
	assert.NotNil(t, res)
	assert.GreaterOrEqual(t, len(res.GetPlatforms()), 0)
}
