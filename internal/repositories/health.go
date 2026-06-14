package repositories

import (
	"context"
	"visit-service/internal/database"
)

func GetHealth(ctx context.Context) error {
	query := `
		SELECT 1
	`
	_, err := database.DB.Exec(ctx, query)

	return err
}
