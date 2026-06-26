package exporter

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/mangalores/case-studies-voiceline/internal/app"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type GoogleSheetExporter struct {
	credentialsFilePath string
	sheetID             string
	sheetRange          string
}

func NewGoogleSheetExporter(credentialsFilePath, sheetID, sheetRange string) *GoogleSheetExporter {
	return &GoogleSheetExporter{
		credentialsFilePath: credentialsFilePath,
		sheetID:             sheetID,
		sheetRange:          sheetRange,
	}
}

func (g *GoogleSheetExporter) Export(ctx context.Context, data app.ExtractedData) error {
	if err := g.validate(); err != nil {
		return err
	}

	srv, err := sheets.NewService(ctx, option.WithCredentialsFile(g.credentialsFilePath))
	if err != nil {
		return fmt.Errorf("create google sheets service: %w", err)
	}

	rows := buildSheetRows(data)
	_, err = srv.Spreadsheets.Values.Append(g.sheetID, g.sheetRange, &sheets.ValueRange{
		Values: rows,
	}).ValueInputOption("RAW").InsertDataOption("INSERT_ROWS").Do()
	if err != nil {
		return fmt.Errorf("append google sheets rows: %w", err)
	}

	slog.Info("GoogleSheetExporter exported data", "rows", len(rows))

	return nil
}

func (g *GoogleSheetExporter) validate() error {
	if strings.TrimSpace(g.credentialsFilePath) == "" {
		return fmt.Errorf("google sheets credentials file path is required")
	}

	if strings.TrimSpace(g.sheetID) == "" {
		return fmt.Errorf("google sheet id is required")
	}

	if strings.TrimSpace(g.sheetRange) == "" {
		return fmt.Errorf("google sheet range is required")
	}

	return nil
}

func buildSheetRows(data app.ExtractedData) [][]interface{} {
	rows := [][]interface{}{
		{"Summary", data.Summary, "", "", "", ""},
	}

	rows = appendStringSection(rows, "Participants", data.Participants)
	rows = appendStringSection(rows, "Decisions", data.Decisions)

	rows = append(rows, []interface{}{"Action Items", "Owner", "Task", "Due", "", ""})
	for _, item := range data.ActionItems {
		rows = append(rows, []interface{}{"", item.Owner, item.Task, item.Due, "", ""})
	}

	return rows
}

func appendStringSection(rows [][]interface{}, label string, values []string) [][]interface{} {
	if len(values) == 0 {
		return rows
	}

	rows = append(rows, []interface{}{label, values[0], "", "", "", ""})
	for _, value := range values[1:] {
		rows = append(rows, []interface{}{"", value, "", "", "", ""})
	}

	return rows
}
