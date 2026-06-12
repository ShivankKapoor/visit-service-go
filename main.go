package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"visit-service/internal/handler"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	slog.Info("Starting the visit service")

	MainHandler := handler.NewMainHandler()
	TrackHandler := handler.NewTrackHandler()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", MainHandler.Home)
	mux.HandleFunc("POST /track", TrackHandler.Track)

	srv := &http.Server{Addr: ":8088", Handler: mux}

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
