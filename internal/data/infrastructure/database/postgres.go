// db/postgres.go

package database

import (
	"database/sql"
	"fmt"
	"github.com/craftizmv/rewards/pkg/logger"
	_ "github.com/lib/pq"
)

func NewPostgresDB(cfg *Config, log logger.ILogger) (*sql.DB, error) {
	// Construct the connection string
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	// Open the database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Error("Failed to open PostgreSQL connection:", err)
		return nil, err
	}

	// Set connection pool parameters
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// Verify the connection with a Ping
	if err := db.Ping(); err != nil {
		log.Error("Failed to ping PostgreSQL:", err)
		return nil, err
	}

	log.Info("Successfully connected to PostgreSQL")
	return db, nil
}
