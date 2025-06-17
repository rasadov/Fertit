package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
)

type Settings struct {
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

	jwtSecret := os.Getenv("JWT_SECRET")
	jwtIssuer := os.Getenv("JWT_ISSUER")

	postgresHost := os.Getenv("POSTGRES_HOST")
	postgresPort := os.Getenv("POSTGRES_PORT")
	postgresUser := os.Getenv("POSTGRES_USER")
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	postgresDB := os.Getenv("POSTGRES_DB")

	adminUsername := os.Getenv("ADMIN_USERNAME")
	adminPassword := os.Getenv("ADMIN_PASSWORD")

	smtpHost := os.Getenv("SMTP_HOST")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPassword := os.Getenv("SMTP_PASSWORD")
	smtpPort := os.Getenv("SMTP_PORT")

	baseUrl := os.Getenv("BASE_URL")

	smtpPortInt, err := strconv.Atoi(smtpPort)
	if err != nil {
		return nil, err
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
	redisPassword := os.Getenv("REDIS_PASSWORD")

	if redisHost == "" || redisPort == "" {
		log.Fatal("REDIS_HOST and REDIS_PASSWORD environment variables must be set")
	}

	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)

	return &Settings{
		Debug:                  debug,
		BaseUrl:                baseUrl,
		JWTSecret:              jwtSecret,
		JWTIssuer:              jwtIssuer,
		TokenExpirationSeconds: 3600,

		PostgresUrl: postgresUrl,

		RedisAddr:     redisAddr,
		RedisPassword: redisPassword,

		AdminUsername: adminUsername,
		AdminPassword: adminPassword,

		SmtpHost:     smtpHost,
		SmtpPort:     smtpPortInt,
		SmtpUsername: smtpUser,
		SmtpPassword: smtpPassword,
	}, nil
}
