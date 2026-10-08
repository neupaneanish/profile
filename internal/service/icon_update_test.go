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
	"neupaneanish.com.np/profile/internal/utils"
)

func TestUpdateIcon(t *testing.T) {
	t.Parallel()

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)
		name := strings.ToLower(rand.Text()[:8])

		icon := getIcon(t, name, nil)

		req := &profilev1.UpdateIconRequest{
			Id: icon.ID.String(),
			Icon: &profilev1.CreateUpdateIcon{
				Name:         icon.Name,
				SiteHostname: icon.SiteHostname,
				SiteSuffix:   utils.StringpbValue(icon.SiteSuffix),
				Hostname:     icon.Hostname,
				Suffix:       icon.Suffix,
				Color:        icon.Color,
			},
			UpdatedAt: timestamppb.New(icon.UpdatedAt),
		}

		res, err := rootProfileServiceClient.UpdateIcon(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrConflict, err)
	})

	t.Run("Unique Name Error", func(t *testing.T) {
		t.Parallel()

		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)

		name := strings.ToLower(rand.Text()[:8])

		icon := getIcon(t, name, nil)
		icon1 := getIcon(t, name+name, nil)

		req := &profilev1.UpdateIconRequest{
			Id: icon.ID.String(),
			Icon: &profilev1.CreateUpdateIcon{
				Name:         icon1.Name,
				SiteHostname: icon.SiteHostname,
				SiteSuffix:   utils.StringpbValue(icon.SiteSuffix),
				Hostname:     icon.Hostname,
				Suffix:       icon.Suffix,
				Color:        icon.Color,
			},
			UpdatedAt: timestamppb.New(icon.UpdatedAt),
		}

		res, err := rootProfileServiceClient.UpdateIcon(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrUniqueViolation("Icon Name"), err)
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)

		name := strings.ToLower(rand.Text()[:8])
		icon := getIcon(t, name, nil)

		var siteSuffix *string
		suffix := "/" + icon.Name
		siteSuffix = &suffix

		req := &profilev1.UpdateIconRequest{
			Id: icon.ID.String(),
			Icon: &profilev1.CreateUpdateIcon{
				Name:         icon.Name,
				SiteHostname: icon.SiteHostname,
				SiteSuffix:   utils.StringpbValue(siteSuffix),
				Hostname:     icon.Hostname,
				Suffix:       icon.Suffix,
				Color:        icon.Color,
			},
			UpdatedAt: timestamppb.New(icon.UpdatedAt),
		}

		res, err := rootProfileServiceClient.UpdateIcon(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, res)
	})
}
