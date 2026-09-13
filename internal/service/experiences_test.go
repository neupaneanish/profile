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

func TestExperiences(t *testing.T) {
	t.Parallel()

	t.Run("Success Gateway", func(t *testing.T) {
		t.Parallel()
		userID := uuid.NewV7()
		ctx := contextWithValue(t, userID, enum.UserRoleUser)

		req := &gatewayProfilev1.ExperiencesRequest{}
		res, err := gatewayProfileServiceClient.Experiences(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, res)
		assert.Empty(t, res.GetExperiences())
	})

	t.Run("Success Root", func(t *testing.T) {
		t.Parallel()
		userID := uuid.NewV7()
		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)
		seedExperience(t, userID)

		req := &rootProfilev1.ExperiencesRequest{UserId: userID.String()}
		res, err := rootProfileServiceClient.Experiences(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, res)
		assert.Len(t, res.GetExperiences(), 1)
	})
}
