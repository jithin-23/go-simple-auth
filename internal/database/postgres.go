package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"simple-auth/internal/config"
)

// Connect opens a GORM connection pool to Postgres and verifies it with a ping.
func Connect(cfg config.DBConfig) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.New(log.New(os.Stdout, "\n", log.LstdFlags), logger.Config{
			SlowThreshold:             200 * time.Millisecond, // log queries slower than this
			LogLevel:                  logger.Warn,            // only log slow queries and errors
			IgnoreRecordNotFoundError: true,                   // "not found" is normal, not an error
			Colorful:                  true,
		}),
		TranslateError: true, // turn Postgres errors into gorm.ErrDuplicatedKey etc.
	})
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	// GORM wraps Go's standard *sql.DB, which owns the connection pool.
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql.DB: %w", err)
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return db, nil
}
