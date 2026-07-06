// internal/middleware/ratelimit_redis.go
package middleware

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

func RateLimitRedis(rdb *redis.Client, requestsPerMinute int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.Background()
			key := "ratelimit:" + r.RemoteAddr

			// Simple sliding window / token bucket logic using Redis
			now := time.Now().Unix()

			pipe := rdb.Pipeline()
			// Remove old entries outside the window
			pipe.ZRemRangeByScore(ctx, key, "0", strconv.FormatInt(now-60, 10))
			// Add current request
			pipe.ZAdd(ctx, key, redis.Z{
				Score:  float64(now),
				Member: now,
			})
			// Count requests in last 60 seconds
			pipe.ZCard(ctx, key)

			cmds, err := pipe.Exec(ctx)
			if err != nil {
				// If Redis fails, allow request (fail open)
				next.ServeHTTP(w, r)
				return
			}

			count := cmds[2].(*redis.IntCmd).Val()

			if count > int64(requestsPerMinute) {
				http.Error(w, "Too many requests", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
