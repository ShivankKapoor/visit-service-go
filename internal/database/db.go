package database

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func Init() error {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return fmt.Errorf("DATABASE_URL not set")
	}

	slog.Info("Initializing database connection")

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return fmt.Errorf("unable to parse database config: %w", err)
	}

	config.ConnConfig.ConnectTimeout = 30 * time.Second

	var dbErr error
	DB, dbErr = pgxpool.NewWithConfig(context.Background(), config)
	if dbErr != nil {
		return fmt.Errorf("unable to create connection pool: %w", dbErr)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := DB.Ping(ctx); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	slog.Info("Database connection established successfully")
	return nil
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}
