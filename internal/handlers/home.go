package handlers

import (
	"net/http"
	"visit-service/internal/models"

	"github.com/gin-gonic/gin"
)

func Home(c *gin.Context) {

	reponse := models.HomeResponse{
		Name:     "Visit Service",
		Platform: "Go",
		Status:   "Up",
	}

	c.JSON(http.StatusOK, reponse)
}
