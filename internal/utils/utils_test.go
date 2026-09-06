//go:build unit

package utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"neupaneanish.com.np/profile/internal/utils"
)

func TestValidationURL(t *testing.T) {
	t.Parallel()

	t.Run("Malformed domain", func(t *testing.T) {
		t.Parallel()

		url, err := utils.ValidateURL("https://neupane:anish.com.np")
		require.Error(t, err)
		assert.Empty(t, url)
	})

	t.Run("Hacked Url", func(t *testing.T) {
		t.Parallel()

		url, err := utils.ValidateURL("hack@hack@hack.com/@hacked")
		require.Error(t, err)
		assert.Empty(t, url)
	})
}
