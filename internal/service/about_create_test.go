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
	profilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
)

func TestCreateAbout(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleUser)

		req := &profilev1.CreateAboutRequest{About: rand.Text()}

		res, err := gatewayProfileServiceClient.CreateAbout(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, res)
		assert.NotEmpty(t, res.GetAbout().GetAbout())
	})

	t.Run("UniqueViolation", func(t *testing.T) {
		t.Parallel()
		about := seedAbout(t)
		ctx := contextWithValue(t, about.UserID, enum.UserRoleUser)

		req := &profilev1.CreateAboutRequest{About: about.About}

		res, err := gatewayProfileServiceClient.CreateAbout(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrUniqueViolation("About"), err)
	})
}
