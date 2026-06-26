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

	exporter, err := buildExporter(cfg)
	if err != nil {
		return nil, err
	}

	transcriptionWorker := app.NewTranscriptionWorker(recordingRepository, transcriber, transcribeCommands, extractCommands)
	extractionWorker := app.NewExtractionWorker(recordingRepository, extractor, extractCommands, exportCommands)
	exportWorker := app.NewExportWorker(recordingRepository, exporter, exportCommands)

	return app.NewApplication(server, transcriptionWorker, extractionWorker, exportWorker), nil
}

func newRecordingRepository(cfg *config.AppConfig) *repository.RecordingRepository {
	return repository.NewRecordingRepository(cfg.StoragePath)
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
		openAIConfig, err := transcriberadapter.LoadOpenAIConfig()
		if err != nil {
			return nil, err
		}

		return transcriberadapter.NewOpenAITranscriber(openAIConfig), nil
	default:
		mockConfig, err := transcriberadapter.LoadMockConfig()
		if err != nil {
			return nil, err
		}

		return transcriberadapter.NewMockTranscriber(mockConfig)
	}
}

func buildExtractor(cfg *config.AppConfig) (app.Extractor, error) {
	switch cfg.Extractor {
	case "openai":
		openAIConfig, err := extractoradapter.LoadOpenAIConfig()
		if err != nil {
			return nil, err
		}

		return extractoradapter.NewOpenAIExtractor(openAIConfig), nil
	default:
		mockConfig, err := extractoradapter.LoadMockConfig()
		if err != nil {
			return nil, err
		}

		return extractoradapter.NewMockExtractor(mockConfig), nil
	}
}

func buildExporter(cfg *config.AppConfig) (app.Exporter, error) {
	switch cfg.Exporter {
	case "webhook":
		webhookConfig, err := exporteradapter.LoadWebhookConfig()
		if err != nil {
			return nil, err
		}

		return exporteradapter.NewWebhookExporter(webhookConfig), nil
	case "googlesheets":
		googleConfig, err := exporteradapter.LoadGoogleConfig()
		if err != nil {
			return nil, err
		}

		return exporteradapter.NewGoogleSheetExporter(googleConfig), nil
	default:
		return exporteradapter.NewMockExporter(), nil
	}
}
