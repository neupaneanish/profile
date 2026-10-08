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
	"neupaneanish.com.np/profile/internal/utils"
)

func TestCreateIcon(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		name := strings.ToLower(rand.Text()[:8])
		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)

		req := &profilev1.CreateIconRequest{
			Icon: &profilev1.CreateUpdateIcon{
				Name:         name,
				SiteHostname: name + ".com",
				Hostname:     name + ".com",
				Suffix:       name,
				Color:        "#FFFFFF",
			},
		}

		res, err := rootProfileServiceClient.CreateIcon(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, res)
	})

	t.Run("Unique Name Error", func(t *testing.T) {
		t.Parallel()

		name := strings.ToLower(rand.Text()[:8])

		icon := getIcon(t, name+name, nil)

		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)

		req := &profilev1.CreateIconRequest{
			Icon: &profilev1.CreateUpdateIcon{
				Name:         icon.Name,
				SiteHostname: name + ".com",
				Hostname:     name + ".com",
				Suffix:       name,
				Color:        "#FFFFFF",
			},
		}

		res, err := rootProfileServiceClient.CreateIcon(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrUniqueViolation("Icon Name"), err)
	})

	t.Run("Unique Icon Site", func(t *testing.T) {
		t.Parallel()

		name := strings.ToLower(rand.Text()[:8])

		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)
		icon := getIcon(t, name+name, nil)

		req := &profilev1.CreateIconRequest{
			Icon: &profilev1.CreateUpdateIcon{
				Name:         name,
				SiteHostname: icon.SiteHostname,
				Hostname:     name + ".com",
				Suffix:       name,
				Color:        "#FFFFFF",
			},
		}
		res, err := rootProfileServiceClient.CreateIcon(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrUniqueViolation("Icon Site"), err)
	})

	t.Run("Unique Icon URL", func(t *testing.T) {
		t.Parallel()

		name := strings.ToLower(rand.Text()[:8])

		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)
		icon := getIcon(t, name+name, nil)

		req := &profilev1.CreateIconRequest{
			Icon: &profilev1.CreateUpdateIcon{
				Name:         name,
				SiteHostname: name + ".com",
				Hostname:     icon.Hostname,
				Suffix:       icon.Suffix,
				Color:        "#FFFFFF",
			},
		}

		res, err := rootProfileServiceClient.CreateIcon(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrUniqueViolation("Icon URL"), err)
	})

	t.Run("Icon Site with suffix", func(t *testing.T) {
		t.Parallel()

		name := strings.ToLower(rand.Text()[:8])
		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)

		suffix := "/" + name

		icon := getIcon(t, name+name, &suffix)

		req := &profilev1.CreateIconRequest{
			Icon: &profilev1.CreateUpdateIcon{
				Name:         name,
				SiteHostname: icon.SiteHostname,
				SiteSuffix:   utils.StringpbValue(icon.SiteSuffix),
				Hostname:     name + ".com",
				Suffix:       name,
				Color:        "#FFFFFF",
			},
		}

		res, err := rootProfileServiceClient.CreateIcon(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrUniqueViolation("Icon Site with suffix"), err)
	})

	t.Run("Invalid Site subdomain not allowed", func(t *testing.T) {
		t.Parallel()

		name := strings.ToLower(rand.Text()[:8])
		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)

		req := &profilev1.CreateIconRequest{
			Icon: &profilev1.CreateUpdateIcon{
				Name:         name,
				SiteHostname: "sub." + name + ".com",
				Hostname:     name + ".com",
				Suffix:       name,
				Color:        "#FFFFFF",
			},
		}

		res, err := rootProfileServiceClient.CreateIcon(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrInvalidURL, err)
	})

	t.Run("Invalid public suffix", func(t *testing.T) {
		t.Parallel()

		name := strings.ToLower(rand.Text()[:8])
		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)

		req := &profilev1.CreateIconRequest{
			Icon: &profilev1.CreateUpdateIcon{
				Name:         name,
				SiteHostname: name + ".com",
				Hostname:     name + ".amc",
				Suffix:       name,
				Color:        "#FFFFFF",
			},
		}

		res, err := rootProfileServiceClient.CreateIcon(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrInvalidURL, err)
	})

	t.Run("Invalid TLDPlusOne", func(t *testing.T) {
		t.Parallel()

		name := strings.ToLower(rand.Text()[:8])
		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)

		req := &profilev1.CreateIconRequest{
			Icon: &profilev1.CreateUpdateIcon{
				Name:         name,
				SiteHostname: "abc.cde" + name + ".com",
				Hostname:     name + ".com",
				Suffix:       name,
				Color:        "#FFFFFF",
			},
		}

		res, err := rootProfileServiceClient.CreateIcon(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrInvalidURL, err)
	})
}
