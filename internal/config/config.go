package config

import (
	"fmt"
	"log"
	"os"
	"sync"
)

type Settings struct {
	Debug                  bool
	TokenExpirationSeconds int
	PostgresUrl            string
	RedisAddr              string
	RedisPassword          string
}

var (
	AppConfig *Settings
	once      sync.Once
)

func init() {
	var err error
	once.Do(func() {
		AppConfig, err = GetSettings()
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("AppConfig:", AppConfig)
}

func GetSettings() (*Settings, error) {
	debug := os.Getenv("DEBUG") == "true"

	postgresHost := os.Getenv("POSTGRES_HOST")
	postgresPort := os.Getenv("POSTGRES_PORT")
	postgresUser := os.Getenv("POSTGRES_USER")
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	postgresDB := os.Getenv("POSTGRES_DB")

	// Default values if not provided
	if postgresHost == "" {
		postgresHost = "localhost"
	}
	if postgresPort == "" {
		postgresPort = "5432"
	}
	if postgresDB == "" {
		postgresDB = "postgres"
	}

	// Construct PostgreSQL URL
	postgresUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		postgresUser, postgresPassword, postgresHost, postgresPort, postgresDB)

	log.Printf("postgresUrl: %s", postgresUrl)

	if postgresUser == "" || postgresPassword == "" {
		log.Fatal("POSTGRES_USER and POSTGRES_PASSWORD environment variables must be set")
	}

	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	if redisHost == "" || redisPort == "" {
		log.Fatal("REDIS_HOST and REDIS_PASSWORD environment variables must be set")
	}

	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)

	return &Settings{
		Debug:                  debug,
		TokenExpirationSeconds: 3600,
		PostgresUrl:            postgresUrl,
		RedisAddr:              redisAddr,
		RedisPassword:          redisPassword,
	}, nil
}
