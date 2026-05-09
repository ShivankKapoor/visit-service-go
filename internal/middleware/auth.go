package middleware

import (
	"net/http"
	"os"
	"strings"
	"visit-service/internal/service"

	"github.com/gin-gonic/gin"
)

func AllowedReferer() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !service.IsProd() {
			c.Next()
			return
		}

		allowedOriginsStr := os.Getenv("ALLOWED_ORIGINS")
		referer := c.GetHeader("Referer")

		if referer == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			c.Abort()
			return
		}

		allowedOrigins := strings.Split(allowedOriginsStr, ",")
		allowed := false
		for _, origin := range allowedOrigins {
			if strings.Contains(referer, origin) {
				allowed = true
				break
			}
		}

		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			c.Abort()
			return
		}

		c.Next()
	}
}
