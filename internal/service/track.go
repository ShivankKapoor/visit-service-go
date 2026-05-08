package service

import (
	"context"
	"log/slog"
	"net/http"
	"visit-service/internal/models"
	"visit-service/internal/repositories"

	"github.com/gin-gonic/gin"
)

func Track(c *gin.Context, visit models.PageVisit) {
	err := repositories.InsertPageVisit(context.Background(), visit)
	if err != nil {
		slog.Error("Failed to insert page visit", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record visit"})
		return
	}

	slog.Info("Visit recorded", "ip", visit.IPAddress, "page", visit.PageVisited, "device", visit.DeviceInfo)

	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}
