package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"visit-service/internal/db"
	"visit-service/internal/handler"
	"visit-service/internal/middleware"
	"visit-service/internal/service"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	slog.Info("Starting the visit service v2")

	if os.Getenv("PROD") == "true" && os.Getenv("ALLOWED_ORIGINS") == "" {
		slog.Error("ALLOWED_ORIGINS must be set in production")
		os.Exit(1)
	}

	pool, err := db.New()
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	trackService := service.NewTrackService(pool)
	service.StartDailyCron(pool)
	defer service.StopDailyCron()

	MainHandler := handler.NewMainHandler()
	TrackHandler := handler.NewTrackHandler(trackService)
	CronHandler := handler.NewCronHandler(pool)
	HealthHandler := handler.NewHealthHandler(pool)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", MainHandler.Home)
	mux.Handle("POST /track", middleware.AllowedReferer(http.HandlerFunc(TrackHandler.Track)))
	if os.Getenv("PROD") != "true" {
		mux.HandleFunc("GET /run-daily-summary", CronHandler.RunDailySummary)
	}
	mux.HandleFunc("GET /health", HealthHandler.Health)

	srv := &http.Server{
		Addr:              ":8088",
		Handler:           middleware.Cors(middleware.RateLimit(mux)),
		ReadHeaderTimeout: 20 * time.Second,
		ReadTimeout:       20 * time.Second,
		WriteTimeout:      20 * time.Second,
		IdleTimeout:       20 * time.Second,
	}

	go func() {
		slog.Info("Meridian web engine starting...", "port", ":8088")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Failed to start the web server", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Graceful shutdown failed", "error", err)
		os.Exit(1)
	}
	slog.Info("Server stopped")
}
