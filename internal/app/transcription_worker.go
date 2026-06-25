package app

import (
	"context"
	"fmt"

	"log/slog"
)

type TranscriptionWorker struct {
	commands     <-chan TranscribeCommand
	nextCommands chan<- ExtractCommand
	recordings   RecordingTranscriptionStore
	transcriber  Transcriber
}

func NewTranscriptionService(
	recordings RecordingTranscriptionStore,
	transcriber Transcriber,
	commands <-chan TranscribeCommand,
	nextCommands chan<- ExtractCommand,
) *TranscriptionWorker {
	return &TranscriptionWorker{
		recordings:   recordings,
		transcriber:  transcriber,
		commands:     commands,
		nextCommands: nextCommands,
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

			if err := s.handleCommand(ctx, command); err != nil {
				slog.Error("transcribe recording", "worker", "TranscriptionWorker", "id", command.ID, "error", err)
				if saveErr := s.recordings.SaveError(command.ID, "TranscriptionWorker", err.Error()); saveErr != nil {
					slog.Error("persist worker error", "worker", "TranscriptionWorker", "id", command.ID, "error", saveErr)
				}
			}
		}
	}
}

func (s *TranscriptionWorker) handleCommand(ctx context.Context, command TranscribeCommand) error {
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

	if s.nextCommands != nil {
		s.nextCommands <- ExtractCommand{ID: command.ID}
	}

	slog.Info("transcribe recording", "worker", "TranscriptionWorker", "id", command.ID, "status", "success")

	return nil
}
