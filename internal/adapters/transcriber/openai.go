package transcriber

import (
	"context"
	"io"
)

type OpenAITranscriber struct {
	apiKey string
}

func NewOpenAITranscriber(apiKey string) *OpenAITranscriber {

	return &OpenAITranscriber{apiKey: apiKey}
}

func (o *OpenAITranscriber) Transcribe(ctx context.Context, audio io.Reader) (string, error) {
	// Placeholder for actual OpenAI API call
	return "", nil
}
