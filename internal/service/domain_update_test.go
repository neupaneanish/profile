//go:build integration

package service_test

import (
	"fmt"
	"testing"
	"uuid"

	"math/rand/v2"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"neupaneanish.com.np/profile/internal/enum"
	"neupaneanish.com.np/profile/internal/errs"
	profilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
)

func TestUpdateDomain(t *testing.T) {
	t.Parallel()
	t.Run("Not Verified", func(t *testing.T) {
		t.Parallel()

		userID := uuid.NewV7()

		ctx := contextWithValue(t, userID, enum.UserRoleUser)

		domain := getDomain(t, userID, fmt.Sprintf("cname%d", rand.Int64N(9999999999)), "neupaneanish.com.np", false)

		req := &profilev1.VerifyDomainRequest{
			Id:        domain.ID.String(),
			Txt:       domain.Txt,
			Hostname:  domain.Hostname,
			UpdatedAt: timestamppb.New(domain.UpdatedAt),
		}

		res, err := gatewayProfileServiceClient.VerifyDomain(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrNotFound("TXT"), err)
	})

	t.Run("Update Template ID Success", func(t *testing.T) {
		t.Parallel()
		userID := uuid.NewV7()
		domain := getDomain(t, userID, getCname(), getHostname(), true)
		templateID := seedTemplate(t)

		ctx := contextWithValue(t, userID, enum.UserRoleUser)

		req := &profilev1.UpdateDomainTemplateRequest{
			Id:         domain.ID.String(),
			TemplateId: templateID.String(),
			UpdatedAt:  timestamppb.New(domain.UpdatedAt),
		}

		res, err := gatewayProfileServiceClient.UpdateDomainTemplate(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, res)
	})

	t.Run("Update Template ID ForeignKeyViolation", func(t *testing.T) {
		t.Parallel()
		userID := uuid.NewV7()
		domain := getDomain(t, userID, getCname(), getHostname(), true)

		ctx := contextWithValue(t, userID, enum.UserRoleUser)

		req := &profilev1.UpdateDomainTemplateRequest{
			Id:         domain.ID.String(),
			TemplateId: uuid.NewV7().String(),
			UpdatedAt:  timestamppb.New(domain.UpdatedAt),
		}

		res, err := gatewayProfileServiceClient.UpdateDomainTemplate(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrForeignKeyViolation("Template ID"), err)
	})

	t.Run("Not Found", func(t *testing.T) {
		t.Parallel()
		userID := uuid.NewV7()
		ctx := contextWithValue(t, userID, enum.UserRoleUser)

		req := &profilev1.UpdateDomainTemplateRequest{
			Id:         uuid.NewV7().String(),
			TemplateId: uuid.NewV7().String(),
			UpdatedAt:  timestamppb.Now(),
		}

		res, err := gatewayProfileServiceClient.UpdateDomainTemplate(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrConflict, err)
	})
}
