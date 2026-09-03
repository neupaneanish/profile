//go:build unit

package enum_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"neupaneanish.com.np/profile/internal/enum"
)

func TestLocationType(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		locationType enum.LocationType
		want         bool
	}{
		{
			name:         "valid onsite",
			locationType: enum.LocationTypeOnSite,
			want:         true,
		},
		{
			name:         "valid remote",
			locationType: enum.LocationTypeRemote,
			want:         true,
		},
		{
			name:         "valid hybrid",
			locationType: enum.LocationTypeHybrid,
			want:         true,
		},
		{
			name:         "invalid empty role",
			locationType: enum.LocationType(""),
			want:         false,
		},
		{
			name:         "invalid random string",
			locationType: enum.LocationType("admin"),
			want:         false,
		},
		{
			name:         "invalid case sensitivity",
			locationType: enum.LocationType("ONSITE"),
			want:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.locationType.Valid()
			assert.Equal(t, tt.want, got)
		})
	}
}
