//go:build integration

package service_test

import (
	"testing"
	"uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"neupaneanish.com.np/profile/internal/enum"
	profilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
)

func TestNameservers(t *testing.T) {
	t.Parallel()

	ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)

	req := &profilev1.NameserversRequest{}

	res, err := rootProfileServiceClient.Nameservers(ctx, req)
	require.NoError(t, err)
	assert.NotNil(t, res)
}
