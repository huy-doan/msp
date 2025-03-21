package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all the configuration for the application
type Config struct {
	// Server configuration
	ServerHost string
	ServerPort string
	
	// Database configuration
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	
	GinMode string // Gin mode for the server

	// Logger configuration
	LogLevel      string
	LogDirectory  string // Directory where log files will be stored
	EnableConsole bool   // Whether to also log to console
	EnableSQLLog  bool   // Whether to log SQL queries

	// Authentication configuration
	JWTSecret   string
	JWTDuration int
}

// LoadConfig loads the configuration from environment variables
func LoadConfig() *Config {
	// Load .env file if it exists
	godotenv.Load()
	
	// Set default values
	config := &Config{
		ServerHost:    "0.0.0.0",
		ServerPort:    "8080",
		LogLevel:      "info",
		LogDirectory:  "./logs",
		EnableConsole: true,
		EnableSQLLog:  false,
		JWTDuration:   24, // Hours
	}
	
	// Map of environment variables to configuration fields
	envVars := map[string]*string{
		"SERVER_HOST":    &config.ServerHost,
		"SERVER_PORT":    &config.ServerPort,
		"GIN_MODE":       &config.GinMode,
		"DB_HOST":        &config.DBHost,
		"DB_PORT":        &config.DBPort,
		"DB_USER":        &config.DBUser,
		"DB_PASSWORD":    &config.DBPassword,
		"DB_NAME":        &config.DBName,
		"LOG_LEVEL":      &config.LogLevel,
		"LOG_DIRECTORY":  &config.LogDirectory,
		"JWT_SECRET":     &config.JWTSecret,
	}

	// Override string fields with environment variables if they exist
	for env, field := range envVars {
		if val := os.Getenv(env); val != "" {
			*field = val
		}
	}

	// Override boolean fields
	boolVars := map[string]*bool{
		"ENABLE_CONSOLE": &config.EnableConsole,
		"ENABLE_SQL_LOG": &config.EnableSQLLog,
	}
	for env, field := range boolVars {
		if val := os.Getenv(env); val != "" {
			parsedVal, err := strconv.ParseBool(val)
			if err == nil {
				*field = parsedVal
			}
		}
	}

	// Override integer fields
	if val := os.Getenv("JWT_DURATION"); val != "" {
		if duration, err := strconv.Atoi(val); err == nil {
			config.JWTDuration = duration
		}
	}
	return config
}

// GetLoggerConfig returns logger configuration
func (c *Config) GetLoggerConfig() map[string]interface{} {
	return map[string]interface{}{
		"log_level":       c.LogLevel,
		"log_directory":   c.LogDirectory,
		"enable_console":  c.EnableConsole,
		"enable_sql_log":  c.EnableSQLLog,
	}
}
