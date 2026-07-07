package middleware

import (
	"net/http"
	"server/internal/logger"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/zap"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap the response writer to capture status code
		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		reqID := middleware.GetReqID(r.Context())

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		// Call the next handler
		next.ServeHTTP(rw, r)

		duration := time.Since(start)

		logger.Log.Info("HTTP Request",
			zap.String("request_id", reqID),
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Int("status", ww.Status()),
			zap.Duration("duration", duration),
			zap.String("remote_ip", r.RemoteAddr),
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
func CORS() func(http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // Change in production
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	})
}

func Compress() func(http.Handler) http.Handler {
	return middleware.Compress(5) // gzip level 5 is good balance
}
