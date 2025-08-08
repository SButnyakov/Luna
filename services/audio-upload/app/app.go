package app

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/SButnyakov/luna/audio-upload/api/http/handlers"
	"github.com/SButnyakov/luna/audio-upload/config"
	applogger "github.com/SButnyakov/luna/audio-upload/internal/logger"
	"github.com/SButnyakov/luna/audio-upload/internal/storage/mongo"
	"github.com/SButnyakov/luna/audio-upload/internal/usecases"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Run() {
	// Config
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("failed to load config:", err)
	}

	// Logger
	logger := applogger.NewLogger(cfg.Env)

	logger.Info("logger initialized")

	cancelCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctxWithLogger := applogger.WithContext(cancelCtx, logger)

	// Storages connection
	mongoDB, err := mongo.Connect(ctxWithLogger, cfg.Mongo)
	if err != nil {
		logger.Error("failed to connect mongo db", slog.Any("err", err)) // TODO: fatal or panic level
		os.Exit(1)
	}

	logger.Info("mongo db connected")

	// TODO: uncomment when repo implemented
	// s3Storage, err := s3.Connect(cfg.S3)
	// if err != nil {
	// 	logger.Error("failed to connect mongo db", slog.Any("err", err)) // TODO: fatal or panic level
	// 	os.Exit(1)
	// }

	logger.Info("s3 storage connected")

	// Repos
	tracksRepo := mongo.NewTracksRepo(ctxWithLogger, mongoDB, cfg.Mongo.TracksCollection)

	// Usecases
	uploadTrackUsecase := usecases.NewUploadTrackUsecase(ctxWithLogger, tracksRepo)
	_ = uploadTrackUsecase // TODO: remove when handlers implemented

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Get("/upload_track", handlers.UploadTrack(ctxWithLogger, uploadTrackUsecase, 30*time.Second)) // TODO
	http.ListenAndServe(fmt.Sprintf(":%d", 8101), r)
}
