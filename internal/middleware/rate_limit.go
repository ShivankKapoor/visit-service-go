package middleware

import (
	"net/http"
	"sync"
	"time"
	"visit-service/internal/network" // Import your IP helper

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type IPRateLimiter struct {
	ips map[string]*rate.Limiter
	mu  sync.RWMutex
}

var limiter = IPRateLimiter{ips: make(map[string]*rate.Limiter)}

func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIP := network.GetClientIP(c)

		limiter.mu.Lock()
		if _, exists := limiter.ips[userIP]; !exists {
			limiter.ips[userIP] = rate.NewLimiter(rate.Every(time.Minute/20), 20)
		}
		l := limiter.ips[userIP]
		limiter.mu.Unlock()

		if !l.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "too many requests",
				"ip":    userIP,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
