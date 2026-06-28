package app

import (
	"fmt"

	"github.com/mangalores/case-studies-voiceline/internal/message"
)

type RecordingStore interface {
	StoreRecording(file []byte) (string, error)
}

type UploadService struct {
	recordings          RecordingStorer
	transcribePublisher TranscribePublisher
}

func NewUploadService(recordings RecordingStorer, transcribePublisher TranscribePublisher) *UploadService {
	return &UploadService{recordings: recordings, transcribePublisher: transcribePublisher}
}

func (s *UploadService) UploadRecording(file []byte) (string, error) {
	id, err := s.recordings.StoreRecording(file)
	if err != nil {
		return "", err
	}

	if s.transcribePublisher != nil {
		if err := s.transcribePublisher.Publish(message.TranscribeCommand{ID: id}); err != nil {
			return "", fmt.Errorf("publish transcribe command: %w", err)
		}
	}

	return id, nil
}
