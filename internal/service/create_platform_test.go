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
	"neupaneanish.com.np/profile/internal/repository"
)

func TestCreatePlatform(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		name := strings.ToLower(rand.Text()[:8])
		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)

		req := &profilev1.CreatePlatformRequest{
			Platform: &profilev1.CreateUpdatePlatform{
				Name:          name,
				Url:           "https://" + name + ".com",
				UrlSuffix:     "/",
				LogoUrl:       "https://" + name + ".com",
				LogoUrlSuffix: "/" + name,
				LogoUrlPath:   name,
				Color:         "#FFFFFF",
			},
		}

		res, err := rootProfileServiceClient.CreatePlatform(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, res)
		assert.NotEmpty(t, res.GetId())
	})

	t.Run("Unique Name Error", func(t *testing.T) {
		t.Parallel()

		name := strings.ToLower(rand.Text()[:8])

		platform := getPlatform(t)

		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)

		req := &profilev1.CreatePlatformRequest{
			Platform: &profilev1.CreateUpdatePlatform{
				Name:          platform.Name,
				Url:           "https://" + name + ".com",
				UrlSuffix:     "/",
				LogoUrl:       "https://" + name + ".com",
				LogoUrlSuffix: "/" + name,
				LogoUrlPath:   name,
				Color:         "#FFFFFF",
			},
		}

		res, err := rootProfileServiceClient.CreatePlatform(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrUniqueViolation("Platform"), err)
	})

	t.Run("Unique URL Error", func(t *testing.T) {
		t.Parallel()

		name := strings.ToLower(rand.Text()[:8])

		platform := getPlatform(t)

		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)

		req := &profilev1.CreatePlatformRequest{
			Platform: &profilev1.CreateUpdatePlatform{
				Name:          name,
				Url:           platform.Url,
				UrlSuffix:     "/",
				LogoUrl:       "https://" + name + ".com",
				LogoUrlSuffix: "/" + name,
				LogoUrlPath:   name,
				Color:         "#FFFFFF",
			},
		}

		res, err := rootProfileServiceClient.CreatePlatform(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrUniqueViolation("URL"), err)
	})

	t.Run("Unique Logo URL Error", func(t *testing.T) {
		t.Parallel()

		name := strings.ToLower(rand.Text()[:8])

		platform := getPlatform(t)

		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)

		req := &profilev1.CreatePlatformRequest{
			Platform: &profilev1.CreateUpdatePlatform{
				Name:          name,
				Url:           "https://" + name + ".com",
				UrlSuffix:     "/",
				LogoUrl:       platform.LogoUrl,
				LogoUrlSuffix: platform.LogoUrlSuffix,
				LogoUrlPath:   platform.LogoUrlPath,
				Color:         "#FFFFFF",
			},
		}

		res, err := rootProfileServiceClient.CreatePlatform(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrUniqueViolation("Logo URL"), err)
	})

	t.Run("Invalid URL", func(t *testing.T) {
		t.Parallel()

		name := strings.ToLower(rand.Text()[:8])
		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)

		req := &profilev1.CreatePlatformRequest{
			Platform: &profilev1.CreateUpdatePlatform{
				Name:          name,
				Url:           "https://127.0.0.1",
				UrlSuffix:     "/",
				LogoUrl:       "https://" + name + ".com",
				LogoUrlSuffix: "/" + name,
				LogoUrlPath:   name,
				Color:         "#FFFFFF",
			},
		}

		res, err := rootProfileServiceClient.CreatePlatform(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrInvalidURL, err)
	})

	t.Run("Invalid Logo URL", func(t *testing.T) {
		t.Parallel()

		name := strings.ToLower(rand.Text()[:8])
		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)

		req := &profilev1.CreatePlatformRequest{
			Platform: &profilev1.CreateUpdatePlatform{
				Name:          name,
				Url:           "https://" + name + ".com",
				UrlSuffix:     "/",
				LogoUrl:       "https://127.0.0.1",
				LogoUrlSuffix: "/" + name,
				LogoUrlPath:   name,
				Color:         "#FFFFFF",
			},
		}

		res, err := rootProfileServiceClient.CreatePlatform(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrInvalidURL, err)
	})

	t.Run("Localhost", func(t *testing.T) {
		t.Parallel()

		name := strings.ToLower(rand.Text()[:8])
		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)

		req := &profilev1.CreatePlatformRequest{
			Platform: &profilev1.CreateUpdatePlatform{
				Name:          name,
				Url:           "https://localhost",
				UrlSuffix:     "/",
				LogoUrl:       "https://127.0.0.1",
				LogoUrlSuffix: "/" + name,
				LogoUrlPath:   name,
				Color:         "#FFFFFF",
			},
		}

		res, err := rootProfileServiceClient.CreatePlatform(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrInvalidURL, err)
	})
}

func seedPlatform(t *testing.T) uuid.UUID {
	name := strings.ToLower(rand.Text()[:8])
	params := &repository.CreatePlatformParams{
		Name:          name,
		Url:           "https://" + name + ".com",
		UrlSuffix:     "/",
		LogoUrl:       "https://" + name + ".com",
		LogoUrlSuffix: "/" + name,
		LogoUrlPath:   name,
		Color:         "#FFFFFF",
		CreatedBy:     uuid.NewV7(),
		UpdatedBy:     uuid.NewV7(),
	}
	id, err := cfg.Repository.CreatePlatform(t.Context(), params)
	require.NoError(t, err)

	return id
}

func getPlatform(t *testing.T) *repository.Platform {
	id := seedPlatform(t)
	params := &repository.PlatformParams{ID: id}
	platform, err := cfg.Repository.Platform(t.Context(), params)
	require.NoError(t, err)

	return platform
}
