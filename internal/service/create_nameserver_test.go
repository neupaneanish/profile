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
	"neupaneanish.com.np/profile/internal/errs"
	profilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
	"neupaneanish.com.np/profile/internal/repository"
)

func TestCreateNameServer(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		name := strings.ToLower(rand.Text()[:8])
		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)

		req := &profilev1.CreateNameServerRequest{
			Nameserver: &profilev1.CreateUpdateNameServer{
				Domain: "https://" + name + ".com",
				Cname:  name,
			},
		}

		res, err := rootProfileServiceClient.CreateNameServer(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, res)
		assert.NotEmpty(t, res.GetId())
	})

	t.Run("Invalid Domain", func(t *testing.T) {
		t.Parallel()

		name := strings.ToLower(rand.Text()[:8])
		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)

		req := &profilev1.CreateNameServerRequest{
			Nameserver: &profilev1.CreateUpdateNameServer{
				Domain: "https://127.0.0.1",
				Cname:  name,
			},
		}

		res, err := rootProfileServiceClient.CreateNameServer(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrInvalidURL, err)
	})

	t.Run("UniqueViolation", func(t *testing.T) {
		t.Parallel()

		name := strings.ToLower(rand.Text()[:8])

		seedNameServer(t, name)

		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)

		req := &profilev1.CreateNameServerRequest{
			Nameserver: &profilev1.CreateUpdateNameServer{
				Domain: "https://" + name + ".com",
				Cname:  name,
			},
		}

		res, err := rootProfileServiceClient.CreateNameServer(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrUniqueViolation("Nameserver"), err)
	})
}

func seedNameServer(t *testing.T, name string) {
	t.Helper()

	userID := uuid.NewV7()
	params := &repository.CreateNameServerParams{
		Domain:    name + ".com",
		Cname:     name,
		CreatedBy: userID,
		UpdatedBy: userID,
	}

	_, err := cfg.Repository.CreateNameServer(t.Context(), params)
	require.NoError(t, err)
}

func getNameServer(t *testing.T, name string) *repository.Nameserver {
	t.Helper()
	seedNameServer(t, name)

	rows, err := cfg.Repository.NameServers(t.Context())
	require.NoError(t, err)

	return rows[0]
}
