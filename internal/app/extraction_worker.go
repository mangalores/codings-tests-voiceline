package app

import (
	"context"
	"encoding/json"

	"log/slog"

	"github.com/mangalores/case-studies-voiceline/internal/message"
)

type ExtractionService struct {
	recordings      RecordingExtractionStore
	extractor       Extractor
	exportPublisher ExportPublisher
}

func NewExtractionService(
	recordings RecordingExtractionStore,
	extractor Extractor,
	exportPublisher ExportPublisher,
) *ExtractionService {
	return &ExtractionService{
		recordings:      recordings,
		extractor:       extractor,
		exportPublisher: exportPublisher,
	}
}

func (w *ExtractionService) Handle(ctx context.Context, command message.ExtractCommand) error {
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

	if w.exportPublisher != nil {
		if err := w.exportPublisher.Publish(message.ExportCommand{ID: command.ID}); err != nil {
			return err
		}
	}

	slog.Info("extract recording", "worker", "ExtractionWorker", "id", command.ID, "status", "success")

	return nil
}

type ExtractionWorker struct {
	commands <-chan message.ExtractCommand
	service  *ExtractionService
}

func NewExtractionWorker(
	commands <-chan message.ExtractCommand,
	service *ExtractionService,
) *ExtractionWorker {
	return &ExtractionWorker{
		commands: commands,
		service:  service,
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

			if err := w.service.Handle(ctx, command); err != nil {
				slog.Error("extract recording", "worker", "ExtractionWorker", "id", command.ID, "error", err)
				if saveErr := w.service.recordings.SaveError(command.ID, "ExtractionWorker", err.Error()); saveErr != nil {
					slog.Error("persist worker error", "worker", "ExtractionWorker", "id", command.ID, "error", saveErr)
				}
			}
		}
	}
}
