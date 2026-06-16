package store

import (
	"context"
	"log/slog"
	"os"

	"github.com/redis/go-redis/v9"
)

func NewRedisConnection() *redis.Client {

	redisAddr:=os.Getenv("REDIS_ADDR")

	if redisAddr==""{
		redisAddr="localhost:6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "", // no password
		DB:       0,  // use default DB
		Protocol: 2,
	})

	ctx := context.Background()

	pong, err := rdb.Ping(ctx).Result()

	if err != nil {
		slog.Warn("Redis unavailable at startup, running in fallback mode")

	}

	slog.Info("Connected to Redis", "pong", pong)

	return rdb

}
