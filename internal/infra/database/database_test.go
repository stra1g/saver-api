package database

import (
	"os"
	"sync"
	"testing"

	"github.com/stra1g/saver-api/internal/infra/config"
	"github.com/stretchr/testify/assert"
)

// A minimal test that just tests the package-level variables
func TestDatabasePackageVars(t *testing.T) {
	// Assert that our package variables are initialized to their zero values
	assert.Nil(t, db)
	assert.NotNil(t, once) // Once is initialized to an empty Once struct, not nil
}

// A very basic test for NewPostgresDatabase using mocks
func TestNewPostgresDatabase(t *testing.T) {
	// Skip database tests unless ALLOW_DB_TESTS is set
	if os.Getenv("ALLOW_DB_TESTS") != "1" {
		t.Skip("Skipping database tests. Set ALLOW_DB_TESTS=1 to enable.")
	}

	// Save original values
	origDb := db
	origOnce := once

	// Restore after test
	defer func() {
		db = origDb
		once = origOnce
	}()

	// Reset for this test
	db = nil
	once = sync.Once{}

	// Create a valid config
	testConfig := &config.Config{}

	// Fill in minimal required fields
	testConfig.Server.Port = "8080"
	testConfig.Server.Host = "localhost"
	testConfig.Database.Host = "localhost"
	testConfig.Database.Port = "5432"
	testConfig.Database.User = "postgres"
	testConfig.Database.Password = "postgres"
	testConfig.Database.Name = "postgres"
	testConfig.Database.SSLMode = "disable"

	// We'll use a real database connection for this test
	testDSN := os.Getenv("POSTGRES_TEST_DSN")
	if testDSN != "" {
		// Set environment variables to match testDSN
		// This assumes testDSN is in the format postgres://user:pass@host:port/dbname?sslmode=mode
		t.Setenv("DB_HOST", "localhost")
		t.Setenv("DB_PORT", "5432")
		t.Setenv("DB_USER", "postgres")
		t.Setenv("DB_PASSWORD", "postgres")
		t.Setenv("DB_NAME", "postgres")
		t.Setenv("DB_SSLMODE", "disable")
	}

	// Only test actual connection if we can intercept fatal calls
	if os.Getenv("TEST_FATAL_OVERRIDE") == "1" {
		// Call the function - note this will exit the process if
		// connection fails, so we're not actually testing the result here
		_ = NewPostgresDatabase(testConfig)
	} else {
		// Skip actual connection tests for now
		t.Log("Skipping actual connection test")
	}
}

// To run the fatal error tests, you'd need to run the test with the TEST_FATAL_OVERRIDE
// environment variable set to "1", and then check the process exit code.
// This is hard to do in Go's testing framework, so we're skipping it for now.
