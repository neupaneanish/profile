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
	gatewayProfilev1 "neupaneanish.com.np/profile/internal/protobuf/gateway/profile/v1"
)

func TestDeleteEducation(t *testing.T) {
	t.Parallel()

	t.Run("No Education", func(t *testing.T) {
		t.Parallel()
		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleUser)
		req := &gatewayProfilev1.DeleteEducationRequest{Id: uuid.NewV7().String(), UpdatedAt: timestamppb.New(time.Now())}

		res, err := gatewayProfileServiceClient.DeleteEducation(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrConflict, err)
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		userID := uuid.NewV7()
		ctx := contextWithValue(t, userID, enum.UserRoleUser)
		education := getEducation(t, userID)

		req := &gatewayProfilev1.DeleteEducationRequest{Id: education.ID.String(), UpdatedAt: timestamppb.New(education.UpdatedAt)}
		res, err := gatewayProfileServiceClient.DeleteEducation(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, res)
	})
}
