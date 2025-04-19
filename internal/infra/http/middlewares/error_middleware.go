package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apperror "github.com/stra1g/saver-api/pkg/error"
	"github.com/stra1g/saver-api/pkg/logger"
)

type ErrorResponse struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

func ErrorHandler(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			handleError(c, err, log)
		}
	}
}

func handleError(c *gin.Context, err error, log logger.Logger) {
	if appErr, ok := err.(*apperror.AppError); ok {
		switch appErr.Type() {
		case apperror.ErrorTypeNotFound:
			c.JSON(http.StatusNotFound, ErrorResponse{
				Code:    string(appErr.Type()),
				Message: appErr.Error(),
				Details: appErr.Context(),
			})
		case apperror.ErrorTypeValidation:
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Code:    string(appErr.Type()),
				Message: appErr.Error(),
				Details: appErr.Context(),
			})
		case apperror.ErrorTypeUnauthorized:
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Code:    string(appErr.Type()),
				Message: appErr.Error(),
				Details: appErr.Context(),
			})
		case apperror.ErrorTypeForbidden:
			c.JSON(http.StatusForbidden, ErrorResponse{
				Code:    string(appErr.Type()),
				Message: appErr.Error(),
				Details: appErr.Context(),
			})
		case apperror.ErrorTypeUnprocessable:
			c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
				Code:    string(appErr.Type()),
				Message: appErr.Error(),
				Details: appErr.Context(),
			})
		case apperror.ErrorTypeDatabase, apperror.ErrorTypeExternalAPI, apperror.ErrorTypeInternal:
			log.Error(appErr, "Internal server error", appErr.Context())
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Code:    "INTERNAL_SERVER_ERROR",
				Message: "An unexpected error occurred",
			})
		default:
			log.Error(appErr, "Unhandled error type", map[string]interface{}{
				"type": appErr.Type(),
			})
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Code:    "INTERNAL_SERVER_ERROR",
				Message: "An unexpected error occurred",
			})
		}
		return
	}

	log.Error(err, "Unexpected error", nil)
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Code:    "INTERNAL_SERVER_ERROR",
		Message: "An unexpected error occurred",
	})
}

func NewErrorHandler(log logger.Logger) gin.HandlerFunc {
	return ErrorHandler(log)
}
