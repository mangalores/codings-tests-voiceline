package bootstrap

import (
	"fmt"
	"log/slog"
	"net/http"

	exporteradapter "github.com/mangalores/case-studies-voiceline/internal/adapters/exporter"
	extractoradapter "github.com/mangalores/case-studies-voiceline/internal/adapters/extractor"
	transcriberadapter "github.com/mangalores/case-studies-voiceline/internal/adapters/transcriber"
	"github.com/mangalores/case-studies-voiceline/internal/app"
	"github.com/mangalores/case-studies-voiceline/internal/config"
	apphttp "github.com/mangalores/case-studies-voiceline/internal/http"
	"github.com/mangalores/case-studies-voiceline/internal/repository"
)

const channelSize = 100

func BuildApplication(cfg *config.AppConfig) (*app.Application, error) {
	recordingRepository := newRecordingRepository(cfg)
	transcribeCommands := make(chan app.TranscribeCommand, channelSize)
	extractCommands := make(chan app.ExtractCommand, channelSize)
	exportCommands := make(chan app.ExportCommand, channelSize)

	uploader := newUploadService(recordingRepository, transcribeCommands)
	router := apphttp.NewRouter(uploader)

	server := newServer(cfg, router)

	transcriber, err := buildTranscriber(cfg)
	if err != nil {
		return nil, err
	}

	extractor, err := buildExtractor(cfg)
	if err != nil {
		return nil, err
	}

	exporter := buildExporter(cfg)

	transcriptionWorker := app.NewTranscriptionWorker(recordingRepository, transcriber, transcribeCommands, extractCommands)
	extractionWorker := app.NewExtractionWorker(recordingRepository, extractor, extractCommands, exportCommands)
	exportWorker := app.NewExportWorker(recordingRepository, exporter, exportCommands)

	return app.NewApplication(server, transcriptionWorker, extractionWorker, exportWorker), nil
}

func newRecordingRepository(cfg *config.AppConfig) *repository.RecordingRepository {
	return repository.NewRecordingRepository(cfg.FileStoreConfig.StoragePath)
}

func newUploadService(recordings app.RecordingStorer, transcribeCommands chan<- app.TranscribeCommand) *app.UploadService {
	return app.NewUploadService(recordings, transcribeCommands)
}

func newServer(cfg *config.AppConfig, router http.Handler) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: router,
	}
}

func buildTranscriber(cfg *config.AppConfig) (app.Transcriber, error) {
	slog.Info("config:", "transcriber", cfg.Transcriber)

	switch cfg.Transcriber {
	case "openai":
		return transcriberadapter.NewOpenAITranscriber(cfg.OpenAIConfig.APIKey), nil
	default:
		return transcriberadapter.NewMockTranscriber(cfg.MockConfig.TranscriptionPath)
	}
}

func buildExtractor(cfg *config.AppConfig) (app.Extractor, error) {
	switch cfg.Extractor {
	case "openai":
		return extractoradapter.NewOpenAIExtractor(cfg.OpenAIConfig.APIKey), nil
	default:
		return extractoradapter.NewMockExtractor(cfg.MockConfig.ExtractionPath), nil
	}
}

func buildExporter(cfg *config.AppConfig) app.Exporter {
	switch cfg.Exporter {
	case "webhook":
		return exporteradapter.NewWebhookExporter(cfg.WebHookConfig.ExportURL)
	case "googlesheets":
		return exporteradapter.NewGoogleSheetExporter(
			cfg.GoogleConfig.CredentialsFilePath,
			cfg.GoogleConfig.SheetID,
			cfg.GoogleConfig.SheetRange,
		)
	default:
		return exporteradapter.NewMockExporter()
	}
}
