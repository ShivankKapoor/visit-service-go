package middleware

import (
	"net/http"
	"os"
	"strings"
)

func AllowedReferer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if os.Getenv("PROD") != "true" {
			next.ServeHTTP(w, r)
			return
		}

		referer := r.Header.Get("Referer")
		if referer == "" {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		for _, origin := range strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",") {
			if strings.Contains(referer, strings.TrimSpace(origin)) {
				next.ServeHTTP(w, r)
				return
			}
		}

		http.Error(w, "forbidden", http.StatusForbidden)
	})
}
