package exporter

import "github.com/kelseyhightower/envconfig"

type WebhookConfig struct {
	ExportURL string `envconfig:"EXPORT_URL"`
}

type GoogleConfig struct {
	CredentialsFilePath string `envconfig:"CREDENTIALS_FILE"`
	SheetID             string `envconfig:"SHEET_ID"`
	SheetRange          string `envconfig:"SHEET_RANGE" default:"meetings!A:F"`
}

func LoadWebhookConfig() (WebhookConfig, error) {
	var cfg WebhookConfig
	if err := envconfig.Process("WEBHOOK", &cfg); err != nil {
		return WebhookConfig{}, err
	}

	return cfg, nil
}

func LoadGoogleConfig() (GoogleConfig, error) {
	var cfg GoogleConfig
	if err := envconfig.Process("GOOGLE", &cfg); err != nil {
		return GoogleConfig{}, err
	}

	return cfg, nil
}
