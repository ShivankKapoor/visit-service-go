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
)

func main() {
	slog.Info("Starting the visit service")

	MainHandler := handler.NewMainHandler()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", MainHandler.Home)

	srv := &http.Server{Addr: ":8080", Handler: mux}

	go func() {
		slog.Info("Meridian web engine starting...", "port", ":8080")
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
