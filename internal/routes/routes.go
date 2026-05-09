package routes

import (
	"visit-service/internal/handlers"
	"visit-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Configure(r *gin.Engine) {
	v1 := r.Group("")
	{
		v1.GET("/", handlers.Home)
		v1.POST("/track", middleware.AllowedReferer(), handlers.Track)
		v1.GET("/health", handlers.Health)
		v1.GET("/admin/run-summary", handlers.TriggerDailySummary)
	}
}
