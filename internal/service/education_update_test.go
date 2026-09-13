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

func TestUpdateEducation(t *testing.T) {
	t.Parallel()

	t.Run("Error gateway", func(t *testing.T) {
		t.Parallel()
		userID := uuid.NewV7()
		ctx := contextWithValue(t, userID, enum.UserRoleUser)

		id, params, updatedAt := updateEducation(t, userID, "")

		res, err := gatewayProfileServiceClient.UpdateEducation(
			ctx, &gatewayProfilev1.UpdateEducationRequest{
				Id:        id.String(),
				Education: params,
				UpdatedAt: timestamppb.New(updatedAt),
			},
		)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrConflict, err)
	})

	t.Run("Success Success", func(t *testing.T) {
		t.Parallel()
		userID := uuid.NewV7()
		ctx := contextWithValue(t, userID, enum.UserRoleUser)

		id, params, updatedAt := updateEducation(t, userID, "Test School")

		res, err := gatewayProfileServiceClient.UpdateEducation(
			ctx, &gatewayProfilev1.UpdateEducationRequest{
				Id:        id.String(),
				Education: params,
				UpdatedAt: timestamppb.New(updatedAt),
			},
		)
		require.NoError(t, err)
		assert.NotNil(t, res)
	})

	t.Run("Error Root", func(t *testing.T) {
		t.Parallel()
		userID := uuid.NewV7()
		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)

		id, params, updatedAt := updateEducation(t, userID, "")

		res, err := rootProfileServiceClient.UpdateEducation(
			ctx, &rootProfilev1.UpdateEducationRequest{
				UserId:    userID.String(),
				Education: params,
				Id:        id.String(),
				UpdatedAt: timestamppb.New(updatedAt),
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

		id, params, updatedAt := updateEducation(t, userID, "Test School")

		res, err := rootProfileServiceClient.UpdateEducation(
			ctx, &rootProfilev1.UpdateEducationRequest{
				Id:        id.String(),
				UserId:    userID.String(),
				Education: params,
				UpdatedAt: timestamppb.New(updatedAt),
			},
		)
		require.NoError(t, err)
		assert.NotNil(t, res)
	})
}

func updateEducation(
	t *testing.T,
	userID uuid.UUID,
	school string,
) (uuid.UUID, *profilev1.CreateUpdateEducation, time.Time) {
	education := getEducation(t, userID)

	var value string

	if school == "" {
		value = education.School
	} else {
		value = school
	}
	params := &profilev1.CreateUpdateEducation{
		School:        value,
		Degree:        education.Degree,
		Affiliation:   utils.StringpbValue(education.Affiliation),
		FieldOfStudy:  utils.StringpbValue(education.FieldOfStudy),
		Concentration: utils.StringpbValue(education.Concentration),
		StartDate:     timestamppb.New(education.StartDate),
		EndDate:       utils.TimestamppbValue(education.EndDate),
		Address:       education.Address,
		Description:   utils.StringpbValue(education.Description),
	}
	return education.ID, params, education.UpdatedAt
}
