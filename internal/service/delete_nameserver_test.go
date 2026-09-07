//go:build integration

package service_test

import (
	"crypto/rand"
	"strings"
	"testing"
	"uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"neupaneanish.com.np/profile/internal/enum"
	"neupaneanish.com.np/profile/internal/errs"
	profilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
)

func TestDeleteNameServer(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)
		name := strings.ToLower(rand.Text()[:8])
		nameserver := getNameServer(t, name)

		req := &profilev1.DeleteNameServerRequest{
			Id:        nameserver.ID.String(),
			UpdatedAt: timestamppb.New(nameserver.UpdatedAt),
		}

		res, err := rootProfileServiceClient.DeleteNameServer(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, res)
		assert.Empty(t, res.String())
	})

	t.Run("Error", func(t *testing.T) {
		t.Parallel()
		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)

		req := &profilev1.DeleteNameServerRequest{
			Id:        uuid.NewV7().String(),
			UpdatedAt: timestamppb.Now(),
		}

		res, err := rootProfileServiceClient.DeleteNameServer(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrConflict, err)
	})
}
