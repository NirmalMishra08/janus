package redis

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

func InitRedis(redisUrl string) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisUrl,
		Password: "", // no password
		DB:       0,  // use default DB
		Protocol: 2,
	})

	ctx := context.Background()

	err := rdb.Ping(ctx).Err()
	if err != nil {
		log.Fatalf("failed to connect to Redis: %v", err)
	}

	log.Println("conneced to redis")

	return rdb
}
