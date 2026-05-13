package middlewares

import (
	"log"
	"net/http"
	"time"
)

func LoggerMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			duration := time.Since(start)
			log.Printf("[API HIT] %s | %s | %v", r.Method, r.URL.Path, duration)
		})
	}
}
