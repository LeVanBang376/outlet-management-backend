package middleware

import (
	"net/http"
	"strings"
)

func CORS(allowedOrigins []string) func(http.Handler) http.Handler {

	allowed := make(map[string]bool)

	for _, origin := range allowedOrigins {
		allowed[strings.TrimSpace(origin)] = true
	}

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			origin := r.Header.Get("Origin")

			if allowed[origin] {
				w.Header().Set(
					"Access-Control-Allow-Origin",
					origin,
				)
			}

			w.Header().Set(
				"Access-Control-Allow-Methods",
				"GET, POST, PUT, PATCH, DELETE, OPTIONS",
			)

			w.Header().Set(
				"Access-Control-Allow-Headers",
				"Content-Type, Authorization",
			)

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
