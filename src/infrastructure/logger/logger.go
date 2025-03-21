// src/infrastructure/logger/logger.go
package logger

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// LogLevel represents logging levels
type LogLevel string

const (
	// DEBUG level for detailed information in development environment
	DEBUG LogLevel = "debug"
	// INFO level for general operational information
	INFO LogLevel = "info"
	// WARN level for warnings that don't cause errors but should be noted
	WARN LogLevel = "warn"
	// ERROR level for system errors
	ERROR LogLevel = "error"
)

// Logger is the custom structured logger interface
type Logger interface {
	Debug(msg string, fields map[string]interface{})
	Info(msg string, fields map[string]interface{})
	Warn(msg string, fields map[string]interface{})
	Error(msg string, fields map[string]interface{})
	WithTraceID(traceID string) Logger
	GetTraceID() string
}

// loggerImpl is the implementation of Logger interface
type loggerImpl struct {
	logger  *logrus.Logger
	traceID string
}

// Config holds the configuration for logger
type Config struct {
	LogLevel      string
	LogDirectory   string
	EnableConsole bool
	EnableSQLLog  bool
}

// NewLogger creates a new logger instance
func NewLogger(config *Config) Logger {
	logger := logrus.New()

	// Set log level
	level, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	// Ensure log directory exists
	logDir := filepath.Dir(config.LogDirectory)
	err = os.MkdirAll(logDir, 0755)
	if err != nil {
		panic(err)
	}

	// Configure JSON formatter for structured logging
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})
	
	// We'll handle different log levels in our custom methods
	if config.EnableConsole {
		logger.SetOutput(os.Stdout)
	} else {
		logger.SetOutput(io.Discard) // Discard default output as we'll use level-specific files
	}

	return &loggerImpl{
		logger:  logger,
		traceID: GenerateTraceID(),
	}
}

// GenerateTraceID generates a unique trace ID for request tracking
func GenerateTraceID() string {
	return uuid.New().String()
}

// WithTraceID creates a new logger instance with the specified trace ID
func (l *loggerImpl) WithTraceID(traceID string) Logger {
	return &loggerImpl{
		logger:  l.logger,
		traceID: traceID,
	}
}

// GetTraceID returns the current trace ID
func (l *loggerImpl) GetTraceID() string {
	return l.traceID
}

// makeFields adds common fields to all log entries
func (l *loggerImpl) makeFields(fields map[string]interface{}) logrus.Fields {
	if fields == nil {
		fields = make(map[string]interface{})
	}
	
	fields["trace_id"] = l.traceID
	fields["timestamp"] = time.Now().UTC().Format(time.RFC3339)
	
	return logrus.Fields(fields)
}

// getLogFile returns the appropriate log file for the given level
func (l *loggerImpl) getLogFile(level logrus.Level) *os.File {
	var logPath string
	logDir := filepath.Dir(l.logger.Out.(*os.File).Name())
	
	switch level {
	case logrus.DebugLevel:
		logPath = filepath.Join(logDir, "debug.log")
	case logrus.InfoLevel:
		logPath = filepath.Join(logDir, "info.log")
	case logrus.WarnLevel:
		logPath = filepath.Join(logDir, "warning.log")
	case logrus.ErrorLevel:
		logPath = filepath.Join(logDir, "error.log")
	default:
		logPath = filepath.Join(logDir, "app.log")
	}
	
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		// If we can't open the log file, fallback to stdout
		l.logger.WithField("error", err.Error()).Error("Failed to open log file")
		return os.Stdout
	}
	
	return file
}

// logToFile logs a message to the appropriate file based on level
func (l *loggerImpl) logToFile(level logrus.Level, entry *logrus.Entry) {
	if !l.logger.IsLevelEnabled(level) {
		return
	}
	
	// Get the appropriate file for this log level
	file := l.getLogFile(level)
	defer file.Close()
	
	// Create a new logger for this specific write
	fileLogger := logrus.New()
	fileLogger.SetOutput(file)
	fileLogger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})
	fileLogger.SetLevel(level)
	
	// Write the log entry to the file
	switch level {
	case logrus.DebugLevel:
		fileLogger.WithFields(entry.Data).Debug(entry.Message)
	case logrus.InfoLevel:
		fileLogger.WithFields(entry.Data).Info(entry.Message)
	case logrus.WarnLevel:
		fileLogger.WithFields(entry.Data).Warn(entry.Message)
	case logrus.ErrorLevel:
		fileLogger.WithFields(entry.Data).Error(entry.Message)
	}
}

// Debug logs a message at the DEBUG level
func (l *loggerImpl) Debug(msg string, fields map[string]interface{}) {
	if l.logger.IsLevelEnabled(logrus.DebugLevel) {
		entry := &logrus.Entry{
			Logger:  l.logger,
			Data:    l.makeFields(fields),
			Time:    time.Now(),
			Level:   logrus.DebugLevel,
			Message: msg,
		}
		
		// Write to the console if enabled
		l.logger.WithFields(entry.Data).Debug(msg)
		
		// Write to the appropriate log file
		l.logToFile(logrus.DebugLevel, entry)
	}
}

// Info logs a message at the INFO level
func (l *loggerImpl) Info(msg string, fields map[string]interface{}) {
	if l.logger.IsLevelEnabled(logrus.InfoLevel) {
		entry := &logrus.Entry{
			Logger:  l.logger,
			Data:    l.makeFields(fields),
			Time:    time.Now(),
			Level:   logrus.InfoLevel,
			Message: msg,
		}
		
		// Write to the console if enabled
		l.logger.WithFields(entry.Data).Info(msg)
		
		// Write to the appropriate log file
		l.logToFile(logrus.InfoLevel, entry)
	}
}

// Warn logs a message at the WARN level
func (l *loggerImpl) Warn(msg string, fields map[string]interface{}) {
	if l.logger.IsLevelEnabled(logrus.WarnLevel) {
		entry := &logrus.Entry{
			Logger:  l.logger,
			Data:    l.makeFields(fields),
			Time:    time.Now(),
			Level:   logrus.WarnLevel,
			Message: msg,
		}
		
		// Write to the console if enabled
		l.logger.WithFields(entry.Data).Warn(msg)
		
		// Write to the appropriate log file
		l.logToFile(logrus.WarnLevel, entry)
	}
}

// Error logs a message at the ERROR level
func (l *loggerImpl) Error(msg string, fields map[string]interface{}) {
	if l.logger.IsLevelEnabled(logrus.ErrorLevel) {
		entry := &logrus.Entry{
			Logger:  l.logger,
			Data:    l.makeFields(fields),
			Time:    time.Now(),
			Level:   logrus.ErrorLevel,
			Message: msg,
		}
		
		// Write to the console if enabled
		l.logger.WithFields(entry.Data).Error(msg)
		
		// Write to the appropriate log file
		l.logToFile(logrus.ErrorLevel, entry)
	}
}
