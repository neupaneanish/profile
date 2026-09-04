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

func TestGatewayEducations(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		userID := uuid.NewV7()
		ctx := contextWithValue(t, userID, enum.UserRoleUser)

		req := &gatewayProfilev1.EducationsRequest{}
		res, err := gatewayProfileServiceClient.Educations(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, res)
		assert.Len(t, res.GetEducations(), 0)
	})
}

func TestRootEducations(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		userID := uuid.NewV7()
		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)
		seedEducation(t, userID)

		req := &rootProfilev1.EducationsRequest{UserId: userID.String()}
		res, err := rootProfileServiceClient.Educations(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, res)
		assert.Len(t, res.GetEducations(), 1)
	})
}
