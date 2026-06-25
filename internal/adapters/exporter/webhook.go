package exporter

import (
	"context"

	"github.com/mangalores/case-studies-voiceline/internal/app"
)

type WebhookExporter struct {
	webhookURL string
}

func NewWebhookExporter(webhookURL string) *WebhookExporter {
	return &WebhookExporter{webhookURL: webhookURL}
}

func (w *WebhookExporter) Export(ctx context.Context, data app.ExtractedData) error {
	// Implement the logic to send data to the webhook URL.
	// This is a placeholder implementation. You would typically use an HTTP client to POST the data to the webhook.
	return nil
}
