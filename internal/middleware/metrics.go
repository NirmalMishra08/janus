package middleware

import (
	"net/http"
	"server/internal/metrics"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func PrometheusMetrics(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		metrics.HttpActiveRequests.Inc()

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		next.ServeHTTP(ww , r)

		duration := time.Since(start).Seconds()

		metrics.HttpRequestDuration.WithLabelValues(
			r.Method,
			r.URL.Path,
		).Observe(duration)

		metrics.HttpRequestsTotal.WithLabelValues(
			r.Method,
			r.URL.Path,
			strconv.Itoa(ww.Status()),
		).Inc()

		metrics.HttpActiveRequests.Dec()
	})
}