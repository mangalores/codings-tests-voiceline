package transcriber

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"strings"
)

type OpenAITranscriber struct {
	apiKey string
}

func NewOpenAITranscriber(cfg OpenAIConfig) *OpenAITranscriber {
	return &OpenAITranscriber{apiKey: cfg.APIKey}
}

func (o *OpenAITranscriber) Transcribe(ctx context.Context, audio io.Reader) (string, error) {
	if o.apiKey == "" {
		return "", fmt.Errorf("openai api key is required")
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormFile("file", "recording.mp3")
	if err != nil {
		return "", err
	}

	if _, err := io.Copy(part, audio); err != nil {
		return "", err
	}

	if err := writer.WriteField("model", "gpt-4o-mini-transcribe"); err != nil {
		return "", err
	}

	if err := writer.WriteField("response_format", "json"); err != nil {
		return "", err
	}

	if err := writer.Close(); err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://api.openai.com/v1/audio/transcriptions",
		&body,
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+o.apiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode >= http.StatusMultipleChoices {
		return "", fmt.Errorf("openai transcription error: %s", strings.TrimSpace(string(respBody)))
	}

	var result struct {
		Text string `json:"text"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", err
	}

	slog.Info("OpenAITranscriber: transcribed audio")

	return result.Text, nil
}
