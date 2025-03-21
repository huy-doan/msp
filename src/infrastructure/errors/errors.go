package errors

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/vnlab/makeshop-payment/src/infrastructure/middleware"
)

// Standard error codes
const (
	CodeInternalError      = "INTERNAL_ERROR"
	CodeNotFound           = "NOT_FOUND"
	CodeBadRequest         = "BAD_REQUEST"
	CodeUnauthorized       = "UNAUTHORIZED"
	CodeForbidden          = "FORBIDDEN"
	CodeValidationError    = "VALIDATION_ERROR"
	CodeDuplicateEntry     = "DUPLICATE_ENTRY"
	CodeDatabaseError      = "DATABASE_ERROR"
	CodeResourceNotFound   = "RESOURCE_NOT_FOUND"
	CodeInvalidCredentials = "INVALID_CREDENTIALS"
)

// New creates a new error with message
func New(message string) error {
	return errors.New(message)
}

// Wrap wraps an error with additional message
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}

// Internal creates a new internal server error
func Internal(message string, err error) *middleware.CustomError {
	return &middleware.CustomError{
		Status:  http.StatusInternalServerError,
		Message: message,
		Code:    CodeInternalError,
		Err:     err,
	}
}

// NotFound creates a new not found error
func NotFound(message string, resourceType string) *middleware.CustomError {
	ce := &middleware.CustomError{
		Status:  http.StatusNotFound,
		Message: message,
		Code:    CodeResourceNotFound,
	}
	if resourceType != "" {
		ce.AddDetail("resource_type", resourceType)
	}
	return ce
}

// BadRequest creates a new bad request error
func BadRequest(message string, err error) *middleware.CustomError {
	return &middleware.CustomError{
		Status:  http.StatusBadRequest,
		Message: message,
		Code:    CodeBadRequest,
		Err:     err,
	}
}

// Unauthorized creates a new unauthorized error
func Unauthorized(message string) *middleware.CustomError {
	if message == "" {
		message = "Authentication required"
	}
	return &middleware.CustomError{
		Status:  http.StatusUnauthorized,
		Message: message,
		Code:    CodeUnauthorized,
	}
}

// Forbidden creates a new forbidden error
func Forbidden(message string) *middleware.CustomError {
	if message == "" {
		message = "You don't have permission to perform this action"
	}
	return &middleware.CustomError{
		Status:  http.StatusForbidden,
		Message: message,
		Code:    CodeForbidden,
	}
}

// Validation creates a new validation error
func Validation(message string, details map[string]string) *middleware.CustomError {
	ce := &middleware.CustomError{
		Status:  http.StatusBadRequest,
		Message: message,
		Code:    CodeValidationError,
		Details: details,
	}
	return ce
}

// Database creates a new database error (typically internal server error)
func Database(message string, err error) *middleware.CustomError {
	return &middleware.CustomError{
		Status:  http.StatusInternalServerError,
		Message: message,
		Code:    CodeDatabaseError,
		Err:     err,
	}
}

// DuplicateEntry creates a new duplicate entry error
func DuplicateEntry(message string, field string) *middleware.CustomError {
	ce := &middleware.CustomError{
		Status:  http.StatusConflict,
		Message: message,
		Code:    CodeDuplicateEntry,
	}
	if field != "" {
		ce.AddDetail("field", field)
	}
	return ce
}

// InvalidCredentials creates a new invalid credentials error
func InvalidCredentials(message string) *middleware.CustomError {
	if message == "" {
		message = "Invalid email or password"
	}
	return &middleware.CustomError{
		Status:  http.StatusUnauthorized,
		Message: message,
		Code:    CodeInvalidCredentials,
	}
}
