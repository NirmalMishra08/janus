package middleware

import (
	"log"
	"net/http"
	"time"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap the response writer to capture status code
		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// Call the next handler
		next.ServeHTTP(rw, r)

		// Calculate latency
		latency := time.Since(start)

		// Log the request
		log.Printf("[%s] %s %s %d %s",
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
			rw.statusCode,
			latency,
		)
	})
}

// responseWriter helps us capture the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}