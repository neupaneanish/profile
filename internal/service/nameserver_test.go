//go:build integration

package service_test

import (
	"testing"
	"uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"neupaneanish.com.np/profile/internal/enum"
	"neupaneanish.com.np/profile/internal/errs"
	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
)

func TestNameserver(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)
		id := seedNameserver(t, getCname(), getHostname())

		req := &rootProfilev1.NameserverRequest{Id: id.String()}
		res, err := rootProfileServiceClient.Nameserver(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, id.String(), req.GetId())
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()
		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)

		req := &rootProfilev1.NameserverRequest{Id: uuid.NewV7().String()}
		res, err := rootProfileServiceClient.Nameserver(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrNotFound("Nameserver"), err)
	})
}
