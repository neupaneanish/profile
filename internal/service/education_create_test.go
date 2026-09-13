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
	profilev1 "neupaneanish.com.np/profile/internal/protobuf/common/profile/v1"
	gatewayProfilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
)

func TestCreateEducation(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleUser)

		req := &gatewayProfilev1.CreateEducationRequest{
			Education: &profilev1.CreateUpdateEducation{
				School:        "Westcliff University",
				Degree:        "Master of Science in Computer Science",
				Affiliation:   nil,
				FieldOfStudy:  nil,
				Concentration: nil,
				StartDate:     timestamppb.New(time.Date(2024, time.September, 1, 0, 0, 0, 0, time.UTC)),
				EndDate:       timestamppb.New(time.Date(2026, time.August, 31, 0, 0, 0, 0, time.UTC)),
				Address:       "United State",
				Description:   nil,
			},
		}

		res, err := gatewayProfileServiceClient.CreateEducation(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, res)
	})
}
