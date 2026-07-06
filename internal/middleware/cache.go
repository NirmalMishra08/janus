// middleware/cache.go
package middleware

// import (
// 	"net/http"
// 	"time"

// 	"github.com/go-chi/chi/v5/middleware"
// 	"github.com/redis/go-redis/v9"
// )

// func ResponseCache(rdb *redis.Client) func(http.Handler) http.Handler {
// 	return middleware.ResponseCache(rdb, middleware.ResponseCacheOptions{
// 		TTL:                  5 * time.Minute,
// 		CacheableMethods:     []string{"GET", "HEAD"},
// 		CacheableStatusCodes: []int{200, 204, 301, 304},
// 	})
// }