package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"
	"visit-service/internal/models"
	"visit-service/internal/network"
	"visit-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Track(c *gin.Context) {
	var req models.TrackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validateTrackRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	clientIP := network.GetClientIP(c)
	userAgent := c.Request.Header.Get("User-Agent")

	deviceInfo := req.DeviceInfo
	if idx := strings.Index(deviceInfo, ","); idx != -1 {
		deviceInfo = deviceInfo[:idx]
	}

	visit := models.PageVisit{
		ID:          uuid.New().String(),
		IPAddress:   clientIP,
		PageVisited: req.PageVisited,
		DeviceInfo:  &deviceInfo,
		UserAgent:   &userAgent,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
	}

	go service.TrackAsync(visit)

	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}

func validateTrackRequest(req models.TrackRequest) error {
	if req.PageVisited == "" {
		return fmt.Errorf("pageVisited is required")
	}

	if len(req.PageVisited) > 2048 {
		return fmt.Errorf("pageVisited exceeds maximum length of 2048")
	}

	if len(req.DeviceInfo) > 500 {
		return fmt.Errorf("deviceInfo exceeds maximum length of 500")
	}

	return nil
}
