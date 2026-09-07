//go:build integration

package service_test

import (
	"crypto/rand"
	"strings"
	"testing"
	"uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"neupaneanish.com.np/profile/internal/enum"
	"neupaneanish.com.np/profile/internal/errs"
	profilev1 "neupaneanish.com.np/profile/internal/protobuf/root/profile/v1"
)

func TestUpdateNameServer(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		name := strings.ToLower(rand.Text()[:8])
		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)

		nameserver := getNameServer(t, name)

		req := &profilev1.UpdateNameServerRequest{
			Id: nameserver.ID.String(),
			Nameserver: &profilev1.CreateUpdateNameServer{
				Domain: "https://" + nameserver.Domain,
				Cname:  nameserver.Cname,
			},
			Active:    false,
			UpdatedAt: timestamppb.New(nameserver.UpdatedAt),
		}

		res, err := rootProfileServiceClient.UpdateNameServer(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, nameserver.ID.String(), res.GetId())
	})

	t.Run("No Update", func(t *testing.T) {
		t.Parallel()
		name := strings.ToLower(rand.Text()[:8])
		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)

		nameserver := getNameServer(t, name)

		req := &profilev1.UpdateNameServerRequest{
			Id: nameserver.ID.String(),
			Nameserver: &profilev1.CreateUpdateNameServer{
				Domain: "https://" + nameserver.Domain,
				Cname:  nameserver.Cname,
			},
			Active:    nameserver.Active,
			UpdatedAt: timestamppb.New(nameserver.UpdatedAt),
		}

		res, err := rootProfileServiceClient.UpdateNameServer(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrConflict, err)
	})

	t.Run("Unique Violation", func(t *testing.T) {
		t.Parallel()
		name := strings.ToLower(rand.Text()[:8])
		ctx := contextWithValue(t, uuid.NewV7(), enum.UserRoleRoot)

		nameserver := getNameServer(t, name)
		nameserver2 := getNameServer(t, name+"1")

		req := &profilev1.UpdateNameServerRequest{
			Id: nameserver.ID.String(),
			Nameserver: &profilev1.CreateUpdateNameServer{
				Domain: "https://" + nameserver2.Domain,
				Cname:  nameserver2.Cname,
			},
			Active:    nameserver.Active,
			UpdatedAt: timestamppb.New(nameserver.UpdatedAt),
		}

		res, err := rootProfileServiceClient.UpdateNameServer(ctx, req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, errs.ErrUniqueViolation("Nameserver"), err)
	})
}
