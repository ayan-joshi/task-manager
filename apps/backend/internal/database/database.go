package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"taskmanager/internal/domain"
)

// Connect opens a PostgreSQL connection pool and verifies connectivity.
func Connect(dsn string, production bool) (*gorm.DB, error) {
	logLevel := logger.Info
	if production {
		logLevel = logger.Warn
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                 logger.Default.LogMode(logLevel),
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("access sql.DB: %w", err)
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}
	return db, nil
}

// Migrate runs GORM auto-migration for every domain model. Auto-migration is
// sufficient for this application's evolving schema; for stricter production
// change control these statements can be exported to versioned SQL migrations.
func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&domain.User{},
		&domain.Task{},
		&domain.Attachment{},
		&domain.ActivityLog{},
	); err != nil {
		return fmt.Errorf("auto-migrate: %w", err)
	}
	log.Println("database migration complete")
	return nil
}
