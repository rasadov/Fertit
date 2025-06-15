package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
)

var Client *redis.Client

func InitRedis(ctx context.Context) {
	Client = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0, // default DB
	})

	// Test connection
	_, err := Client.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}

	log.Println("Connected to Redis successfully")
}
