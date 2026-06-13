package repository

import (
	"context"
	"visit-service/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PageVisitRepository struct {
	db *pgxpool.Pool
}

func NewPageVisitRepository(db *pgxpool.Pool) *PageVisitRepository {
	return &PageVisitRepository{db: db}
}

func (r *PageVisitRepository) InsertPageVisit(ctx context.Context, visit model.PageVisit) error {
	query := `
		INSERT INTO page_visits (id, ip_address, page_visited, device_info, user_agent, timestamp)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(ctx, query,
		visit.ID,
		visit.IPAddress,
		visit.PageVisited,
		visit.DeviceInfo,
		visit.UserAgent,
		visit.Timestamp,
	)

	return err
}
