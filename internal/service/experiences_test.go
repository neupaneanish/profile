//go:build integration

package service_test

import (
	"crypto/rand"
	"testing"
	"uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"neupaneanish.com.np/profile/internal/enum"
	externalProfilev1 "neupaneanish.com.np/profile/internal/protobuf/external/profile/v1"
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

	t.Run("Success External", func(t *testing.T) {
		t.Parallel()
		hostname := rand.Text()[:8] + ".com"

		userID := uuid.NewV7()
		seedExperience(t, userID)

		ctx := externalContextWithValue(t, userID, hostname)

		req := &externalProfilev1.ExperiencesRequest{}
		res, err := externalProfileServiceClient.Experiences(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, res)
		assert.Len(t, res.GetExperiences(), 1)
	})
}
