package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"
	"visit-service/internal/models"
	"visit-service/internal/network"
	"visit-service/internal/repositories"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Track(c *gin.Context) {
	var req models.TrackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	clientIP := network.GetClientIP(c)
	userAgent := c.Request.Header.Get("User-Agent")

	deviceInfo := req.DeviceInfo
	if idx := strings.Index(deviceInfo, ","); idx != -1 {
		deviceInfo = deviceInfo[:idx]
	}

	visit := repositories.PageVisit{
		ID:          uuid.New().String(),
		IPAddress:   clientIP,
		PageVisited: req.PageVisited,
		DeviceInfo:  &deviceInfo,
		UserAgent:   &userAgent,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
	}

	err := repositories.InsertPageVisit(context.Background(), visit)
	if err != nil {
		slog.Error("Failed to insert page visit", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record visit"})
		return
	}

	slog.Info("Visit recorded", "ip", clientIP, "page", req.PageVisited, "device", req.DeviceInfo)

	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}
