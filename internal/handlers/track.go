package handlers

import (
	"log/slog"
	"net/http"
	"visit-service/internal/models"
	"visit-service/internal/network"

	"github.com/gin-gonic/gin"
)

func Track(c *gin.Context) {
	var req models.TrackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	clientIP := network.GetClientIP(c)
	DeviceInfo := req.DeviceInfo
	PageVisited := req.PageVisited

	slog.Info("\nVisit from \n IP: " + clientIP + "\n Device Info: " + DeviceInfo + "\nPage Visited: " + PageVisited)

	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}
