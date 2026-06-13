package middleware

import (
	"net/http"
	"os"
	"strings"
)

func Cors(next http.Handler) http.Handler {
	prod := os.Getenv("PROD") == "true"

	var allowedOrigins map[string]bool
	if prod {
		allowedOrigins = make(map[string]bool)
		for _, o := range strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",") {
			allowedOrigins[strings.TrimSpace(o)] = true
		}
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if prod {
			if origin != "" && allowedOrigins[origin] {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
