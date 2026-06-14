package handlers

import (
	"visit-service/internal/service"

	"github.com/gin-gonic/gin"
)

func TriggerDailySummary(c *gin.Context) {

	if !service.IsProd() {
		service.RunDailySummary()
		c.JSON(200, gin.H{"message": "Daily summary triggered"})
	} else {
		c.JSON(401, gin.H{"message": "Unauthorized"})
	}
}
