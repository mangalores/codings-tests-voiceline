package exporter

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/mangalores/case-studies-voiceline/internal/app"
)

type WebhookExporter struct {
	webhookURL string
}

func NewWebhookExporter(cfg WebhookConfig) *WebhookExporter {
	return &WebhookExporter{webhookURL: cfg.ExportURL}
}

func (w *WebhookExporter) Export(ctx context.Context, data app.ExtractedData) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, w.webhookURL, strings.NewReader(string(payload)))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("webhook export error: %s", strings.TrimSpace(string(respBody)))
	}

	slog.Info("WebhookExporter exported data", "status", resp.StatusCode)

	return nil
}
