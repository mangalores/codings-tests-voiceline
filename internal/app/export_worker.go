package app

import (
	"context"

	"log/slog"

	"github.com/mangalores/case-studies-voiceline/internal/message"
)

type ExportService struct {
	recordings RecordingExportStore
	exporter   Exporter
}

func NewExportService(
	recordings RecordingExportStore,
	exporter Exporter,
) *ExportService {
	return &ExportService{
		recordings: recordings,
		exporter:   exporter,
	}
}

func (w *ExportService) Handle(ctx context.Context, command message.ExportCommand) error {
	data, err := w.recordings.GetExtraction(command.ID)
	if err != nil {
		return err
	}

	if err := w.exporter.Export(ctx, data); err != nil {
		return err
	}

	slog.Info("SUCCESS export recording", "worker", "ExportWorker", "id", command.ID, "status", "success")

	return nil
}

type ExportWorker struct {
	commands <-chan message.ExportCommand
	service  *ExportService
}

func NewExportWorker(
	commands <-chan message.ExportCommand,
	service *ExportService,
) *ExportWorker {
	return &ExportWorker{
		commands: commands,
		service:  service,
	}
}

func (w *ExportWorker) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case command, ok := <-w.commands:
			if !ok {
				return
			}

			if err := w.service.Handle(ctx, command); err != nil {
				slog.Error("FAILED export recording", "worker", "ExportWorker", "id", command.ID, "error", err)
				if saveErr := w.service.recordings.SaveError(command.ID, "ExportWorker", err.Error()); saveErr != nil {
					slog.Error("FAILED persist worker error", "worker", "ExportWorker", "id", command.ID, "error", saveErr)
				}
			}
		}
	}
}
