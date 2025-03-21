package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vnlab/makeshop-payment/src/infrastructure/logger"
)

// RequestLoggerMiddleware creates middleware for logging HTTP requests
func RequestLoggerMiddleware(logger logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Generate trace ID for this request
		traceID := logger.GetTraceID()
		if existingTraceID := c.GetHeader("X-Trace-ID"); existingTraceID != "" {
			traceID = existingTraceID
		}

		// Create request-scoped logger with trace ID
		requestLogger := logger.WithTraceID(traceID)

		// Store logger in context
		c.Set("logger", requestLogger)

		// Set trace ID header in response
		c.Header("X-Trace-ID", traceID)

		// Read request body if it's a POST, PUT, or PATCH request
		var requestBody []byte
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			if c.Request.Body != nil {
				requestBody, _ = io.ReadAll(c.Request.Body)
				// Reset the request body so it can be read again by handlers
				c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
			}
		}

		// Create a custom response writer to capture response body and status
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// Log request details
		requestLogger.Info("Request received", map[string]interface{}{
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"query":      c.Request.URL.RawQuery,
			"ip":         c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
			"body":       string(requestBody),
		})

		// Process request
		c.Next()

		// Get status code and error messages after request is processed
		statusCode := c.Writer.Status()
		errors := c.Errors.Errors()

		// Calculate request duration
		duration := time.Since(start)

		// Log response details based on status code
		logFields := map[string]interface{}{
			"status_code": statusCode,
			"duration_ms": duration.Milliseconds(),
			"path":        c.Request.URL.Path,
			"method":      c.Request.Method,
		}

		// For APIs returning less sensitive data, we might want to log response bodies
		// Don't log response bodies for files or large responses
		if blw.body.Len() < 1024 && c.Writer.Header().Get("Content-Type") != "application/octet-stream" {
			logFields["response_body"] = blw.body.String()
		}

		if len(errors) > 0 {
			logFields["errors"] = errors
		}

		// Log based on status code
		if statusCode >= 500 {
			requestLogger.Error("Server error response", logFields)
		} else if statusCode >= 400 {
			requestLogger.Warn("Client error response", logFields)
		} else {
			requestLogger.Info("Request completed successfully", logFields)
		}
	}
}

// bodyLogWriter is a custom gin.ResponseWriter that captures the response body
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write captures the response body for logging
func (w *bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
