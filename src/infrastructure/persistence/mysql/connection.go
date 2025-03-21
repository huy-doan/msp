package mysql

import (
	"fmt"
	"net/url"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/vnlab/makeshop-payment/src/infrastructure/config"
	"github.com/vnlab/makeshop-payment/src/infrastructure/logger"
)

const (
	CONN_MAX_LIFETIME = time.Minute * 10
	MAX_IDLE_CONNS    = 500
	MAX_OPEN_CONNS    = 250
)

// NewConnection creates a new MySQL database connection using GORM
func NewConnection(appConfig *config.Config, appLogger logger.Logger) (*gorm.DB, error) {
	dbHost     := appConfig.DBHost
	dbPort     := appConfig.DBPort
	dbUser     := appConfig.DBUser
	dbPassword := appConfig.DBPassword
	dbName     := appConfig.DBName

	// Configure connection string with Tokyo timezone
	loc := url.QueryEscape("Asia/Tokyo")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=%s",
		dbUser, dbPassword, dbHost, dbPort, dbName, loc)

	// Create SQL logger that integrates with our custom logger
	sqlLogger := logger.NewSQLLogger(&logger.Config{
		LogLevel:      appConfig.LogLevel,
		LogDirectory:  appConfig.LogDirectory,
		EnableConsole: appConfig.EnableConsole,
		EnableSQLLog:  appConfig.EnableSQLLog,
	}, appLogger)

	// Open database connection with our custom SQL logger
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: sqlLogger,
	})

	if err != nil {
		// appLogger.Error("Failed to connect to database", map[string]interface{}{
		// 	"error": err.Error(),
		// 	"host":  dbHost,
		// 	"port":  dbPort,
		// 	"user":  dbUser,
		// 	"name":  dbName,
		// })
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		// appLogger.Error("Failed to get SQL DB handle", map[string]interface{}{
		// 	"error": err.Error(),
		// })
		return nil, fmt.Errorf("failed to get SQL DB: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(MAX_IDLE_CONNS)
	sqlDB.SetMaxOpenConns(MAX_OPEN_CONNS)
	sqlDB.SetConnMaxLifetime(CONN_MAX_LIFETIME)

	// appLogger.Info("Database connection established", map[string]interface{}{
	// 	"host":             dbHost,
	// 	"name":             dbName,
	// 	"max_idle_conns":   MAX_IDLE_CONNS,
	// 	"max_open_conns":   MAX_OPEN_CONNS,
	// 	"conn_max_lifetime": CONN_MAX_LIFETIME.String(),
	// })

	return db, nil
}
