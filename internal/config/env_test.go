package config_test

import (
	"crypto/rand"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"neupaneanish.com.np/profile/internal/config"
)

func TestLoadEnv(t *testing.T) {
	cleanup := func() {
		_ = os.Unsetenv("DATABASE_HOST")
		_ = os.Unsetenv("DATABASE_NAME")
		_ = os.Unsetenv("DATABASE_USER")
		_ = os.Unsetenv("DATABASE_PASSWORD")
		_ = os.Unsetenv("DATABASE_PORT")
		_ = os.Unsetenv("DATABASE_SSL")
		_ = os.Unsetenv("VALKEY_URL")
		_ = os.Unsetenv("PORT")
		_ = os.Unsetenv("ENVIRONMENT")
		_ = os.Unsetenv("TELEMETRY_URL")
		_ = os.Unsetenv("SERVICE_NAME")
	}

	t.Run("Success with all variables", func(t *testing.T) {
		cleanup()

		t.Setenv("DATABASE_HOST", "127.0.0.1")
		t.Setenv("DATABASE_NAME", "postgres")
		t.Setenv("DATABASE_USER", "postgres")
		t.Setenv("DATABASE_PASSWORD", "postgres")
		t.Setenv("VALKEY_URL", "localhost:6379")
		t.Setenv("PORT", "50051")
		t.Setenv("SERVICE_NAME", "Test Service")
		t.Setenv("ENVIRONMENT", "production")
		t.Setenv("TELEMETRY_URL", "127.0.0.1:4317")

		env, envErr := config.LoadEnv()
		require.NoError(t, envErr)
		assert.NotNil(t, env)
		assert.Equal(t, "50051", env.Port)
	})

	t.Run("Default Environment", func(t *testing.T) {
		cleanup()
		t.Setenv("DATABASE_HOST", "127.0.0.1")
		t.Setenv("DATABASE_NAME", "postgres")
		t.Setenv("DATABASE_USER", "postgres")
		t.Setenv("DATABASE_PASSWORD", "postgres")
		t.Setenv("VALKEY_URL", "localhost:6379")
		t.Setenv("TELEMETRY_URL", "127.0.0.1:4317")

		env, envErr := config.LoadEnv()
		require.NoError(t, envErr)
		assert.NotNil(t, env)
		assert.Equal(t, "50051", env.Port)

		t.Run("Invalid port", func(t *testing.T) {
			t.Setenv("PORT", "79")
			pEnv, pEnvErr := config.LoadEnv()
			require.Error(t, pEnvErr)
			assert.Nil(t, pEnv)
		})

		t.Run("Invalid environment", func(t *testing.T) {
			t.Setenv("ENVIRONMENT", "staging")
			pEnv, pEnvErr := config.LoadEnv()
			require.Error(t, pEnvErr)
			assert.Nil(t, pEnv)
		})

		t.Run("Invalid Database PORT", func(t *testing.T) {
			t.Setenv("DATABASE_PORT", "79")
			pEnv, pEnvErr := config.LoadEnv()
			require.Error(t, pEnvErr)
			assert.Nil(t, pEnv)
		})

		t.Run("Database SSL False", func(t *testing.T) {
			t.Setenv("DATABASE_SSL", "False")
			pEnv, pEnvErr := config.LoadEnv()
			require.NoError(t, pEnvErr)
			assert.NotNil(t, pEnv)
		})

		t.Run("Database SSl Parse Failed", func(t *testing.T) {
			t.Setenv("DATABASE_SSL", rand.Text())
			pEnv, pEnvErr := config.LoadEnv()
			require.NoError(t, pEnvErr)
			assert.NotNil(t, pEnv)
		})
	})

	t.Run("Missing Required Environment", func(t *testing.T) {
		requiredVariables := []string{
			"DATABASE_HOST",
			"DATABASE_NAME",
			"DATABASE_USER",
			"DATABASE_PASSWORD",
			"VALKEY_URL",
			"TELEMETRY_URL",
		}

		for _, v := range requiredVariables {
			t.Run("Missing "+v, func(t *testing.T) {
				for _, all := range requiredVariables {
					t.Setenv(all, rand.Text())
				}

				_ = os.Unsetenv(v)

				env, err := config.LoadEnv()
				require.Error(t, err)
				assert.Nil(t, env)
				assert.Contains(t, err.Error(), v)
			})
		}
	})
}
