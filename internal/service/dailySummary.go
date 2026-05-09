package service

import (
	"context"
	"log/slog"
	"time"
	"visit-service/internal/models"
	"visit-service/internal/repositories"

	"github.com/robfig/cron/v3"
)

var cronScheduler *cron.Cron

func RunDailySummary() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	slog.Info("Starting daily summary job")

	count, err := repositories.GetYesterdayVisitCount(ctx)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			slog.Error("Daily summary job timeout: failed to get visit count", "timeout", "30s")
		} else {
			slog.Error("Error getting daily summary", "error", err)
		}
		return
	}

	cst, err := time.LoadLocation("America/Chicago")
	if err != nil {
		slog.Error("Error loading CST timezone", "error", err)
		return
	}

	now := time.Now().In(cst)
	yesterday := now.AddDate(0, 0, -1)
	summaryDate := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, time.UTC)

	stats := models.DailyVisitStats{
		SummaryDate: summaryDate,
		TotalVisits: count,
	}

	err = repositories.SaveDailyStats(ctx, stats)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			slog.Error("Daily summary job timeout: failed to save stats", "timeout", "30s")
		} else {
			slog.Error("Error saving daily stats", "error", err)
		}
		return
	}

	SendDailyVisitsMessage(int(count))
	slog.Info("Daily report saved", "visits", count)
}

func SetupDailyReportTask() error {
	cronScheduler = cron.New(cron.WithLocation(
		func() *time.Location {
			loc, _ := time.LoadLocation("America/Chicago")
			return loc
		}(),
	))

	_, err := cronScheduler.AddFunc("0 0 * * *", RunDailySummary)
	if err != nil {
		return err
	}

	cronScheduler.Start()
	slog.Info("Daily report cron task started")
	return nil
}

func StopDailyReportTask() {
	if cronScheduler != nil {
		cronScheduler.Stop()
		slog.Info("Daily report cron task stopped")
	}
}
