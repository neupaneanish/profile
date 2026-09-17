//go:build unit

package enum_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"neupaneanish.com.np/profile/internal/enum"
)

func TestMethod(t *testing.T) {
	t.Parallel()

	t.Run("Method", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name   string
			method enum.DBMethod
			want   bool
		}{
			{
				name:   "valid create",
				method: enum.DBMethodCreate,
				want:   true,
			},
			{
				name:   "valid update",
				method: enum.DBMethodUpdate,
				want:   true,
			},
			{
				name:   "valid delete",
				method: enum.DBMethodDelete,
				want:   true,
			},
			{
				name:   "invalid empty method",
				method: enum.DBMethod(""),
				want:   false,
			},
			{
				name:   "invalid random string",
				method: enum.DBMethod("admin"),
				want:   false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				got := tt.method.Valid()
				assert.Equal(t, tt.want, got)
			})
		}
	})

	t.Run("Table", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name  string
			table enum.DBTable
			want  bool
		}{
			{
				name:  "valid profile",
				table: enum.DBTableProfile,
				want:  true,
			},
			{
				name:  "valid about",
				table: enum.DBTableAbout,
				want:  true,
			},
			{
				name:  "valid education",
				table: enum.DBTableEducation,
				want:  true,
			},
			{
				name:  "valid experience",
				table: enum.DBTableExperience,
				want:  true,
			},
			{
				name:  "valid icon",
				table: enum.DBTableIcon,
				want:  true,
			},
			{
				name:  "valid social",
				table: enum.DBTableSocial,
				want:  true,
			},
			{
				name:  "valid nameserver",
				table: enum.DBTableNameserver,
				want:  true,
			},
			{
				name:  "valid domain",
				table: enum.DBTableDomain,
				want:  true,
			},
			{
				name:  "invalid empty method",
				table: enum.DBTable(""),
				want:  false,
			},
			{
				name:  "invalid random string",
				table: enum.DBTable("admin"),
				want:  false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				got := tt.table.Valid()
				assert.Equal(t, tt.want, got)
			})
		}
	})
}
