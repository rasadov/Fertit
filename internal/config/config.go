package config

import (
	"os"
	"sync"
)

type Settings struct {
	Debug                  bool
	TokenExpirationSeconds int
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

	return &Settings{
		Debug:                  debug,
		TokenExpirationSeconds: 3600,
	}, nil
}
