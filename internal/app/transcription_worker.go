package app

import (
	"context"
	"fmt"

	"log/slog"

	"github.com/mangalores/case-studies-voiceline/internal/message"
)

type TranscriptionService struct {
	recordings       RecordingTranscriptionStore
	transcriber      Transcriber
	extractPublisher ExtractPublisher
}

func NewTranscriptionService(
	recordings RecordingTranscriptionStore,
	transcriber Transcriber,
	extractPublisher ExtractPublisher,
) *TranscriptionService {
	return &TranscriptionService{
		recordings:       recordings,
		transcriber:      transcriber,
		extractPublisher: extractPublisher,
	}
}

func (s *TranscriptionService) Handle(ctx context.Context, command message.TranscribeCommand) error {
	audio, err := s.recordings.OpenRecording(command.ID)
	if err != nil {
		return err
	}
	defer audio.Close()

	transcription, err := s.transcriber.Transcribe(ctx, audio)
	if err != nil {
		return err
	}

	if err := s.recordings.SaveTranscription(command.ID, []byte(transcription)); err != nil {
		return fmt.Errorf("save transcription: %w", err)
	}

	if s.extractPublisher != nil {
		if err := s.extractPublisher.Publish(message.ExtractCommand{ID: command.ID}); err != nil {
			return fmt.Errorf("publish extract command: %w", err)
		}
	}

	slog.Info("transcribe recording", "worker", "TranscriptionWorker", "id", command.ID, "status", "success")

	return nil
}

type TranscriptionWorker struct {
	commands <-chan message.TranscribeCommand
	service  *TranscriptionService
}

func NewTranscriptionWorker(
	commands <-chan message.TranscribeCommand,
	service *TranscriptionService,
) *TranscriptionWorker {
	return &TranscriptionWorker{
		commands: commands,
		service:  service,
	}
}

func (s *TranscriptionWorker) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case command, ok := <-s.commands:
			if !ok {
				return
			}

			if err := s.service.Handle(ctx, command); err != nil {
				slog.Error("transcribe recording", "worker", "TranscriptionWorker", "id", command.ID, "error", err)
				if saveErr := s.service.recordings.SaveError(command.ID, "TranscriptionWorker", err.Error()); saveErr != nil {
					slog.Error("persist worker error", "worker", "TranscriptionWorker", "id", command.ID, "error", saveErr)
				}
			}
		}
	}
}
