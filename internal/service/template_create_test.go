//go:build integration

package service_test

import (
	"crypto/rand"
	"fmt"
	rand2 "math/rand/v2"
	"testing"
	"uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"neupaneanish.com.np/profile/internal/enum"
	"neupaneanish.com.np/profile/internal/errs"
	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
)

func TestCreateTemplate(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		iconID := seedIcon(t, rand.Text()[:8], nil)

		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)

		req := &rootProfilev1.CreateTemplateRequest{Template: &rootProfilev1.Template{
			IconId:      iconID.String(),
			Name:        fmt.Sprintf("Name V%d", rand2.Int64N(999999)),
			Description: rand.Text(),
		}}

		res, err := rootProfileServiceClient.CreateTemplate(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, res)
	})

	t.Run("UniqueViolation", func(t *testing.T) {
		t.Parallel()

		template := getTemplate(t)

		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)

		req := &rootProfilev1.CreateTemplateRequest{Template: &rootProfilev1.Template{
			IconId:      template.IconID.String(),
			Name:        template.Name,
			Description: template.Description,
		}}
		res, err := rootProfileServiceClient.CreateTemplate(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrUniqueViolation("Name"), err)
	})

	t.Run("ForeignKeyViolation", func(t *testing.T) {
		t.Parallel()

		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)

		req := &rootProfilev1.CreateTemplateRequest{Template: &rootProfilev1.Template{
			IconId:      uuid.NewV7().String(),
			Name:        fmt.Sprintf("Name V%d", rand2.Int64N(999999)),
			Description: rand.Text(),
		}}
		res, err := rootProfileServiceClient.CreateTemplate(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrForeignKeyViolation("Icon ID"), err)
	})
}
