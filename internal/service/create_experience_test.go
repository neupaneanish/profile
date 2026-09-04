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
	gatewayProfilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
)

func TestCreateExperience(t *testing.T) {
	t.Parallel()
	ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleUser)

	req := &gatewayProfilev1.CreateExperienceRequest{
		Title:        "Teacher",
		CompanyName:  "Oxford College of Engineering and Management",
		Location:     "Nepal",
		LocationType: "hybrid",
		StartDate:    timestamppb.New(time.Date(2023, time.June, 15, 0, 0, 0, 0, time.UTC)),
		EndDate:      timestamppb.New(time.Date(2024, time.August, 14, 0, 0, 0, 0, time.UTC)),
		Description:  &wrapperspb.StringValue{Value: "Computer Science Teacher"},
	}

	res, err := gatewayProfileServiceClient.CreateExperience(ctx, req)
	require.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotEmpty(t, res.GetId())
}
