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

func TestDeleteDomain(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		userID := uuid.NewV7()

		ctx := contextWithValue(t, userID, enum.UserRoleUser)
		url := strings.ToLower(rand.Text()[:8]) + ".com"

		domain := getDomain(t, userID, url, "254.132.100", "A")

		req := &profilev1.DeleteDomainRequest{
			Id:        domain.ID.String(),
			UpdatedAt: timestamppb.New(domain.UpdatedAt),
		}

		res, err := gatewayProfileServiceClient.DeleteDomain(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, res)
	})

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		userID := uuid.NewV7()

		ctx := contextWithValue(t, userID, enum.UserRoleUser)

		req := &profilev1.DeleteDomainRequest{
			Id:        userID.String(),
			UpdatedAt: timestamppb.Now(),
		}

		res, err := gatewayProfileServiceClient.DeleteDomain(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrConflict, err)
	})
}
