package service

import (
	"context"
	"log/slog"
	"time"
	"visit-service/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetDBHealth(db *pgxpool.Pool) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := repository.NewHealthRepository(db).Ping(ctx); err != nil {
		slog.Error("DB health check failed", "err", err)
		return false
	}
	return true
}
