package handlers

import (
	"net/http"
	"visit-service/internal/models"

	"github.com/gin-gonic/gin"
)

func HealthCheck(c *gin.Context) {

	reponse := models.HealthCheckResponse{
		Name:     "Visit Service",
		Platform: "Go",
		Status:   "Up",
	}

	c.JSON(http.StatusOK, reponse)
}
