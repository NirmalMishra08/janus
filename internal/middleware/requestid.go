package middleware

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

func RequestId(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get("X-Request-ID")

		if reqID == ""{
			reqID = uuid.New().String()
		}

		ctx := context.WithValue(r.Context(), middleware.RequestIDKey, reqID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
