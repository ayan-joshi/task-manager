package database

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"taskmanager/internal/domain"
)

// Connect opens a database connection pool and verifies connectivity.
//
// The driver is selected from the DSN: anything beginning with "file:" or
// "sqlite:" (or ending in .db) uses the pure-Go SQLite driver, which is handy
// for zero-dependency local development. Everything else uses PostgreSQL, the
// production database. The domain models are written to work on both.
func Connect(dsn string, production bool) (*gorm.DB, error) {
	logLevel := logger.Info
	if production {
		logLevel = logger.Warn
	}

	dialector := postgres.Open(dsn)
	usingSQLite := isSQLiteDSN(dsn)
	if usingSQLite {
		dialector = sqlite.Open(dsn)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
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
	if usingSQLite {
		// SQLite permits a single writer; a single connection avoids
		// "database is locked" errors under concurrent requests.
		sqlDB.SetMaxOpenConns(1)
	} else {
		sqlDB.SetMaxOpenConns(25)
		sqlDB.SetMaxIdleConns(5)
		sqlDB.SetConnMaxLifetime(30 * time.Minute)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}
	return db, nil
}

// isSQLiteDSN reports whether the DSN should be served by the SQLite driver.
func isSQLiteDSN(dsn string) bool {
	lower := strings.ToLower(dsn)
	return strings.HasPrefix(lower, "file:") ||
		strings.HasPrefix(lower, "sqlite:") ||
		strings.HasSuffix(strings.SplitN(lower, "?", 2)[0], ".db")
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
