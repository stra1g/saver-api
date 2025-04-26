package config_test

import (
	"os"
	"testing"

	"github.com/stra1g/saver-api/internal/infra/config"
	"github.com/stretchr/testify/assert"
)

func setupEnvVars(t *testing.T) func() {
	// Save original env vars to restore them later
	originalVars := map[string]string{
		"PORT":        os.Getenv("PORT"),
		"HOST":        os.Getenv("HOST"),
		"DB_HOST":     os.Getenv("DB_HOST"),
		"DB_PORT":     os.Getenv("DB_PORT"),
		"DB_USER":     os.Getenv("DB_USER"),
		"DB_PASSWORD": os.Getenv("DB_PASSWORD"),
		"DB_NAME":     os.Getenv("DB_NAME"),
		"DB_SSLMODE":  os.Getenv("DB_SSLMODE"),
	}

	// Set test environment variables
	t.Setenv("PORT", "3000")
	t.Setenv("HOST", "127.0.0.1")
	t.Setenv("DB_HOST", "test-db-host")
	t.Setenv("DB_PORT", "5433")
	t.Setenv("DB_USER", "testuser")
	t.Setenv("DB_PASSWORD", "testpass")
	t.Setenv("DB_NAME", "testdb")
	t.Setenv("DB_SSLMODE", "disable")

	// Return cleanup function
	return func() {
		for key, value := range originalVars {
			if value == "" {
				os.Unsetenv(key)
			} else {
				os.Setenv(key, value)
			}
		}
	}
}

func TestNewConfig_WithEnvVars(t *testing.T) {
	// Setup and defer cleanup of environment
	cleanup := setupEnvVars(t)
	defer cleanup()

	// Act
	cfg, err := config.NewConfig()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	// Server config assertions
	assert.Equal(t, "3000", cfg.Server.Port)
	assert.Equal(t, "127.0.0.1", cfg.Server.Host)

	// Database config assertions
	assert.Equal(t, "test-db-host", cfg.Database.Host)
	assert.Equal(t, "5433", cfg.Database.Port)
	assert.Equal(t, "testuser", cfg.Database.User)
	assert.Equal(t, "testpass", cfg.Database.Password)
	assert.Equal(t, "testdb", cfg.Database.Name)
	assert.Equal(t, "disable", cfg.Database.SSLMode)
}

func TestNewConfig_WithDefaults(t *testing.T) {
	// Clear all relevant environment variables
	vars := []string{"PORT", "HOST", "DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_SSLMODE"}
	originalVars := make(map[string]string, len(vars))

	for _, key := range vars {
		originalVars[key] = os.Getenv(key)
		os.Unsetenv(key)
	}
	defer func() {
		for key, value := range originalVars {
			if value != "" {
				os.Setenv(key, value)
			}
		}
	}()

	// We need to set these to non-empty or validation will fail
	t.Setenv("DB_USER", "defaultuser")
	t.Setenv("DB_PASSWORD", "defaultpass")
	t.Setenv("DB_NAME", "defaultdb")

	// Act
	cfg, err := config.NewConfig()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	// Check default values
	assert.Equal(t, "8080", cfg.Server.Port)
	assert.Equal(t, "0.0.0.0", cfg.Server.Host)
	assert.Equal(t, "localhost", cfg.Database.Host)
	assert.Equal(t, "5432", cfg.Database.Port)
	assert.Equal(t, "defaultuser", cfg.Database.User)
	assert.Equal(t, "defaultpass", cfg.Database.Password)
	assert.Equal(t, "defaultdb", cfg.Database.Name)
	assert.Equal(t, "disable", cfg.Database.SSLMode)
}

func TestNewConfig_ValidationError_MissingRequired(t *testing.T) {
	// Clear specific environment variables to trigger validation errors
	vars := []string{"DB_USER", "DB_PASSWORD", "DB_NAME"}
	originalVars := make(map[string]string, len(vars))

	for _, key := range vars {
		originalVars[key] = os.Getenv(key)
		os.Unsetenv(key)
	}
	defer func() {
		for key, value := range originalVars {
			if value != "" {
				os.Setenv(key, value)
			}
		}
	}()

	// Act
	cfg, err := config.NewConfig()

	// Assert
	assert.Error(t, err)
	assert.Nil(t, cfg)
	assert.Contains(t, err.Error(), "Configuration validation failed")
	assert.Contains(t, err.Error(), "User is required")
	assert.Contains(t, err.Error(), "Password is required")
	assert.Contains(t, err.Error(), "Name is required")
}

