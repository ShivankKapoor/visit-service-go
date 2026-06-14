package service

import (
	"context"
	"log/slog"
	"time"
	"visit-service/internal/repositories"
)

func GetDBHealth() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := repositories.GetHealth(ctx)

	if err != nil {
		slog.Error("Error when checking DB Health", "error", err)
		return false
	}

	return true

}
