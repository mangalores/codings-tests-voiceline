package exporter

import (
	"context"

	"log/slog"

	"github.com/mangalores/case-studies-voiceline/internal/app"
)

type MockExporter struct{}

func NewMockExporter() *MockExporter {
	return &MockExporter{}
}

func (m *MockExporter) Export(ctx context.Context, data app.ExtractedData) error {
	slog.Info("MockExporter exported data", "data", data)
	return nil
}
