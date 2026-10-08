//go:build integration

package service_test

import (
	"crypto/rand"
	"testing"
	"uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"neupaneanish.com.np/profile/internal/enum"
	"neupaneanish.com.np/profile/internal/errs"
	profilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
)

func TestCreateNameServer(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)

		name := rand.Text()[:8]

		req := &profilev1.CreateNameserverRequest{
			Cname:    getCname(),
			Hostname: name + ".com",
		}

		res, err := rootProfileServiceClient.CreateNameserver(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, res)
	})

	t.Run("Invalid Hostname", func(t *testing.T) {
		t.Parallel()
		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)

		name := rand.Text()[:8]

		req := &profilev1.CreateNameserverRequest{
			Cname:    getCname(),
			Hostname: "sub" + name + ".akjchn",
		}

		res, err := rootProfileServiceClient.CreateNameserver(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrInvalidHostname, err)
	})

	t.Run("UniqueViolation", func(t *testing.T) {
		t.Parallel()

		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)

		name := rand.Text()[:8]

		cname := getCname()
		hostname := name + ".com"

		seedNameserver(t, cname, hostname)

		req := &profilev1.CreateNameserverRequest{
			Cname:    cname,
			Hostname: hostname,
		}

		res, err := rootProfileServiceClient.CreateNameserver(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrUniqueViolation("Nameserver"), err)
	})
}
