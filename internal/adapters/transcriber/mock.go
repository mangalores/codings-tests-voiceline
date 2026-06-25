package transcriber

import (
	"context"
	"fmt"
	"io"
	"os"
)

type MockTranscriber struct {
	transcriptionPath string
}

type MockConfig struct {
	TranscriptionPath string `envconfig:"TRANSCRIPTION_PATH" default:"assets/mock_transcription.md"`
}

func NewMockTranscriber(transcriptionPath string) (*MockTranscriber, error) {
	if transcriptionPath == "" {
		transcriptionPath = "assets/mock_transcription.md"
	}

	return &MockTranscriber{transcriptionPath: transcriptionPath}, nil
}

func (m *MockTranscriber) Transcribe(ctx context.Context, audio io.Reader) (string, error) {
	content, err := os.ReadFile(m.transcriptionPath)
	if err != nil {
		return "", fmt.Errorf("read mock transcription: %w", err)
	}

	return string(content), nil
}
