package error_test

import (
	"errors"
	"fmt"
	"testing"

	apperror "github.com/stra1g/saver-api/pkg/error"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	// Arrange
	errType := apperror.ErrorTypeNotFound
	message := "resource not found"

	// Act
	err := apperror.New(errType, message)

	// Assert
	assert.NotNil(t, err)
	assert.Equal(t, errType, err.Type())
	assert.Equal(t, fmt.Sprintf("[%s] %s", errType, message), err.Error())
	assert.Empty(t, err.Context())
}

func TestWrap(t *testing.T) {
	// Arrange
	errType := apperror.ErrorTypeValidation
	originalErr := errors.New("original error")

	// Act
	err := apperror.Wrap(errType, originalErr)

	// Assert
	assert.NotNil(t, err)
	assert.Equal(t, errType, err.Type())
	assert.Equal(t, fmt.Sprintf("[%s] %s", errType, originalErr.Error()), err.Error())
	assert.Empty(t, err.Context())

	// Test Unwrap functionality
	unwrappedErr := errors.Unwrap(err)
	assert.Equal(t, originalErr, unwrappedErr)
}

func TestWrapWithContext(t *testing.T) {
	// Arrange
	errType := apperror.ErrorTypeDatabase
	originalErr := errors.New("database error")
	context := map[string]interface{}{
		"table":  "users",
		"action": "insert",
	}

	// Act
	err := apperror.WrapWithContext(errType, originalErr, context)

	// Assert
	assert.NotNil(t, err)
	assert.Equal(t, errType, err.Type())
	assert.Equal(t, fmt.Sprintf("[%s] %s", errType, originalErr.Error()), err.Error())
	assert.Equal(t, context, err.Context())
}

func TestAddContext(t *testing.T) {
	// Arrange
	errType := apperror.ErrorTypeInternal
	err := apperror.New(errType, "internal error")

	// Act
	err.AddContext("function", "TestAddContext")
	err.AddContext("time", "now")

	// Assert
	expectedContext := map[string]interface{}{
		"function": "TestAddContext",
		"time":     "now",
	}
	assert.Equal(t, expectedContext, err.Context())
}

func TestIsErrorType(t *testing.T) {
	tests := []struct {
		name          string
		err           error
		errType       apperror.ErrorType
		expectedMatch bool
	}{
		{
			name:          "same error type",
			err:           apperror.New(apperror.ErrorTypeNotFound, "not found"),
			errType:       apperror.ErrorTypeNotFound,
			expectedMatch: true,
		},
		{
			name:          "different error type",
			err:           apperror.New(apperror.ErrorTypeValidation, "validation error"),
			errType:       apperror.ErrorTypeNotFound,
			expectedMatch: false,
		},
		{
			name:          "wrapped error with same type",
			err:           apperror.Wrap(apperror.ErrorTypeForbidden, errors.New("access denied")),
			errType:       apperror.ErrorTypeForbidden,
			expectedMatch: true,
		},
		{
			name:          "non-app error",
			err:           errors.New("standard error"),
			errType:       apperror.ErrorTypeInternal,
			expectedMatch: false,
		},
		{
			name:          "nil error",
			err:           nil,
			errType:       apperror.ErrorTypeInternal,
			expectedMatch: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := apperror.IsErrorType(tc.err, tc.errType)
			assert.Equal(t, tc.expectedMatch, result)
		})
	}
}

func TestWrapWithEmptyOriginalError(t *testing.T) {
	// Arrange
	errType := apperror.ErrorTypeDatabase
	originalErr := errors.New("")

	// Act
	err := apperror.Wrap(errType, originalErr)

	// Assert
	assert.NotNil(t, err)
	assert.Equal(t, errType, err.Type())
	assert.Equal(t, fmt.Sprintf("[%s] ", errType), err.Error())
}

func TestWrapWithNilError(t *testing.T) {
	// Arrange
	errType := apperror.ErrorTypeDatabase

	// Act
	err := apperror.Wrap(errType, nil)

	// Assert
	assert.Nil(t, err)
}

func TestWrapWithContextNilError(t *testing.T) {
	// Arrange
	errType := apperror.ErrorTypeDatabase
	context := map[string]interface{}{
		"table": "users",
	}

	// Act
	err := apperror.WrapWithContext(errType, nil, context)

	// Assert
	assert.Nil(t, err)
}
