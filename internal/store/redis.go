package store

import (
	"context"
	"fmt"
	"log"

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
		log.Fatal("Error connecting to Redis:", err)

	}

	fmt.Println("Connected to Redis:", pong)

	return rdb

}
