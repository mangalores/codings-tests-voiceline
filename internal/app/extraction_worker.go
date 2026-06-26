package app

import (
	"context"
	"encoding/json"

	"log/slog"
)

type ExtractionWorker struct {
	commands     <-chan ExtractCommand
	recordings   RecordingExtractionStore
	extractor    Extractor
	nextCommands chan<- ExportCommand
}

func NewExtractionWorker(
	recordings RecordingExtractionStore,
	extractor Extractor,
	commands <-chan ExtractCommand,
	nextCommands chan<- ExportCommand,
) *ExtractionWorker {
	return &ExtractionWorker{
		recordings:   recordings,
		extractor:    extractor,
		commands:     commands,
		nextCommands: nextCommands,
	}
}

func (w *ExtractionWorker) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case command, ok := <-w.commands:
			if !ok {
				return
			}

			if err := w.handleCommand(ctx, command); err != nil {
				slog.Error("extract recording", "worker", "ExtractionWorker", "id", command.ID, "error", err)
				if saveErr := w.recordings.SaveError(command.ID, "ExtractionWorker", err.Error()); saveErr != nil {
					slog.Error("persist worker error", "worker", "ExtractionWorker", "id", command.ID, "error", saveErr)
				}
			}
		}
	}
}

func (w *ExtractionWorker) handleCommand(ctx context.Context, command ExtractCommand) error {
	transcription, err := w.recordings.GetTranscription(command.ID)
	if err != nil {
		return err
	}

	extracted, err := w.extractor.Extract(ctx, transcription)
	if err != nil {
		return err
	}

	payload, err := json.Marshal(extracted)
	if err != nil {
		return err
	}

	if err := w.recordings.SaveExtraction(command.ID, payload); err != nil {
		return err
	}

	if w.nextCommands != nil {
		w.nextCommands <- ExportCommand{ID: command.ID}
	}

	slog.Info("extract recording", "worker", "ExtractionWorker", "id", command.ID, "status", "success")

	return nil
}
