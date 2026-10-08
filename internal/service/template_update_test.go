//go:build integration

package service_test

import (
	"fmt"
	rand2 "math/rand/v2"
	"testing"
	"uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"neupaneanish.com.np/profile/internal/enum"
	"neupaneanish.com.np/profile/internal/errs"
	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
)

func TestUpdateTemplate(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		template := getTemplate(t)

		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)

		req := &rootProfilev1.UpdateTemplateRequest{
			Id: template.ID.String(),
			Template: &rootProfilev1.Template{
				IconId:      template.IconID.String(),
				Name:        fmt.Sprintf("Name V%d", rand2.Int64N(999999)),
				Description: template.Description,
			},
			UpdatedAt: timestamppb.New(template.UpdatedAt),
		}

		res, err := rootProfileServiceClient.UpdateTemplate(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, res)
	})

	t.Run("UniqueViolation", func(t *testing.T) {
		t.Parallel()

		template := getTemplate(t)
		template2 := getTemplate(t)

		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)

		req := &rootProfilev1.UpdateTemplateRequest{
			Id: template.ID.String(),
			Template: &rootProfilev1.Template{
				IconId:      template.IconID.String(),
				Name:        template2.Name,
				Description: template.Description,
			},
			UpdatedAt: timestamppb.New(template.UpdatedAt),
		}

		res, err := rootProfileServiceClient.UpdateTemplate(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrUniqueViolation("Name"), err)
	})

	t.Run("ForeignKeyViolation", func(t *testing.T) {
		t.Parallel()

		template := getTemplate(t)

		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)

		req := &rootProfilev1.UpdateTemplateRequest{
			Id: template.ID.String(),
			Template: &rootProfilev1.Template{
				IconId:      uuid.NewV7().String(),
				Name:        template.Name,
				Description: template.Description,
			},
			UpdatedAt: timestamppb.New(template.UpdatedAt),
		}

		res, err := rootProfileServiceClient.UpdateTemplate(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrForeignKeyViolation("Icon ID"), err)
	})

	t.Run("No Update", func(t *testing.T) {
		t.Parallel()

		template := getTemplate(t)

		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)

		req := &rootProfilev1.UpdateTemplateRequest{
			Id: template.ID.String(),
			Template: &rootProfilev1.Template{
				IconId:      template.IconID.String(),
				Name:        template.Name,
				Description: template.Description,
			},
			UpdatedAt: timestamppb.New(template.UpdatedAt),
		}

		res, err := rootProfileServiceClient.UpdateTemplate(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrConflict, err)
	})
}
