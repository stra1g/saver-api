package middlewares_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stra1g/saver-api/internal/infra/http/middlewares"
	apperror "github.com/stra1g/saver-api/pkg/error"
	mocks "github.com/stra1g/saver-api/pkg/testutils/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupRouter(mockLogger *mocks.MockLogger) (*gin.Engine, *httptest.ResponseRecorder) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a mock recorder
	recorder := httptest.NewRecorder()

	// Create a new Gin router
	router := gin.New()

	// Use the error handler middleware
	router.Use(middlewares.NewErrorHandler(mockLogger))

	return router, recorder
}

func TestErrorHandler_NoErrors(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	router, recorder := setupRouter(mockLogger)

	// Set up a test route that doesn't generate errors
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Create a test request
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)

	// Act
	router.ServeHTTP(recorder, req)

	// Assert
	assert.Equal(t, http.StatusOK, recorder.Code)
	var response map[string]interface{}
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "ok", response["status"])

	// Verify logger wasn't called
	mockLogger.AssertNotCalled(t, "Error")
}

func TestErrorHandler_NotFoundError(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	router, recorder := setupRouter(mockLogger)

	// Set up a test route that generates a not found error
	router.GET("/test", func(c *gin.Context) {
		err := apperror.New(apperror.ErrorTypeNotFound, "Resource not found")
		c.Error(err)
	})

	// Create a test request
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)

	// Act
	router.ServeHTTP(recorder, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, recorder.Code)

	var response middlewares.ErrorResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, string(apperror.ErrorTypeNotFound), response.Code)
	assert.Contains(t, response.Message, "Resource not found")

	// Verify logger wasn't called for not found errors
	mockLogger.AssertNotCalled(t, "Error")
}

func TestErrorHandler_ValidationError(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	router, recorder := setupRouter(mockLogger)

	// Set up a test route that generates a validation error
	router.GET("/test", func(c *gin.Context) {
		err := apperror.New(apperror.ErrorTypeValidation, "Validation failed")
		c.Error(err)
	})

	// Create a test request
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)

	// Act
	router.ServeHTTP(recorder, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	var response middlewares.ErrorResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, string(apperror.ErrorTypeValidation), response.Code)
	assert.Contains(t, response.Message, "Validation failed")
}

func TestErrorHandler_UnauthorizedError(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	router, recorder := setupRouter(mockLogger)

	// Set up a test route that generates an unauthorized error
	router.GET("/test", func(c *gin.Context) {
		err := apperror.New(apperror.ErrorTypeUnauthorized, "Unauthorized access")
		c.Error(err)
	})

	// Create a test request
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)

	// Act
	router.ServeHTTP(recorder, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, recorder.Code)

	var response middlewares.ErrorResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, string(apperror.ErrorTypeUnauthorized), response.Code)
	assert.Contains(t, response.Message, "Unauthorized access")
}

func TestErrorHandler_ForbiddenError(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	router, recorder := setupRouter(mockLogger)

	// Set up a test route that generates a forbidden error
	router.GET("/test", func(c *gin.Context) {
		err := apperror.New(apperror.ErrorTypeForbidden, "Access forbidden")
		c.Error(err)
	})

	// Create a test request
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)

	// Act
	router.ServeHTTP(recorder, req)

	// Assert
	assert.Equal(t, http.StatusForbidden, recorder.Code)

	var response middlewares.ErrorResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, string(apperror.ErrorTypeForbidden), response.Code)
	assert.Contains(t, response.Message, "Access forbidden")
}

func TestErrorHandler_UnprocessableError(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	router, recorder := setupRouter(mockLogger)

	// Set up a test route that generates an unprocessable entity error
	router.GET("/test", func(c *gin.Context) {
		err := apperror.New(apperror.ErrorTypeUnprocessable, "Unprocessable entity")
		c.Error(err)
	})

	// Create a test request
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)

	// Act
	router.ServeHTTP(recorder, req)

	// Assert
	assert.Equal(t, http.StatusUnprocessableEntity, recorder.Code)

	var response middlewares.ErrorResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, string(apperror.ErrorTypeUnprocessable), response.Code)
	assert.Contains(t, response.Message, "Unprocessable entity")
}

func TestErrorHandler_DatabaseError(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockLogger.On("Error", mock.Anything, "Internal server error", mock.Anything).Return()

	router, recorder := setupRouter(mockLogger)

	// Set up a test route that generates a database error
	router.GET("/test", func(c *gin.Context) {
		err := apperror.New(apperror.ErrorTypeDatabase, "Database connection failed")
		c.Error(err)
	})

	// Create a test request
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)

	// Act
	router.ServeHTTP(recorder, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)

	var response middlewares.ErrorResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "INTERNAL_SERVER_ERROR", response.Code)
	assert.Equal(t, "An unexpected error occurred", response.Message)

	// Verify logger was called
	mockLogger.AssertCalled(t, "Error", mock.Anything, "Internal server error", mock.Anything)
}

