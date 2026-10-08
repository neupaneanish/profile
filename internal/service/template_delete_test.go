//go:build integration

package service_test

import (
	"testing"
	"uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"neupaneanish.com.np/profile/internal/enum"
	"neupaneanish.com.np/profile/internal/errs"
	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
)

func TestDeleteTemplate(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		template := getTemplate(t)
		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)

		req := &rootProfilev1.DeleteTemplateRequest{
			Id:        template.ID.String(),
			UpdatedAt: timestamppb.New(template.UpdatedAt),
		}

		res, err := rootProfileServiceClient.DeleteTemplate(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, res)
	})

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)

		req := &rootProfilev1.DeleteTemplateRequest{
			Id:        uuid.NewV7().String(),
			UpdatedAt: timestamppb.Now(),
		}

		res, err := rootProfileServiceClient.DeleteTemplate(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrConflict, err)
	})
}
