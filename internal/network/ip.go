package network

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// GetClientIP extracts the real user IP, prioritizing Cloudflare headers
func GetClientIP(c *gin.Context) string {
	// 1. Check Cloudflare's specific header
	cfIP := c.GetHeader("CF-Connecting-IP")
	if cfIP != "" {
		return cfIP
	}

	// 2. Check X-Forwarded-For (can be a comma-separated list)
	xForwardedFor := c.GetHeader("X-Forwarded-For")
	if xForwardedFor != "" {
		// The first IP in the list is the original client
		parts := strings.Split(xForwardedFor, ",")
		return strings.TrimSpace(parts[0])
	}

	// 3. Fallback to the direct connection IP (RemoteAddr)
	return c.ClientIP()
}
