package app

import (
	"context"

	"log/slog"
)

type ExportWorker struct {
	commands   <-chan ExportCommand
	recordings RecordingExportStore
	exporter   Exporter
}

func NewExportWorker(
	recordings RecordingExportStore,
	exporter Exporter,
	commands <-chan ExportCommand,
) *ExportWorker {
	return &ExportWorker{
		recordings: recordings,
		exporter:   exporter,
		commands:   commands,
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

			if err := w.handleCommand(ctx, command); err != nil {
				slog.Error("export recording", "worker", "ExportWorker", "id", command.ID, "error", err)
				if saveErr := w.recordings.SaveError(command.ID, "ExportWorker", err.Error()); saveErr != nil {
					slog.Error("persist worker error", "worker", "ExportWorker", "id", command.ID, "error", saveErr)
				}
			}
		}
	}
}

func (w *ExportWorker) handleCommand(ctx context.Context, command ExportCommand) error {
	data, err := w.recordings.GetExtraction(command.ID)
	if err != nil {
		return err
	}

	if err := w.exporter.Export(ctx, data); err != nil {
		return err
	}

	slog.Info("export recording", "worker", "ExportWorker", "id", command.ID, "status", "success")

	return nil
}
