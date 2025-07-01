package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
)

var (
	Debug                  bool
	BaseUrl                string
	JWTSecret              string
	JWTIssuer              string
	TokenExpirationSeconds int
	PostgresUrl            string
	RedisAddr              string
	RedisPassword          string
	AdminUsername          string
	AdminPassword          string
	SmtpHost               string
	SmtpPort               int
	SmtpUsername           string
	SmtpPassword           string
)

var (
	once sync.Once
)

func init() {
	var err error
	once.Do(func() {
		err = GetSettings()
	})
	if err != nil {
		log.Fatal(err)
	}
}

func GetSettings() error {
	var err error
	Debug = os.Getenv("DEBUG") == "true"

	JWTSecret = os.Getenv("JWT_SECRET")
	JWTIssuer = os.Getenv("JWT_ISSUER")

	postgresHost := os.Getenv("POSTGRES_HOST")
	postgresPort := os.Getenv("POSTGRES_PORT")
	postgresUser := os.Getenv("POSTGRES_USER")
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	postgresDB := os.Getenv("POSTGRES_DB")

	AdminUsername = os.Getenv("ADMIN_USERNAME")
	AdminPassword = os.Getenv("ADMIN_PASSWORD")

	SmtpHost = os.Getenv("SMTP_HOST")
	SmtpUsername = os.Getenv("SMTP_USER")
	SmtpPassword = os.Getenv("SMTP_PASSWORD")
	smtpPortStr := os.Getenv("SMTP_PORT")

	BaseUrl = os.Getenv("BASE_URL")

	SmtpPort, err = strconv.Atoi(smtpPortStr)
	if err != nil {
		return err
	}

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
	RedisPassword = os.Getenv("REDIS_PASSWORD")

	if redisHost == "" || redisPort == "" {
		log.Fatal("REDIS_HOST and REDIS_PASSWORD environment variables must be set")
	}

	RedisAddr = fmt.Sprintf("%s:%s", redisHost, redisPort)

	return nil
}
