//go:build integration

package service_test

import (
	"testing"
	"time"
	"uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"neupaneanish.com.np/profile/internal/enum"
	"neupaneanish.com.np/profile/internal/errs"
	gatewayProfilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
)

func TestDeleteExperience(t *testing.T) {
	t.Parallel()

	t.Run("No Experience", func(t *testing.T) {
		t.Parallel()
		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleUser)
		req := &gatewayProfilev1.DeleteExperienceRequest{
			Id:        uuid.NewV7().String(),
			UpdatedAt: timestamppb.New(time.Now()),
		}

		res, err := gatewayProfileServiceClient.DeleteExperience(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrConflict, err)
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		userID := uuid.NewV7()
		ctx := contextWithValue(t, userID, enum.UserRoleUser)
		experience := getExperience(t, userID)

		req := &gatewayProfilev1.DeleteExperienceRequest{
			Id:        experience.ID.String(),
			UpdatedAt: timestamppb.New(experience.UpdatedAt),
		}
		res, err := gatewayProfileServiceClient.DeleteExperience(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, res)
	})
}
