package redis

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

func InitRedis(addr string, db int) {
	Client = redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   db,
	})

	ctx := context.Background()
	_, err := Client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("failed to connect to Redis %v", err)
	}
}
