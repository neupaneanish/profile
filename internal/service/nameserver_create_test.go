//go:build integration

package service_test

import (
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
		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)

		req := &profilev1.CreateNameserverRequest{
			IpType: "A",
			Ip:     "183.54.120.1",
		}

		res, err := rootProfileServiceClient.CreateNameserver(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, res)
	})

	t.Run("Private IP", func(t *testing.T) {
		t.Parallel()
		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)

		req := &profilev1.CreateNameserverRequest{
			IpType: "A",
			Ip:     "192.168.1.1",
		}

		res, err := rootProfileServiceClient.CreateNameserver(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrInvalidIP, err)
	})

	t.Run("UniqueViolation", func(t *testing.T) {
		t.Parallel()

		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)

		ip := "8.8.8.8"
		ipType := "A"

		seedNameserver(t, ip, ipType)

		req := &profilev1.CreateNameserverRequest{
			IpType: ipType,
			Ip:     ip,
		}

		res, err := rootProfileServiceClient.CreateNameserver(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrUniqueViolation("Nameserver"), err)
	})
}

func seedNameserver(t *testing.T, ip, ipType string) uuid.UUID {
	t.Helper()

	params := &repository.CreateNameserverParams{
		Ip:        ip,
		IpType:    ipType,
		CreatedBy: uuid.Nil(),
		UpdatedBy: uuid.Nil(),
	}

	id, err := cfg.Repository.CreateNameserver(t.Context(), params)
	require.NoError(t, err)
	return id
}

func getNameserver(t *testing.T, ip, ipType string) *repository.Nameserver {
	t.Helper()

	id := seedNameserver(t, ip, ipType)

	ns, err := cfg.Repository.Nameservers(t.Context())
	require.NoError(t, err)

	for _, n := range ns {
		if n.ID == id {
			return n
		}
	}
	return nil
}
