package repositories

import (
	"context"
	"visit-service/internal/database"

	"github.com/jackc/pgx/v5"
)

type PageVisit struct {
	ID          string
	IPAddress   string
	PageVisited string
	DeviceInfo  *string
	UserAgent   *string
	Timestamp   string
}

func InsertPageVisit(ctx context.Context, visit PageVisit) error {
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

// Bulk insert for multiple records
func InsertPageVisitsBatch(ctx context.Context, visits []PageVisit) error {
	batch := &pgx.Batch{}

	for _, visit := range visits {
		batch.Queue(
			`INSERT INTO page_visits (id, ip_address, page_visited, device_info, user_agent, timestamp)
			 VALUES ($1, $2, $3, $4, $5, $6)`,
			visit.ID,
			visit.IPAddress,
			visit.PageVisited,
			visit.DeviceInfo,
			visit.UserAgent,
			visit.Timestamp,
		)
	}

	results := database.DB.SendBatch(ctx, batch)
	defer results.Close()

	for range visits {
		_, err := results.Exec()
		if err != nil {
			return err
		}
	}

	return nil
}
