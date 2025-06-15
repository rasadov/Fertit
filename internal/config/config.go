package config

import (
	"log"
	"os"
	"sync"
)

type Settings struct {
	Debug                  bool
	TokenExpirationSeconds int
	PostgresUrl            string
}

var (
	AppConfig *Settings
	once      sync.Once
)

func Init() error {
	var err error
	once.Do(func() {
		AppConfig, err = GetSettings()
	})
	return err
}

func GetSettings() (*Settings, error) {
	debug := os.Getenv("DEBUG") == "true"

	PostgresUrl := os.Getenv("POSTGRES_URL")

	if PostgresUrl == "" {
		log.Fatal("POSTGRES_URL environment variable not set")
	}

	return &Settings{
		Debug:                  debug,
		TokenExpirationSeconds: 3600,
		PostgresUrl:            PostgresUrl,
	}, nil
}