func TestNewConfig_ValidationError_NonNumeric(t *testing.T) {
	// Setup environment with invalid numeric values
	cleanup := setupEnvVars(t)
	defer cleanup()

	t.Setenv("PORT", "not-numeric")
	t.Setenv("DB_PORT", "also-not-numeric")

	// Act
	cfg, err := config.NewConfig()

	// Assert
	assert.Error(t, err)
	assert.Nil(t, cfg)
	assert.Contains(t, err.Error(), "Configuration validation failed")
	assert.Contains(t, err.Error(), "Port must be a numeric value")
}

func TestNewConfig_ValidationError_InvalidSSLMode(t *testing.T) {
	// Setup environment with invalid SSL mode
	cleanup := setupEnvVars(t)
	defer cleanup()

	t.Setenv("DB_SSLMODE", "invalid-mode")

	// Act
	cfg, err := config.NewConfig()

	// Assert
	assert.Error(t, err)
	assert.Nil(t, cfg)
	assert.Contains(t, err.Error(), "Configuration validation failed")
	assert.Contains(t, err.Error(), "SSLMode failed validation: oneof")
}

func TestGetDatabaseDSN(t *testing.T) {
	// Setup test config
	cfg := &config.Config{}
	cfg.Database.Host = "testhost"
	cfg.Database.Port = "5432"
	cfg.Database.User = "testuser"
	cfg.Database.Password = "testpass"
	cfg.Database.Name = "testdb"
	cfg.Database.SSLMode = "disable"

	// Act
	dsn := cfg.GetDatabaseDSN()

	// Assert
	expectedDSN := "postgres://testuser:testpass@testhost:5432/testdb?sslmode=disable"
	assert.Equal(t, expectedDSN, dsn)
}

func TestGetEnvWithDefault(t *testing.T) {
	// Test with environment variable set
	originalValue := os.Getenv("TEST_KEY")
	defer func() {
		if originalValue == "" {
			os.Unsetenv("TEST_KEY")
		} else {
			os.Setenv("TEST_KEY", originalValue)
		}
	}()

	// Test 1: Env var is set
	os.Setenv("TEST_KEY", "test-value")
	result := config.GetEnvWithDefault("TEST_KEY", "default-value")
	assert.Equal(t, "test-value", result)

	// Test 2: Env var is not set
	os.Unsetenv("TEST_KEY")
	result = config.GetEnvWithDefault("TEST_KEY", "default-value")
	assert.Equal(t, "default-value", result)

	// Test 3: Env var is empty
	os.Setenv("TEST_KEY", "")
	result = config.GetEnvWithDefault("TEST_KEY", "default-value")
	assert.Equal(t, "default-value", result)
}

func TestValidateConfig(t *testing.T) {
	// Test 1: Valid configuration
	validCfg := &config.Config{}
	validCfg.Server.Port = "8080"
	validCfg.Server.Host = "localhost"
	validCfg.Database.Host = "testhost"
	validCfg.Database.Port = "5432"
	validCfg.Database.User = "testuser"
	validCfg.Database.Password = "testpass"
	validCfg.Database.Name = "testdb"
	validCfg.Database.SSLMode = "disable"

	err := config.ValidateConfig(validCfg)
	assert.NoError(t, err)

	// Test 2: Missing required field
	invalidCfg := &config.Config{}
	invalidCfg.Server.Port = "8080"
	invalidCfg.Server.Host = "localhost"
	invalidCfg.Database.Host = "testhost"
	invalidCfg.Database.Port = "5432"
	// Missing User field
	invalidCfg.Database.Password = "testpass"
	invalidCfg.Database.Name = "testdb"
	invalidCfg.Database.SSLMode = "disable"

	err = config.ValidateConfig(invalidCfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "User is required")

	// Test 3: Invalid numeric field
	invalidCfg = &config.Config{}
	invalidCfg.Server.Port = "not-a-number" // Invalid numeric value
	invalidCfg.Server.Host = "localhost"
	invalidCfg.Database.Host = "testhost"
	invalidCfg.Database.Port = "5432"
	invalidCfg.Database.User = "testuser"
	invalidCfg.Database.Password = "testpass"
	invalidCfg.Database.Name = "testdb"
	invalidCfg.Database.SSLMode = "disable"

	err = config.ValidateConfig(invalidCfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Port must be a numeric value")

	// Test 4: Invalid enum value
	invalidCfg = &config.Config{}
	invalidCfg.Server.Port = "8080"
	invalidCfg.Server.Host = "localhost"
	invalidCfg.Database.Host = "testhost"
	invalidCfg.Database.Port = "5432"
	invalidCfg.Database.User = "testuser"
	invalidCfg.Database.Password = "testpass"
	invalidCfg.Database.Name = "testdb"
	invalidCfg.Database.SSLMode = "invalid-mode" // Invalid enum value

	err = config.ValidateConfig(invalidCfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "SSLMode failed validation: oneof")
}
