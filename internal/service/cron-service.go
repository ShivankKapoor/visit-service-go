package service

import (
	"context"
	"log/slog"
	"time"
	"visit-service/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

var cst = mustLoadLocation("America/Chicago")

var stopCron = make(chan struct{})

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
			slog.Info("Daily cron scheduled", "next_run", nextMidnight.Format("2006-01-02 15:04:05 MST"))
			select {
			case <-time.After(time.Until(nextMidnight)):
				RunDailySummary(db)
			case <-stopCron:
				slog.Info("Daily cron stopped")
				return
			}
		}
	}()
}

func StopDailyCron() {
	close(stopCron)
}

func RunDailySummary(db *pgxpool.Pool) {
	now := time.Now().In(cst)
	yesterdayStart := time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, cst)
	yesterdayEnd := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, cst)

	repo := repository.NewDailyVisitStatsRepository(db)

	countCtx, countCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer countCancel()
	count, err := repo.CountVisitsForDay(countCtx, yesterdayStart, yesterdayEnd)
	if err != nil {
		slog.Error("Failed to count daily visits", "err", err)
		return
	}

	insertCtx, insertCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer insertCancel()
	if err := repo.InsertDailyVisitStat(insertCtx, yesterdayStart, count); err != nil {
		slog.Error("Failed to insert daily visit stat", "err", err)
		return
	}

	slog.Info("Daily visit summary saved", "date", yesterdayStart.Format("2006-01-02"), "visits", count)
	go SendDailyVisitsMessage(int(count))
}
