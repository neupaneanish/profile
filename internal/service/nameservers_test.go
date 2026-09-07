//go:build integration

package service_test

import (
	"crypto/rand"
	"strings"
	"testing"
	"uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"neupaneanish.com.np/profile/internal/enum"
	profilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
)

func TestNameservers(t *testing.T) {
	t.Parallel()

	seedNameServer(t, strings.ToLower(rand.Text()[:8]))
	seedNameServer(t, strings.ToLower(rand.Text()[:8]))
	seedNameServer(t, strings.ToLower(rand.Text()[:8]))

	ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)

	req := &profilev1.NameServersRequest{}

	res, err := rootProfileServiceClient.NameServers(ctx, req)
	require.NoError(t, err)
	assert.NotNil(t, res)
	assert.GreaterOrEqual(t, len(res.GetNameservers()), 3)
}
