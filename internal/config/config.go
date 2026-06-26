package config

import (
	"fmt"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type AppConfig struct {
	Port string `envconfig:"PORT" default:"8080"`

	Transcriber string `envconfig:"TRANSCRIBER" default:"mock"`
	Extractor   string `envconfig:"EXTRACTOR" default:"mock"`
	Exporter    string `envconfig:"EXPORTER" default:"mock"`

	MockConfig
	FileStoreConfig
	OpenAIConfig
	WebHookConfig
	GoogleConfig
}

// todo: the implementation specific configs should be moved to their respective packages, and the config should be passed to the build functions of the adapters.
type MockConfig struct {
	TranscriptionPath string `envconfig:"TRANSCRIPTION_PATH" default:"assets/mock_transcription.md"`
	ExtractionPath    string `envconfig:"EXTRACTION_PATH" default:"assets/mock_extraction.json"`
}

type FileStoreConfig struct {
	StoragePath string `envconfig:"STORAGE_PATH" default:"./storage"`
}

type OpenAIConfig struct {
	APIKey string `envconfig:"API_KEY"`
}

type WebHookConfig struct {
	ExportURL string `envconfig:"EXPORT_URL"`
}

type GoogleConfig struct {
	CredentialsFilePath string `envconfig:"CREDENTIALS_FILE"`
	SheetID             string `envconfig:"SHEET_ID"`
	SheetRange          string `envconfig:"SHEET_RANGE" default:"meetings!A:F"`
}

func Load() (*AppConfig, error) {
	_ = godotenv.Load()

	cfg := &AppConfig{}

	if err := envconfig.Process("", cfg); err != nil {
		return nil, err
	}
	if err := envconfig.Process("MOCK", &cfg.MockConfig); err != nil {
		return nil, err
	}

	if err := envconfig.Process("FILE", &cfg.FileStoreConfig); err != nil {
		return nil, err
	}

	if err := envconfig.Process("OPENAI", &cfg.OpenAIConfig); err != nil {
		return nil, err
	}

	if err := envconfig.Process("WEBHOOK", &cfg.WebHookConfig); err != nil {
		return nil, err
	}

	if err := envconfig.Process("GOOGLE", &cfg.GoogleConfig); err != nil {
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
