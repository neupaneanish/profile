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
	profilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
)

func TestUpdateDomain(t *testing.T) {
	t.Parallel()

	t.Run("No Domain", func(t *testing.T) {
		t.Parallel()

		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleUser)

		req := &profilev1.VerifyDomainRequest{
			Id:        uuid.NewV7().String(),
			UpdatedAt: timestamppb.Now(),
		}

		res, err := gatewayProfileServiceClient.VerifyDomain(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrConflict, err)
	})

	t.Run("Fake Domain", func(t *testing.T) {
		t.Parallel()

		url := strings.ToLower(rand.Text()[:8]) + ".com"
		err := verifyDomainError(t, url, "123.123.123.1")
		assert.Equal(t, errs.ErrNotFound("TXT"), err)
	})

	t.Run("No TXT", func(t *testing.T) {
		t.Parallel()

		err := verifyDomainError(t, "neupaneanish.com.np", "123.123.123.2")
		assert.Equal(t, errs.ErrNotFound("TXT"), err)
	})
}

func verifyDomainError(t *testing.T, url, ip string) error {
	userID := uuid.NewV7()

	ctx := contextWithValue(t, userID, enum.UserRoleUser)

	domain := getDomain(t, userID, url, ip, "A")

	req := &profilev1.VerifyDomainRequest{
		Id:        domain.ID.String(),
		UpdatedAt: timestamppb.New(domain.UpdatedAt),
	}

	res, err := gatewayProfileServiceClient.VerifyDomain(ctx, req)
	require.Error(t, err)
	assert.Nil(t, res)
	return err
}
