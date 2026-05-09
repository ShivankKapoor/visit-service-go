package handlers

import (
	"net/http"
	"visit-service/internal/models"
	"visit-service/internal/service"

	"github.com/gin-gonic/gin"
)

func Health(c *gin.Context) {
	dbHealth := "Unhealthy"
	responseHealth := http.StatusInternalServerError
	if service.GetDBHealth() {
		dbHealth = "Health"
		responseHealth = http.StatusOK
	}

	reponse := models.HealthResponse{
		Name:      "Visit Service",
		Platform:  "GO",
		APIStatus: "UP",
		DBStatus:  dbHealth,
	}

	c.JSON(responseHealth, reponse)
}
