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

func TestUpdatePlatform(t *testing.T) {
	t.Parallel()

	t.Run("Same Data", func(t *testing.T) {
		t.Parallel()

		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)

		platform := getPlatform(t)

		req := &profilev1.UpdatePlatformRequest{
			Id: platform.ID.String(),
			Platform: &profilev1.CreateUpdatePlatform{
				Name:          platform.Name,
				Url:           platform.Url,
				UrlSuffix:     platform.UrlSuffix,
				LogoUrl:       platform.LogoUrl,
				LogoUrlSuffix: platform.LogoUrlSuffix,
				LogoUrlPath:   platform.LogoUrlPath,
				Color:         platform.Color,
			},
			UpdatedAt: timestamppb.New(platform.UpdatedAt),
		}

		res, err := rootProfileServiceClient.UpdatePlatform(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrConflict, err)
	})

	t.Run("Unique Name Error", func(t *testing.T) {
		t.Parallel()

		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)

		platform := getPlatform(t)
		platform1 := getPlatform(t)

		req := &profilev1.UpdatePlatformRequest{
			Id: platform.ID.String(),
			Platform: &profilev1.CreateUpdatePlatform{
				Name:          platform1.Name,
				Url:           platform.Url,
				UrlSuffix:     platform.UrlSuffix,
				LogoUrl:       platform.LogoUrl,
				LogoUrlSuffix: platform.LogoUrlSuffix,
				LogoUrlPath:   platform.LogoUrlPath,
				Color:         platform.Color,
			},
			UpdatedAt: timestamppb.New(platform.UpdatedAt),
		}

		res, err := rootProfileServiceClient.UpdatePlatform(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrUniqueViolation("Platform"), err)
	})

	t.Run("Unique URL Error", func(t *testing.T) {
		t.Parallel()

		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)

		platform := getPlatform(t)
		platform1 := getPlatform(t)

		req := &profilev1.UpdatePlatformRequest{
			Id: platform.ID.String(),
			Platform: &profilev1.CreateUpdatePlatform{
				Name:          platform.Name,
				Url:           platform1.Url,
				UrlSuffix:     platform.UrlSuffix,
				LogoUrl:       platform.LogoUrl,
				LogoUrlSuffix: platform.LogoUrlSuffix,
				LogoUrlPath:   platform.LogoUrlPath,
				Color:         platform.Color,
			},
			UpdatedAt: timestamppb.New(platform.UpdatedAt),
		}

		res, err := rootProfileServiceClient.UpdatePlatform(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrUniqueViolation("URL"), err)
	})

	t.Run("Unique Logo URL Error", func(t *testing.T) {
		t.Parallel()

		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)

		platform := getPlatform(t)
		platform1 := getPlatform(t)

		req := &profilev1.UpdatePlatformRequest{
			Id: platform.ID.String(),
			Platform: &profilev1.CreateUpdatePlatform{
				Name:          platform.Name,
				Url:           platform.Url,
				UrlSuffix:     platform.UrlSuffix,
				LogoUrl:       platform1.LogoUrl,
				LogoUrlSuffix: platform1.LogoUrlSuffix,
				LogoUrlPath:   platform1.LogoUrlPath,
				Color:         platform.Color,
			},
			UpdatedAt: timestamppb.New(platform.UpdatedAt),
		}

		res, err := rootProfileServiceClient.UpdatePlatform(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrUniqueViolation("Logo URL"), err)
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)

		platform := getPlatform(t)
		name := strings.ToLower(rand.Text()[:8])

		req := &profilev1.UpdatePlatformRequest{
			Id: platform.ID.String(),
			Platform: &profilev1.CreateUpdatePlatform{
				Name:          name,
				Url:           platform.Url,
				UrlSuffix:     platform.UrlSuffix,
				LogoUrl:       platform.LogoUrl,
				LogoUrlSuffix: platform.LogoUrlSuffix,
				LogoUrlPath:   platform.LogoUrlPath,
				Color:         platform.Color,
			},
			UpdatedAt: timestamppb.New(platform.UpdatedAt),
		}

		res, err := rootProfileServiceClient.UpdatePlatform(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, platform.ID.String(), res.GetId())
	})
}
