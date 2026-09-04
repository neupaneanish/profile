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
	"neupaneanish.com.np/profile/internal/service"
)

func TestGatewayUpdateExperience(t *testing.T) {
	t.Parallel()

	t.Run("No Update Same Data", func(t *testing.T) {
		t.Parallel()
		userID := uuid.NewV7()
		ctx := contextWithValue(t, userID, enum.UserRoleUser)
		experience := getExperience(t, userID)

		params := &profilev1.UpdateExperience{
			Id:           experience.ID.String(),
			Title:        experience.Title,
			CompanyName:  experience.CompanyName,
			Location:     experience.Location,
			LocationType: string(experience.LocationType),
			StartDate:    timestamppb.New(experience.StartDate),
			EndDate:      service.TimestamppbValue(experience.EndDate),
			Description:  service.StringpbValue(experience.Description),
			UpdatedAt:    timestamppb.New(experience.UpdatedAt),
		}

		res, err := gatewayProfileServiceClient.UpdateExperience(
			ctx, &gatewayProfilev1.UpdateExperienceRequest{
				Experience: params,
			},
		)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrConflict, err)
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		userID := uuid.NewV7()
		ctx := contextWithValue(t, userID, enum.UserRoleUser)
		experience := getExperience(t, userID)

		params := &profilev1.UpdateExperience{
			Id:           experience.ID.String(),
			Title:        "Professor",
			CompanyName:  experience.CompanyName,
			Location:     experience.Location,
			LocationType: string(experience.LocationType),
			StartDate:    timestamppb.New(experience.StartDate),
			EndDate:      service.TimestamppbValue(experience.EndDate),
			Description:  service.StringpbValue(experience.Description),
			UpdatedAt:    timestamppb.New(experience.UpdatedAt),
		}

		res, err := gatewayProfileServiceClient.UpdateExperience(
			ctx, &gatewayProfilev1.UpdateExperienceRequest{
				Experience: params,
			},
		)
		require.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, experience.ID.String(), res.GetId())
	})
}

func TestRootUpdateExperience(t *testing.T) {
	t.Parallel()

	t.Run("Invalid Updated At", func(t *testing.T) {
		t.Parallel()
		userID := uuid.NewV7()
		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)
		experience := getExperience(t, userID)

		params := &profilev1.UpdateExperience{
			Id:           experience.ID.String(),
			Title:        experience.Title,
			CompanyName:  experience.CompanyName,
			Location:     experience.Location,
			LocationType: string(experience.LocationType),
			StartDate:    timestamppb.New(experience.StartDate),
			EndDate:      service.TimestamppbValue(experience.EndDate),
			Description:  service.StringpbValue(experience.Description),
			UpdatedAt:    timestamppb.New(time.Now()),
		}

		res, err := rootProfileServiceClient.UpdateExperience(
			ctx, &rootProfilev1.UpdateExperienceRequest{
				UserId:     userID.String(),
				Experience: params,
			},
		)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrConflict, err)
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		userID := uuid.NewV7()
		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)
		experience := getExperience(t, userID)

		params := &profilev1.UpdateExperience{
			Id:           experience.ID.String(),
			Title:        "Professor",
			CompanyName:  experience.CompanyName,
			Location:     experience.Location,
			LocationType: string(experience.LocationType),
			StartDate:    timestamppb.New(experience.StartDate),
			EndDate:      service.TimestamppbValue(experience.EndDate),
			Description:  service.StringpbValue(experience.Description),
			UpdatedAt:    timestamppb.New(experience.UpdatedAt),
		}

		res, err := rootProfileServiceClient.UpdateExperience(
			ctx, &rootProfilev1.UpdateExperienceRequest{
				UserId:     userID.String(),
				Experience: params,
			},
		)
		require.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, experience.ID.String(), res.GetId())
	})
}
