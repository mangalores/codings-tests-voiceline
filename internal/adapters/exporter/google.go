package exporter

import (
	"context"
	"log/slog"

	"github.com/mangalores/case-studies-voiceline/internal/app"
)

type GoogleSheetExporter struct {
	CredentialsFilePath string
	SheetID             string
	SheetRange          string
}

func NewGoogleSheetExporter(credentialsFilePath, sheetID, sheetRange string) *GoogleSheetExporter {
	return &GoogleSheetExporter{
		CredentialsFilePath: credentialsFilePath,
		SheetID:             sheetID,
		SheetRange:          sheetRange,
	}
}

func (g *GoogleSheetExporter) Export(ctx context.Context, data app.ExtractedData) error {
	// Implement the logic to export data to Google Sheets using the provided credentials and sheet information.
	// This is a placeholder implementation. You would typically use the Google Sheets API here.
	slog.Info("GoogleSheetExporter exported data")

	return nil
}
