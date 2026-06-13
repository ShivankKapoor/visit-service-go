package middleware

import (
	"log/slog"
	"net/http"
	"time"
	"visit-service/internal/service"

	gocache "github.com/patrickmn/go-cache"
	"golang.org/x/time/rate"
)

var cache = gocache.New(5*time.Minute, 10*time.Minute)

func getLimiter(ip string) *rate.Limiter {
	if l, found := cache.Get(ip); found {
		return l.(*rate.Limiter)
	}
	l := rate.NewLimiter(rate.Every(time.Minute/20), 20)
	cache.Set(ip, l, gocache.DefaultExpiration)
	return l
}

func RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := service.GetClientIP(r)
		if !getLimiter(ip).Allow() {
			slog.Warn("Rate limit exceeded", "ip", ip)
			http.Error(w, "too many requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
