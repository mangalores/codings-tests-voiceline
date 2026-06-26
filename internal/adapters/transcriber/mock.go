package transcriber

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
)

type MockTranscriber struct {
	transcriptionPath string
}

func NewMockTranscriber(cfg MockConfig) (*MockTranscriber, error) {
	transcriptionPath := strings.TrimSpace(cfg.TranscriptionPath)
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

	slog.Info("MockTranscriber: transcribed audio")

	return string(content), nil
}
