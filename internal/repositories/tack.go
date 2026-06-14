package repositories

import (
	"context"
	"visit-service/internal/database"
	"visit-service/internal/models"
)

func InsertPageVisit(ctx context.Context, visit models.PageVisit) error {
	query := `
		INSERT INTO page_visits (id, ip_address, page_visited, device_info, user_agent, timestamp)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := database.DB.Exec(ctx, query,
		visit.ID,
		visit.IPAddress,
		visit.PageVisited,
		visit.DeviceInfo,
		visit.UserAgent,
		visit.Timestamp,
	)

	return err
}
