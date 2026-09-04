//go:build integration

package service_test

import (
	"testing"
	"time"
	"uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"neupaneanish.com.np/profile/internal/enum"
	"neupaneanish.com.np/profile/internal/errs"
	profilev1 "neupaneanish.com.np/profile/internal/protobuf/common/profile/v1"
	gatewayProfilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
	rootProfilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
)

func TestGatewayUpdateEducation(t *testing.T) {
	t.Parallel()

	t.Run("UpdateAt Mismatch", func(t *testing.T) {
		t.Parallel()
		userID := uuid.NewV7()
		ctx := contextWithValue(t, userID, enum.UserRoleUser)
		education := getEducation(t, userID)

		params := &profilev1.UpdateEducation{
			Id:            education.ID.String(),
			School:        education.School,
			Degree:        education.Degree,
			Affiliation:   nil,
			FieldOfStudy:  nil,
			Concentration: nil,
			StartDate:     timestamppb.New(education.StartDate),
			EndDate:       nil,
			Address:       education.Address,
			Description:   nil,
			UpdatedAt:     timestamppb.New(time.Now()),
		}

		res, err := gatewayProfileServiceClient.UpdateEducation(
			ctx, &gatewayProfilev1.UpdateEducationRequest{
				Education: params,
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
		education := getEducation(t, userID)

		params := &profilev1.UpdateEducation{
			Id:            education.ID.String(),
			School:        education.School,
			Degree:        education.Degree,
			Affiliation:   nil,
			FieldOfStudy:  nil,
			Concentration: nil,
			StartDate:     timestamppb.New(education.StartDate),
			EndDate:       timestamppb.New(time.Date(2026, time.August, 31, 0, 0, 0, 0, time.UTC)),
			Address:       education.Address,
			Description:   nil,
			UpdatedAt:     timestamppb.New(education.UpdatedAt),
		}

		res, err := gatewayProfileServiceClient.UpdateEducation(
			ctx, &gatewayProfilev1.UpdateEducationRequest{
				Education: params,
			},
		)
		require.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, education.ID.String(), res.GetId())
	})
}

func TestRootUpdateEducation(t *testing.T) {
	t.Parallel()

	t.Run("No Update Same Data", func(t *testing.T) {
		t.Parallel()
		userID := uuid.NewV7()
		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)
		education := getEducation(t, userID)

		params := &profilev1.UpdateEducation{
			Id:            education.ID.String(),
			School:        education.School,
			Degree:        education.Degree,
			Affiliation:   nil,
			FieldOfStudy:  nil,
			Concentration: nil,
			StartDate:     timestamppb.New(education.StartDate),
			EndDate:       nil,
			Address:       education.Address,
			Description:   nil,
			UpdatedAt:     timestamppb.New(education.UpdatedAt),
		}

		res, err := rootProfileServiceClient.UpdateEducation(
			ctx, &rootProfilev1.UpdateEducationRequest{
				UserId:    userID.String(),
				Education: params,
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
		education := getEducation(t, userID)

		params := &profilev1.UpdateEducation{
			Id:            education.ID.String(),
			School:        education.School,
			Degree:        education.Degree,
			Affiliation:   nil,
			FieldOfStudy:  nil,
			Concentration: &wrapperspb.StringValue{Value: "Artificial Intelligence"},
			StartDate:     timestamppb.New(education.StartDate),
			EndDate:       timestamppb.New(time.Date(2026, time.August, 31, 0, 0, 0, 0, time.UTC)),
			Address:       education.Address,
			Description:   nil,
			UpdatedAt:     timestamppb.New(education.UpdatedAt),
		}

		res, err := rootProfileServiceClient.UpdateEducation(
			ctx, &rootProfilev1.UpdateEducationRequest{
				UserId:    userID.String(),
				Education: params,
			},
		)
		require.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, education.ID.String(), res.GetId())
	})
}
