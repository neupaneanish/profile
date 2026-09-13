//go:build integration

package service_test

import (
	"testing"
	"uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"neupaneanish.com.np/profile/internal/enum"
	"neupaneanish.com.np/profile/internal/errs"
	gatewayProfilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
)

func TestExperience(t *testing.T) {
	t.Parallel()

	t.Run("Not Found Gateway", func(t *testing.T) {
		t.Parallel()
		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleUser)

		req := &gatewayProfilev1.ExperienceRequest{Id: uuid.NewV7().String()}

		res, err := gatewayProfileServiceClient.Experience(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrNotFound("Experience"), err)
	})

	t.Run("Success gateway", func(t *testing.T) {
		t.Parallel()
		userID := uuid.NewV7()

		ctx := contextWithValue(t, userID, enum.UserRoleUser)

		exp := getExperience(t, userID)

		req := &gatewayProfilev1.ExperienceRequest{Id: exp.ID.String()}
		res, err := gatewayProfileServiceClient.Experience(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, res)
	})

	t.Run("Not Found Root", func(t *testing.T) {
		t.Parallel()
		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)

		req := &rootProfilev1.ExperienceRequest{Id: uuid.NewV7().String(), UserId: uuid.NewV7().String()}

		res, err := rootProfileServiceClient.Experience(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrNotFound("Experience"), err)
	})

	t.Run("Success Root", func(t *testing.T) {
		t.Parallel()
		userID := uuid.NewV7()

		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)

		exp := getExperience(t, userID)

		req := &rootProfilev1.ExperienceRequest{Id: exp.ID.String(), UserId: userID.String()}
		res, err := rootProfileServiceClient.Experience(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, res)
	})
}
