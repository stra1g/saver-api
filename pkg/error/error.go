package error

import (
	"errors"
	"fmt"
)

type ErrorType string

const (
	ErrorTypeNotFound      ErrorType = "NOT_FOUND"
	ErrorTypeValidation    ErrorType = "VALIDATION"
	ErrorTypeUnauthorized  ErrorType = "UNAUTHORIZED"
	ErrorTypeForbidden     ErrorType = "FORBIDDEN"
	ErrorTypeInternal      ErrorType = "INTERNAL"
	ErrorTypeDatabase      ErrorType = "DATABASE"
	ErrorTypeExternalAPI   ErrorType = "EXTERNAL_API"
	ErrorTypeUnprocessable ErrorType = "UNPROCESSABLE"
)

type AppError struct {
	errType     ErrorType
	originalErr error
	context     map[string]interface{}
}

func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s", e.errType, e.originalErr.Error())
}

func (e *AppError) Unwrap() error {
	return e.originalErr
}

func (e *AppError) Type() ErrorType {
	return e.errType
}

func (e *AppError) Context() map[string]interface{} {
	return e.context
}

func New(errType ErrorType, message string) *AppError {
	return &AppError{
		errType:     errType,
		originalErr: errors.New(message),
		context:     make(map[string]interface{}),
	}
}

func Wrap(errType ErrorType, err error) *AppError {
	if err == nil {
		return nil
	}
	
	if appErr, ok := err.(*AppError); ok {
		return &AppError{
			errType:     errType,
			originalErr: appErr.originalErr,
			context:     appErr.context,
		}
	}
	
	return &AppError{
		errType:     errType,
		originalErr: err,
		context:     make(map[string]interface{}),
	}
}

func WrapWithContext(errType ErrorType, err error, context map[string]interface{}) *AppError {
	wrappedErr := Wrap(errType, err)
	if wrappedErr == nil {
		return nil
	}
	
	for k, v := range context {
		wrappedErr.context[k] = v
	}
	
	return wrappedErr
}

func (e *AppError) AddContext(key string, value interface{}) *AppError {
	e.context[key] = value
	return e
}

func IsErrorType(err error, errType ErrorType) bool {
	if err == nil {
		return false
	}
	
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.errType == errType
	}
	
	return false
}
