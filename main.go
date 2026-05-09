package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
	"visit-service/internal/database"
	"visit-service/internal/middleware"
	"visit-service/internal/routes"
	"visit-service/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if service.IsProd() {
		gin.SetMode(gin.ReleaseMode)
	}
	godotenv.Load()

	if err := database.Init(); err != nil {
		slog.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}
	defer database.Close()

	if err := service.SetupDailyReportTask(); err != nil {
		slog.Error("Failed to setup daily report task", "error", err)
		os.Exit(1)
	}
	defer service.StopDailyReportTask()

	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowMethods = []string{"GET", "POST", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Content-Type", "Authorization"}

	if service.IsProd() {
		allowedOriginsStr := os.Getenv("ALLOWED_ORIGINS")
		if allowedOriginsStr == "" {
			slog.Error("ALLOWED_ORIGINS not set in production")
			os.Exit(1)
		}
		allowedOrigins := strings.Split(allowedOriginsStr, ",")
		corsConfig.AllowOrigins = allowedOrigins
		slog.Info("CORS restricted", "origins", allowedOrigins)
	} else {
		corsConfig.AllowAllOrigins = true
		slog.Info("CORS allows all origins (dev mode)")
	}

	r.Use(cors.New(corsConfig))

	r.Use(middleware.RateLimit())

	routes.Configure(r)

	server := &http.Server{
		Addr:         ":8088",
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		slog.Info("Server starting", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server error", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("Server stopped")
}
