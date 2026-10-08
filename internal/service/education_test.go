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

func TestEducation(t *testing.T) {
	t.Parallel()

	t.Run("Not Found Gateway", func(t *testing.T) {
		t.Parallel()
		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleUser)

		req := &gatewayProfilev1.EducationRequest{Id: uuid.NewV7().String()}

		res, err := gatewayProfileServiceClient.Education(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrNotFound("Education"), err)
	})

	t.Run("Success Gateway", func(t *testing.T) {
		t.Parallel()
		userID := uuid.NewV7()

		ctx := contextWithValue(t, userID, enum.UserRoleUser)

		edu := getEducation(t, userID)

		req := &gatewayProfilev1.EducationRequest{Id: edu.ID.String()}
		res, err := gatewayProfileServiceClient.Education(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, res)
		assert.NotEmpty(t, res.GetId())
	})

	t.Run("Not Found Root", func(t *testing.T) {
		t.Parallel()
		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)

		req := &rootProfilev1.EducationRequest{Id: uuid.NewV7().String(), UserId: uuid.NewV7().String()}

		res, err := rootProfileServiceClient.Education(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrNotFound("Education"), err)
	})

	t.Run("Success Root", func(t *testing.T) {
		t.Parallel()
		userID := uuid.NewV7()

		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)

		edu := getEducation(t, userID)

		req := &rootProfilev1.EducationRequest{Id: edu.ID.String(), UserId: userID.String()}
		res, err := rootProfileServiceClient.Education(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, res)
	})
}
