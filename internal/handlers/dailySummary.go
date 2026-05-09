package handlers

import (
	"os"
	"strconv"
	"visit-service/internal/service"

	"github.com/gin-gonic/gin"
)

func TriggerDailySummary(c *gin.Context) {
	isProdStr := os.Getenv("PROD")
	isProd, err := strconv.ParseBool(isProdStr)
	if err != nil {
		isProd = true
	}

	if !isProd {
		service.RunDailySummary()
		c.JSON(200, gin.H{"message": "Daily summary triggered"})
	} else {
		c.JSON(401, gin.H{"message": "Unauthorized"})
	}
}
