package redis

import (
	"context"
	"github.com/rasadov/MailManagerApp/internal/config"
	"github.com/redis/go-redis/v9"
	"log"
)

func GetRedisClient(ctx context.Context, redisAddr, redisPassword string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		Password: config.RedisPassword,
		DB:       0, // default DB
	})

	// Test connection
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}

	log.Println("Connected to Redis successfully")
	return client
}
