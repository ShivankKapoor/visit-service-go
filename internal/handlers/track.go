package handlers

import (
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

	service.Track(c, visit)
}
