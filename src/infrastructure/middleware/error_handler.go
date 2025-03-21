package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vnlab/makeshop-payment/src/infrastructure/logger"
)

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Status  int               `json:"status"`
	Message string            `json:"message"`
	Code    string            `json:"code,omitempty"`
	Details map[string]string `json:"details,omitempty"`
}

// CustomError represents a structured application error
type CustomError struct {
	Status  int               // HTTP status code
	Message string            // User-friendly error message
	Code    string            // Error code for frontend handling
	Details map[string]string // Additional error details
	Err     error             // Original error (not exposed to client)
}

// Error implements the error interface
func (e *CustomError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// ErrorHandlerMiddleware creates middleware for global error handling and logging
func ErrorHandlerMiddleware(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get request-specific logger from context
		requestLogger, exists := c.Get("logger")
		if !exists {
			requestLogger = log
		}
		logger := requestLogger.(logger.Logger)

		// Recover from any panics
		defer func() {
			if err := recover(); err != nil {
				// Log the stack trace
				stack := debug.Stack()
				logger.Error("Panic recovered", map[string]interface{}{
					"error": fmt.Sprintf("%v", err),
					"stack": string(stack),
				})

				// Return a 500 error to the client
				c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{
					Status:  http.StatusInternalServerError,
					Message: "Internal server error",
					Code:    "INTERNAL_ERROR",
				})
			}
		}()

		c.Next()

		// If there are errors collected in the Gin context, log and return them
		if len(c.Errors) > 0 {
			var status int
			var customErr *CustomError
			lastErr := c.Errors.Last().Err

			// Try to cast to CustomError
			if ce, ok := lastErr.(*CustomError); ok {
				customErr = ce
				status = ce.Status
			} else {
				// Default error handling
				status = http.StatusInternalServerError
				customErr = &CustomError{
					Status:  status,
					Message: "An unexpected error occurred",
					Code:    "INTERNAL_ERROR",
					Err:     lastErr,
				}
			}

			// Create error message for logging
			var errMsgs []string
			for _, e := range c.Errors {
				errMsgs = append(errMsgs, e.Error())
			}

			// Log the error with appropriate level based on status code
			logFields := map[string]interface{}{
				"status_code": status,
				"error_code":  customErr.Code,
				"request_url": c.Request.URL.String(),
				"method":      c.Request.Method,
				"user_id":     getUserID(c),
				"details":     customErr.Details,
			}

			if customErr.Err != nil {
				logFields["original_error"] = customErr.Err.Error()
			}

			if status >= 500 {
				logger.Error(strings.Join(errMsgs, "; "), logFields)
			} else {
				logger.Warn(strings.Join(errMsgs, "; "), logFields)
			}

			// Return the error response to the client
			c.JSON(status, ErrorResponse{
				Status:  status,
				Message: customErr.Message,
				Code:    customErr.Code,
				Details: customErr.Details,
			})
		}
	}
}

// NewCustomError creates a new CustomError
func NewCustomError(status int, code string, message string, err error) *CustomError {
	return &CustomError{
		Status:  status,
		Message: message,
		Code:    code,
		Err:     err,
	}
}

// AddDetail adds a detail to a CustomError
func (e *CustomError) AddDetail(key, value string) *CustomError {
	if e.Details == nil {
		e.Details = make(map[string]string)
	}
	e.Details[key] = value
	return e
}

// getUserID extracts user ID from Gin context (set by auth middleware)
func getUserID(c *gin.Context) interface{} {
	userID, exists := c.Get("userId")
	if exists {
		return userID
	}
	return nil
}
