package logger_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stra1g/saver-api/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func TestLoggerDebug(t *testing.T) {
	// Arrange
	var buf bytes.Buffer
	log := logger.Initialize(&buf, true)

	// Act
	log.Debug("test debug message", map[string]interface{}{
		"key1": "value1",
		"key2": 123,
	})

	// Assert
	var logEntry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logEntry)

	assert.NoError(t, err)
	assert.Equal(t, "test debug message", logEntry["message"])
	assert.Equal(t, "debug", logEntry["level"])
	assert.Equal(t, "value1", logEntry["key1"])
	assert.Equal(t, float64(123), logEntry["key2"])
}

func TestLoggerInfo(t *testing.T) {
	// Arrange
	var buf bytes.Buffer
	log := logger.Initialize(&buf, true)

	// Act
	log.Info("test info message", map[string]interface{}{
		"user_id": "abc-123",
	})

	// Assert
	var logEntry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logEntry)

	assert.NoError(t, err)
	assert.Equal(t, "test info message", logEntry["message"])
	assert.Equal(t, "info", logEntry["level"])
	assert.Equal(t, "abc-123", logEntry["user_id"])
}

func TestLoggerWarn(t *testing.T) {
	// Arrange
	var buf bytes.Buffer
	log := logger.Initialize(&buf, true)

	// Act
	log.Warn("test warning message", map[string]interface{}{
		"attempt": 3,
	})

	// Assert
	var logEntry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logEntry)

	assert.NoError(t, err)
	assert.Equal(t, "test warning message", logEntry["message"])
	assert.Equal(t, "warn", logEntry["level"])
	assert.Equal(t, float64(3), logEntry["attempt"])
}

func TestLoggerError(t *testing.T) {
	// Arrange
	var buf bytes.Buffer
	log := logger.Initialize(&buf, true)
	testErr := errors.New("test error")

	// Act
	log.Error(testErr, "failed operation", map[string]interface{}{
		"operation": "save_user",
	})

	// Assert
	var logEntry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logEntry)

	assert.NoError(t, err)
	assert.Equal(t, "failed operation", logEntry["message"])
	assert.Equal(t, "error", logEntry["level"])
	assert.Equal(t, "test error", logEntry["error"])
	assert.Equal(t, "save_user", logEntry["operation"])
}

func TestLoggerWithDebugDisabled(t *testing.T) {
	// Arrange
	var buf bytes.Buffer
	log := logger.Initialize(&buf, false) // Debug disabled

	// Act
	log.Debug("this should not appear", nil)

	// Assert
	assert.Empty(t, buf.String(), "Debug log should not be output when debug is disabled")

	// Verify other levels still work
	buf.Reset()
	log.Info("info still works", nil)

	var logEntry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logEntry)

	assert.NoError(t, err)
	assert.Equal(t, "info still works", logEntry["message"])
}

func TestLoggerWithNilFields(t *testing.T) {
	// Arrange
	var buf bytes.Buffer
	log := logger.Initialize(&buf, true)

	// Act
	log.Info("message with nil fields", nil)

	// Assert
	var logEntry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logEntry)

	assert.NoError(t, err)
	assert.Equal(t, "message with nil fields", logEntry["message"])
	assert.Equal(t, "info", logEntry["level"])
}

// Note: Testing Fatal is challenging because it calls os.Exit()
// We're skipping an actual Fatal test to avoid terminating the test process
