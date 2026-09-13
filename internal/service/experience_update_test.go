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
	"neupaneanish.com.np/profile/internal/utils"
)

func TestUpdateExperience(t *testing.T) {
	t.Parallel()

	t.Run("Error Gateway", func(t *testing.T) {
		t.Parallel()
		userID := uuid.NewV7()
		ctx := contextWithValue(t, userID, enum.UserRoleUser)
		id, params, updatedAt := updateExperience(t, userID, "")

		res, err := gatewayProfileServiceClient.UpdateExperience(
			ctx, &gatewayProfilev1.UpdateExperienceRequest{
				Id:         id.String(),
				Experience: params,
				UpdatedAt:  timestamppb.New(updatedAt),
			},
		)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrConflict, err)
	})

	t.Run("Success Gateway", func(t *testing.T) {
		t.Parallel()
		userID := uuid.NewV7()
		ctx := contextWithValue(t, userID, enum.UserRoleUser)
		id, params, updatedAt := updateExperience(t, userID, "Test Title")

		res, err := gatewayProfileServiceClient.UpdateExperience(
			ctx, &gatewayProfilev1.UpdateExperienceRequest{
				Id:         id.String(),
				Experience: params,
				UpdatedAt:  timestamppb.New(updatedAt),
			},
		)
		require.NoError(t, err)
		assert.NotNil(t, res)
	})

	t.Run("Error Root", func(t *testing.T) {
		t.Parallel()
		userID := uuid.NewV7()
		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)

		id, params, updatedAt := updateExperience(t, userID, "")

		res, err := rootProfileServiceClient.UpdateExperience(
			ctx, &rootProfilev1.UpdateExperienceRequest{
				Id:         id.String(),
				UserId:     userID.String(),
				Experience: params,
				UpdatedAt:  timestamppb.New(updatedAt),
			},
		)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrConflict, err)
	})

	t.Run("Success Root", func(t *testing.T) {
		t.Parallel()
		userID := uuid.NewV7()
		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)

		id, params, updatedAt := updateExperience(t, userID, "Test Title")

		res, err := rootProfileServiceClient.UpdateExperience(
			ctx, &rootProfilev1.UpdateExperienceRequest{
				Id:         id.String(),
				UserId:     userID.String(),
				Experience: params,
				UpdatedAt:  timestamppb.New(updatedAt),
			},
		)
		require.NoError(t, err)
		assert.NotNil(t, res)
	})
}

func updateExperience(
	t *testing.T,
	userID uuid.UUID,
	title string,
) (uuid.UUID, *profilev1.CreateUpdateExperience, time.Time) {
	experience := getExperience(t, userID)

	var value string

	if title == "" {
		value = experience.Title
	} else {
		value = title
	}
	params := &profilev1.CreateUpdateExperience{
		Title:        value,
		CompanyName:  experience.CompanyName,
		Location:     experience.Location,
		LocationType: string(experience.LocationType),
		StartDate:    timestamppb.New(experience.StartDate),
		EndDate:      utils.TimestamppbValue(experience.EndDate),
		Description:  utils.StringpbValue(experience.Description),
	}
	return experience.ID, params, experience.UpdatedAt
}
