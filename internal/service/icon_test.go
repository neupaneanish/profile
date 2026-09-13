//go:build integration

package service_test

import (
	"crypto/rand"
	"strings"
	"testing"
	"uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"neupaneanish.com.np/profile/internal/enum"
	"neupaneanish.com.np/profile/internal/errs"
	profilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
)

func TestIcon(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)
		name := strings.ToLower(rand.Text()[:8])
		icon := getIcon(t, name, nil)

		req := &profilev1.IconRequest{Id: icon.ID.String()}

		res, err := rootProfileServiceClient.Icon(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, icon.Name, res.GetName())
	})

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)

		req := &profilev1.IconRequest{Id: uuid.NewV7().String()}

		res, err := rootProfileServiceClient.Icon(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrNotFound("Icon"), err)
	})
}
