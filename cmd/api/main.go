package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/mangalores/case-studies-voiceline/internal/bootstrap"
	"github.com/mangalores/case-studies-voiceline/internal/config"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		slog.Error("load config", "error", err)
		os.Exit(1)
	}
	slog.Debug("loaded config", "config", cfg)

	application, err := bootstrap.BuildApplication(cfg)
	if err != nil {
		logger.Error("build application", "error", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := application.Run(ctx); err != nil {
		logger.Error("run application", "error", err)
		os.Exit(1)
	}
}
