package service

import (
	"log/slog"
	"net"
	"net/http"
	"strings"
)

func GetClientIP(r *http.Request) string {
	cfIP := r.Header.Get("CF-Connecting-IP")
	if cfIP != "" {
		return cfIP
	}

	slog.Warn("No Cloudflare IP found falling back to X-Forwarded-For")
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	if xForwardedFor != "" {
		parts := strings.Split(xForwardedFor, ",")
		return strings.TrimSpace(parts[0])
	}

	slog.Warn("No X-Forwarded-For found falling back to regular IP header")
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}
