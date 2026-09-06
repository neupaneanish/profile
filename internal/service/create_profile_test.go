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
)

func TestCreateProfile(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleUser)
		req := &gatewayProfilev1.CreateProfileRequest{
			Profile: &profilev1.CreateUpdateProfile{
				Name:  "Anish Neupane",
				Title: "Backend Developer",
			},
			Dob: timestamppb.New(time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC)),
		}

		res, err := gatewayProfileServiceClient.CreateProfile(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, res)
	})

	t.Run("Already Exists", func(t *testing.T) {
		t.Parallel()
		userID := uuid.NewV7()

		seedProfile(t, userID)

		ctx := contextWithValue(t, userID, enum.UserRoleUser)
		req := &gatewayProfilev1.CreateProfileRequest{
			Profile: &profilev1.CreateUpdateProfile{
				Name:  "Anish Neupane",
				Title: "Backend Developer",
			},
			Dob: timestamppb.New(time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC)),
		}

		res, err := gatewayProfileServiceClient.CreateProfile(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrProfileAlreadyExists, err)
	})
}
