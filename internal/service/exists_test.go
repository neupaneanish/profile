//go:build integration

package service_test

import (
	"crypto/rand"
	"testing"
	"uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	externalProfilev1 "neupaneanish.com.np/profile/internal/protobuf/external/profile/v1"
)

func TestExists(t *testing.T) {
	t.Parallel()
	userID := uuid.NewV7()

	ctx := externalContextWithValue(t, userID, rand.Text()[:8]+".com")
	req := &externalProfilev1.ExistsRequest{}
	res, err := externalProfileServiceClient.Exists(ctx, req)
	require.NoError(t, err)
	assert.NotNil(t, res)
	assert.False(t, res.GetExists().GetSocials())
}
