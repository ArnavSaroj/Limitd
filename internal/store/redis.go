package store

import (
	"context"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

func NewRedisConnection() *redis.Client {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
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
