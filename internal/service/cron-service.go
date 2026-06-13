package service

import (
	"context"
	"log/slog"
	"time"
	"visit-service/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

var cst = mustLoadLocation("America/Chicago")

func mustLoadLocation(name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		panic(err)
	}
	return loc
}

func StartDailyCron(db *pgxpool.Pool) {
	go func() {
		for {
			now := time.Now().In(cst)
			nextMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, cst)
			time.Sleep(time.Until(nextMidnight))
			RunDailySummary(db)
		}
	}()
}

func RunDailySummary(db *pgxpool.Pool) {
	now := time.Now().In(cst)
	yesterdayStart := time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, cst)
	yesterdayEnd := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, cst)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	repo := repository.NewDailyVisitStatsRepository(db)

	count, err := repo.CountVisitsForDay(ctx, yesterdayStart, yesterdayEnd)
	if err != nil {
		slog.Error("Failed to count daily visits", "err", err)
		return
	}

	if err := repo.InsertDailyVisitStat(ctx, yesterdayStart, count); err != nil {
		slog.Error("Failed to insert daily visit stat", "err", err)
		return
	}

	slog.Info("Daily visit summary saved", "date", yesterdayStart.Format("2006-01-02"), "visits", count)
	SendDailyVisitsMessage(int(count))
}
