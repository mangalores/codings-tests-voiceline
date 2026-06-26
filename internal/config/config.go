package config

import (
	"fmt"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type AppConfig struct {
	Port string `envconfig:"PORT" default:"8080"`

	StoragePath string `envconfig:"FILE_STORAGE_PATH" default:"./storage"`

	Transcriber string `envconfig:"TRANSCRIBER" default:"mock"`
	Extractor   string `envconfig:"EXTRACTOR" default:"mock"`
	Exporter    string `envconfig:"EXPORTER" default:"mock"`
}

func Load() (*AppConfig, error) {
	_ = godotenv.Load()

	cfg := &AppConfig{}

	if err := envconfig.Process("", cfg); err != nil {
		return nil, err
	}

	if err := validatePort(cfg.Port); err != nil {
		return nil, err
	}

	return cfg, nil
}

func validatePort(port string) error {
	value, err := strconv.Atoi(port)
	if err != nil {
		return fmt.Errorf("port must be numeric")
	}

	if value < 1 || value > 65535 {
		return fmt.Errorf("port must be between 1 and 65535")
	}

	return nil
}
