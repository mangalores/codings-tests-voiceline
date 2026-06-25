package extractor

import (
	"context"

	"github.com/mangalores/case-studies-voiceline/internal/app"
)

type OpenAIExtractor struct {
	APIKey string
}

func NewOpenAIExtractor(apiKey string) *OpenAIExtractor {
	return &OpenAIExtractor{APIKey: apiKey}
}

func (e *OpenAIExtractor) Extract(ctx context.Context, transcription string) (app.ExtractedData, error) {
	// Placeholder for actual OpenAI API call
	return app.ExtractedData{}, nil

}
