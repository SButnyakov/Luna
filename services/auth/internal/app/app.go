package app

import (
	"context"
	"log"
	"log/slog"

	"github.com/SButnyakov/luna/id/internal/config"
	applogger "github.com/SButnyakov/luna/id/internal/lib/logger"
	"github.com/SButnyakov/luna/id/internal/storage/postgres"
)

func Run() {
	// Config
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Logger
	logger := applogger.NewLogger(cfg.Env)

	logger.Info("logger initialized")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx = applogger.WithContext(ctx, logger)

	pg, err := postgres.Connect(ctx, cfg.PG)
	if err != nil {
		log.Fatal(err)
	}
	_ = pg
	logger.Info("postgres connected")
	logger.Debug("postgres args", slog.String("uri", pg.Config().ConnString()))
}
