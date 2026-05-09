package repositories

import (
	"context"
	"log/slog"
	"time"
	"visit-service/internal/database"
	"visit-service/internal/models"
)

func GetYesterdayVisitCount(ctx context.Context) (int64, error) {
	cst, err := time.LoadLocation("America/Chicago")
	if err != nil {
		return 0, err
	}

	now := time.Now().In(cst)
	yesterday := now.AddDate(0, 0, -1)

	startOfYesterday := time.Date(
		yesterday.Year(), yesterday.Month(), yesterday.Day(),
		0, 0, 0, 0, cst,
	)

	endOfYesterday := time.Date(
		now.Year(), now.Month(), now.Day(),
		0, 0, 0, 0, cst,
	)

	var count int64
	err = database.DB.QueryRow(ctx,
		"SELECT COUNT(*) FROM public.page_visits WHERE timestamp >= $1 AND timestamp < $2",
		startOfYesterday.UTC(), endOfYesterday.UTC(),
	).Scan(&count)

	if err != nil {
		slog.Error("Failed to count yesterday's visits", "error", err)
		return 0, err
	}

	return count, nil
}

func SaveDailyStats(ctx context.Context, stats models.DailyVisitStats) error {
	_, err := database.DB.Exec(ctx,
		`INSERT INTO public.daily_visit_stats (summary_date, total_visits)
		 VALUES ($1, $2)
		 ON CONFLICT (summary_date) DO UPDATE
		 SET total_visits = EXCLUDED.total_visits`,
		stats.SummaryDate, stats.TotalVisits,
	)

	if err != nil {
		slog.Error("Failed to save daily stats", "error", err)
		return err
	}

	return nil
}
