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
	locationReq, err := GetLocation(visit.IPAddress)
	var location = "unknown"
	if err == nil {
		location = locationReq.City + ", " + locationReq.RegionName + ", " + locationReq.Country
	}
	slog.Info("Visit recorded", "IP:", visit.IPAddress, "Page:", visit.PageVisited, "Device:", visit.DeviceInfo, "Location:", location)

	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}
