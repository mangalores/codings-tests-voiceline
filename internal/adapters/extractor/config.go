package extractor

import "github.com/kelseyhightower/envconfig"

type MockConfig struct {
	ExtractionPath string `envconfig:"EXTRACTION_PATH" default:"assets/mock_extraction.json"`
}

type OpenAIConfig struct {
	APIKey string `envconfig:"API_KEY"`
}

func LoadMockConfig() (MockConfig, error) {
	var cfg MockConfig
	if err := envconfig.Process("MOCK", &cfg); err != nil {
		return MockConfig{}, err
	}

	return cfg, nil
}

func LoadOpenAIConfig() (OpenAIConfig, error) {
	var cfg OpenAIConfig
	if err := envconfig.Process("OPENAI", &cfg); err != nil {
		return OpenAIConfig{}, err
	}

	return cfg, nil
}
