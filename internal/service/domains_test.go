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

func TestDomains(t *testing.T) {
	t.Parallel()

	t.Run("Gateway Success", func(t *testing.T) {
		t.Parallel()
		userID := uuid.NewV7()

		ctx := contextWithValue(t, userID, enum.UserRoleUser)
		url := strings.ToLower(rand.Text()[:8]) + ".com"

		seedDomain(t, userID, url, "254.132.101", "A")

		req := &gatewayProfilev1.DomainsRequest{}
		res, err := gatewayProfileServiceClient.Domains(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, res)
		assert.Len(t, res.GetDomains(), 1)
	})

	t.Run("Root Success", func(t *testing.T) {
		t.Parallel()
		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)

		req := &rootProfilev1.DomainsRequest{UserId: uuid.NewV7().String()}
		res, err := rootProfileServiceClient.Domains(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, res)
		assert.Len(t, res.GetDomains(), 0)
	})
}
