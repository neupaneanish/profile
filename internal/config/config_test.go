//go:build integration

package config_test

import (
	"os"
	"testing"

	"neupaneanish.com.np/profile/tests"
)

var (
	databaseURL string
	valkeyURL   string
)

func TestMain(m *testing.M) {
	dbURL, dbCleanup, dbErr := tests.Postgres()
	if dbErr != nil {
		panic(dbErr)
	}

	vkURL, valkeyCleanup, valkeyErr := tests.Valkey()
	if valkeyErr != nil {
		dbCleanup()
		panic(valkeyErr)
	}

	databaseURL = dbURL
	valkeyURL = vkURL

	code := m.Run()

	valkeyCleanup()
	dbCleanup()

	os.Exit(code)
}
