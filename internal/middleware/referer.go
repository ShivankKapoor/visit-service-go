package middleware

import (
	"net/http"
	"net/url"
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

		refURL, err := url.Parse(referer)
		if err != nil {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		for _, origin := range strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",") {
			allowedURL, err := url.Parse(strings.TrimSpace(origin))
			if err != nil {
				continue
			}
			if refURL.Host == allowedURL.Host {
				next.ServeHTTP(w, r)
				return
			}
		}

		http.Error(w, "forbidden", http.StatusForbidden)
	})
}
