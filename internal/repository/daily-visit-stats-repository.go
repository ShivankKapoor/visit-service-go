package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DailyVisitStatsRepository struct {
	db *pgxpool.Pool
}

func NewDailyVisitStatsRepository(db *pgxpool.Pool) *DailyVisitStatsRepository {
	return &DailyVisitStatsRepository{db: db}
}

func (r *DailyVisitStatsRepository) CountVisitsForDay(ctx context.Context, dayStart, dayEnd time.Time) (int64, error) {
	var count int64
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM page_visits WHERE timestamp >= $1 AND timestamp < $2`,
		dayStart, dayEnd,
	).Scan(&count)
	return count, err
}

func (r *DailyVisitStatsRepository) InsertDailyVisitStat(ctx context.Context, date time.Time, totalVisits int64) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO daily_visit_stats (summary_date, total_visits)
		 VALUES ($1, $2)
		 ON CONFLICT (summary_date) DO UPDATE SET total_visits = EXCLUDED.total_visits`,
		date, totalVisits,
	)
	return err
}
