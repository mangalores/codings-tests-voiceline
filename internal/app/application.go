package app

import (
	"context"
	"errors"
	"net/http"
	"time"

	"log/slog"
)

type Application struct {
	server              *http.Server
	transcriptionWorker *TranscriptionWorker
	extractionWorker    *ExtractionWorker
	exportWorker        *ExportWorker
}

func NewApplication(
	server *http.Server,
	transcriptionWorker *TranscriptionWorker,
	extractionWorker *ExtractionWorker,
	exportWorker *ExportWorker,
) *Application {
	return &Application{
		server:              server,
		transcriptionWorker: transcriptionWorker,
		extractionWorker:    extractionWorker,
		exportWorker:        exportWorker,
	}
}

func (a *Application) Run(ctx context.Context) error {
	slog.Info("start application", "worker", "Application")

	serverErrCh := make(chan error, 1)
	serverDone := make(chan struct{})

	go func() {
		defer close(serverDone)
		err := a.server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("start server", "worker", "Application", "error", err)
			serverErrCh <- err
			return
		}

		serverErrCh <- nil
	}()

	go a.transcriptionWorker.Run(ctx)
	go a.extractionWorker.Run(ctx)
	go a.exportWorker.Run(ctx)

	go func() {
		<-ctx.Done()
		slog.Info("shutdown application", "worker", "Application")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := a.server.Shutdown(shutdownCtx); err != nil {
			slog.Error("shutdown server", "error", err)
		}
	}()

	select {
	case err := <-serverErrCh:
		<-serverDone
		return err
	case <-ctx.Done():
		<-serverDone
		return nil
	}
}

func (a *Application) Addr() string {
	if a.server == nil {
		return ""
	}

	return a.server.Addr
}
