package redis

import (
	"context"
	"github.com/rasadov/MailManagerApp/internal/config"
	"github.com/redis/go-redis/v9"
	"log"
)

var Client *redis.Client

func InitRedis(ctx context.Context) {
	Client = redis.NewClient(&redis.Options{
		Addr:     config.AppConfig.RedisAddr,
		Password: config.AppConfig.RedisPassword,
		DB:       0, // default DB
	})

	// Test connection
	_, err := Client.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}

	log.Println("Connected to Redis successfully")
}
