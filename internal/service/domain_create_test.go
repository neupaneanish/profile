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

func TestCreateDomain(t *testing.T) {
	t.Parallel()

	// No nameserver found (ErrInternalServer) because it does already have nameservers because of Parallel

	t.Run("UniqueViolation", func(t *testing.T) {
		t.Parallel()

		userID := uuid.NewV7()

		ctx := contextWithValue(t, userID, enum.UserRoleUser)

		url := strings.ToLower(rand.Text()[:8]) + ".com"

		seedDomain(t, userID, url, "254.132.122", "A")

		req := &profilev1.CreateDomainRequest{Url: url}

		res, err := gatewayProfileServiceClient.CreateDomain(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrUniqueViolation("Domain"), err)
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleUser)

		url := strings.ToLower(rand.Text()[:8]) + ".com"

		req := &profilev1.CreateDomainRequest{Url: url}
		seedNameserver(t, "254.132.127", "A")

		res, err := gatewayProfileServiceClient.CreateDomain(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, res)
	})

	t.Run("Invalid Domain", func(t *testing.T) {
		t.Parallel()

		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleUser)

		url := strings.ToLower(rand.Text()[:8]) + ".abcd"

		req := &profilev1.CreateDomainRequest{Url: url}

		res, err := gatewayProfileServiceClient.CreateDomain(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrInvalidURL, err)
	})
}
