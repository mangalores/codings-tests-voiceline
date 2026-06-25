package extractor

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/mangalores/case-studies-voiceline/internal/app"
)

type MockExtractor struct {
	extractionPath string
}

func NewMockExtractor(extractionPath string) *MockExtractor {
	return &MockExtractor{extractionPath: extractionPath}
}

func (m *MockExtractor) Extract(ctx context.Context, transcription string) (app.ExtractedData, error) {
	content, err := os.ReadFile(m.extractionPath)
	if err != nil {
		return app.ExtractedData{}, fmt.Errorf("read mock extraction: %w", err)
	}

	var data app.ExtractedData
	if err := json.Unmarshal(content, &data); err != nil {
		return app.ExtractedData{}, fmt.Errorf("decode mock extraction: %w", err)
	}

	return data, nil
}