func TestErrorHandler_ExternalAPIError(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockLogger.On("Error", mock.Anything, "Internal server error", mock.Anything).Return()

	router, recorder := setupRouter(mockLogger)

	// Set up a test route that generates an external API error
	router.GET("/test", func(c *gin.Context) {
		err := apperror.New(apperror.ErrorTypeExternalAPI, "External API failed")
		c.Error(err)
	})

	// Create a test request
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)

	// Act
	router.ServeHTTP(recorder, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)

	var response middlewares.ErrorResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "INTERNAL_SERVER_ERROR", response.Code)
	assert.Equal(t, "An unexpected error occurred", response.Message)

	// Verify logger was called
	mockLogger.AssertCalled(t, "Error", mock.Anything, "Internal server error", mock.Anything)
}

func TestErrorHandler_InternalError(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockLogger.On("Error", mock.Anything, "Internal server error", mock.Anything).Return()

	router, recorder := setupRouter(mockLogger)

	// Set up a test route that generates an internal error
	router.GET("/test", func(c *gin.Context) {
		err := apperror.New(apperror.ErrorTypeInternal, "Internal error occurred")
		c.Error(err)
	})

	// Create a test request
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)

	// Act
	router.ServeHTTP(recorder, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)

	var response middlewares.ErrorResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "INTERNAL_SERVER_ERROR", response.Code)
	assert.Equal(t, "An unexpected error occurred", response.Message)

	// Verify logger was called
	mockLogger.AssertCalled(t, "Error", mock.Anything, "Internal server error", mock.Anything)
}

func TestErrorHandler_UnknownAppErrorType(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	mockLogger.On("Error", mock.Anything, "Unhandled error type", mock.Anything).Return()

	router, recorder := setupRouter(mockLogger)

	// Create a custom error type for this test
	const customErrorType apperror.ErrorType = "CUSTOM_ERROR_TYPE"

	// Set up a test route that generates an error with an unknown type
	router.GET("/test", func(c *gin.Context) {
		err := apperror.New(customErrorType, "Custom error occurred")
		c.Error(err)
	})

	// Create a test request
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)

	// Act
	router.ServeHTTP(recorder, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)

	var response middlewares.ErrorResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "INTERNAL_SERVER_ERROR", response.Code)
	assert.Equal(t, "An unexpected error occurred", response.Message)

	// Verify logger was called with the right parameters
	mockLogger.AssertCalled(t, "Error", mock.Anything, "Unhandled error type", mock.Anything)
}

func TestErrorHandler_NonAppError(t *testing.T) {
	// Create a special wrapper to capture the actual call
	mockLogger := mocks.NewMockLogger()
	var capturedErr error
	var capturedMsg string

	// Override the Error method to capture parameters
	mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		capturedErr = args.Get(0).(error)
		capturedMsg = args.Get(1).(string)
	}).Return()

	router, recorder := setupRouter(mockLogger)

	// Set up a test route that generates a standard error
	router.GET("/test", func(c *gin.Context) {
		err := errors.New("Standard error occurred")
		c.Error(err)
	})

	// Create a test request
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)

	// Act
	router.ServeHTTP(recorder, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)

	var response middlewares.ErrorResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "INTERNAL_SERVER_ERROR", response.Code)
	assert.Equal(t, "An unexpected error occurred", response.Message)

	// Verify the call happened and check captured values
	mockLogger.AssertCalled(t, "Error", mock.Anything, mock.Anything, mock.Anything)
	assert.Equal(t, "Standard error occurred", capturedErr.Error())
	assert.Equal(t, "Unexpected error", capturedMsg)
	// We don't strictly check the map - the important thing is that the Error method was called
}

func TestErrorHandler_MultipleErrors(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()
	router, recorder := setupRouter(mockLogger)

	// Set up a test route that generates multiple errors
	router.GET("/test", func(c *gin.Context) {
		// Add first error
		err1 := apperror.New(apperror.ErrorTypeValidation, "First error")
		c.Error(err1)

		// Add second error - this one should be handled as it's the last one
		err2 := apperror.New(apperror.ErrorTypeNotFound, "Last error")
		c.Error(err2)
	})

	// Create a test request
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)

	// Act
	router.ServeHTTP(recorder, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, recorder.Code)

	var response middlewares.ErrorResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, string(apperror.ErrorTypeNotFound), response.Code)
	assert.Contains(t, response.Message, "Last error")
}

func TestNewErrorHandler(t *testing.T) {
	// Arrange
	mockLogger := mocks.NewMockLogger()

	// Act
	handler := middlewares.NewErrorHandler(mockLogger)

	// Assert
	assert.NotNil(t, handler)
}
