package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`

	jwt.RegisteredClaims
}

const userContextKey string = "user"

func JWTAuth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if parts[0] != "Bearer" || len(parts) != 2 {
				http.Error(w, "Invalid Token", http.StatusUnauthorized)
				return
			}

			authToken := parts[1]

			token, err := jwt.ParseWithClaims(authToken, &Claims{}, func(t *jwt.Token) (any, error) {
				return []byte(secret), nil
			})
			if err != nil || !token.Valid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			claims := token.Claims.(*Claims)

			ctx := context.WithValue(r.Context(), userContextKey, claims)

			next.ServeHTTP(w, r.WithContext(ctx))

		})
	}
}

func RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(userContextKey).(*Claims)
			if !ok || claims.Role != role {
				http.Error(w, "Unauthorized", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
