//go:build integration

package service_test

import (
	"testing"
	"time"
	"uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"neupaneanish.com.np/profile/internal/enum"
	"neupaneanish.com.np/profile/internal/errs"
	profilev1 "neupaneanish.com.np/profile/internal/protobuf/common/profile/v1"
	gatewayProfilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
)

func TestGatewayUpdateProfile(t *testing.T) {
	t.Parallel()
	t.Run("Error", func(t *testing.T) {
		t.Parallel()
		profile := getProfile(t, uuid.NewV7())
		gatewayUpdateProfileErr(t, profile.UserID, profile.Name, profile.Title, profile.UpdatedAt)
	})

	t.Run("Different UpdateAt", func(t *testing.T) {
		t.Parallel()
		profile := getProfile(t, uuid.NewV7())
		gatewayUpdateProfileErr(t, profile.UserID, "Neupane Anish", "Founder CEO", time.Now())
	})

	t.Run("No User", func(t *testing.T) {
		t.Parallel()
		gatewayUpdateProfileErr(t, uuid.NewV7(), "Neupane Anish", "Founder CEO", time.Now())
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		profile := getProfile(t, uuid.NewV7())
		ctx := contextWithValue(t, profile.UserID, enum.UserRoleUser)

		name := "Neupane Anish"

		req := &gatewayProfilev1.UpdateProfileRequest{
			Profile: &profilev1.CreateUpdateProfile{
				Name:  name,
				Title: profile.Title,
			},
			UpdatedAt: timestamppb.New(profile.UpdatedAt),
		}

		res, err := gatewayProfileServiceClient.UpdateProfile(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, name, res.GetProfile().GetName())
	})
}

func TestRootUpdateProfile(t *testing.T) {
	t.Parallel()
	t.Run("Same Data Error", func(t *testing.T) {
		t.Parallel()
		profile := getProfile(t, uuid.NewV7())
		rootUpdateProfileErr(t, uuid.NewV7(), profile.UserID, profile.Name, profile.Title, profile.UpdatedAt)
	})

	t.Run("Different UpdateAt", func(t *testing.T) {
		t.Parallel()
		profile := getProfile(t, uuid.NewV7())
		rootUpdateProfileErr(t, uuid.NewV7(), profile.UserID, "Neupane Anish", "Founder COO", time.Now())
	})

	t.Run("No User", func(t *testing.T) {
		t.Parallel()
		rootUpdateProfileErr(t, uuid.NewV7(), uuid.NewV7(), "Neupane Anish", "Founder COO", time.Now())
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		profile := getProfile(t, uuid.NewV7())

		adminID := uuid.NewV7()

		ctx := contextWithValue(t, adminID, enum.UserRoleRoot)

		name := "Neupane Anish"

		req := &rootProfilev1.UpdateProfileRequest{
			UserId: profile.UserID.String(),
			Profile: &profilev1.CreateUpdateProfile{
				Name:  name,
				Title: profile.Title,
			},
			UpdatedAt: timestamppb.New(profile.UpdatedAt),
		}

		res, err := rootProfileServiceClient.UpdateProfile(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, name, res.GetProfile().GetName())
		assert.Equal(t, adminID.String(), res.GetProfile().GetUpdatedBy())
	})
}

func gatewayUpdateProfileErr(t *testing.T, userID uuid.UUID, name, title string, updatedAt time.Time) {
	ctx := contextWithValue(t, userID, enum.UserRoleUser)

	req := &gatewayProfilev1.UpdateProfileRequest{
		Profile: &profilev1.CreateUpdateProfile{
			Name:  name,
			Title: title,
		},
		UpdatedAt: timestamppb.New(updatedAt),
	}

	res, err := gatewayProfileServiceClient.UpdateProfile(ctx, req)
	require.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, errs.ErrConflict, err)
}

func rootUpdateProfileErr(t *testing.T, adminID, userID uuid.UUID, name, title string, updatedAt time.Time) {
	ctx := contextWithValue(t, adminID, enum.UserRoleRoot)

	req := &rootProfilev1.UpdateProfileRequest{
		UserId: userID.String(),
		Profile: &profilev1.CreateUpdateProfile{
			Name:  name,
			Title: title,
		},
		UpdatedAt: timestamppb.New(updatedAt),
	}

	res, err := rootProfileServiceClient.UpdateProfile(ctx, req)
	require.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, errs.ErrConflict, err)
}
