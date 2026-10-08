//go:build integration

package service_test

import (
	"crypto/rand"
	"testing"
	"uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"neupaneanish.com.np/profile/internal/enum"
	"neupaneanish.com.np/profile/internal/errs"
	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
)

func TestSocial(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)
		userID := uuid.NewV7()
		_, id := seedSocial(t, userID, rand.Text()[:8])

		req := &rootProfilev1.SocialRequest{
			Id:     id.String(),
			UserId: userID.String(),
		}

		res, err := rootProfileServiceClient.Social(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, id.String(), res.GetId())
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)
		userID := uuid.NewV7()

		req := &rootProfilev1.SocialRequest{
			Id:     uuid.NewV7().String(),
			UserId: userID.String(),
		}

		res, err := rootProfileServiceClient.Social(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrNotFound("Social"), err)
	})
}
