package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/SButnyakov/luna/audio-processing/config"
	httphandlers "github.com/SButnyakov/luna/audio-processing/internal/http_handlers"
	applogger "github.com/SButnyakov/luna/audio-processing/internal/lib/logger"
	"github.com/SButnyakov/luna/audio-processing/internal/repository"
	"github.com/SButnyakov/luna/audio-processing/internal/services"
	"github.com/SButnyakov/luna/audio-processing/internal/storage"
	"github.com/SButnyakov/luna/audio-processing/internal/usecases"
	"github.com/gin-gonic/gin"
)

func Run() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	logger := applogger.NewLogger(cfg.Env)
	ctx := applogger.WithContext(context.Background(), logger)

	// storage
	s3Storage, err := storage.NewS3Storage(cfg.S3.Storage)
	if err != nil {
		log.Fatalf("failed to create S3 storage: %v", err)
	}

	// repositories
	fileRepository := repository.NewFileRepository(ctx, s3Storage, cfg.S3.Audio)

	// services
	configService := services.NewConfigService(ctx, cfg)

	// usecases
	processAudioUsecase := usecases.NewProcessAudioUsecase(ctx, fileRepository, configService)

	// router
	router := gin.Default()
	api := router.Group("/api")

	// v1
	v1 := api.Group(cfg.API.V1.Root)
	{
		audio := v1.Group(cfg.API.V1.Audio.Root)
		audio.POST(cfg.API.V1.Audio.Process, httphandlers.ProcessAudio(processAudioUsecase.Execute))
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := router.Run(fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port)); err != nil {
			log.Fatalf("failed to run server: %v", err)
		}
	}()

	<-interrupt

	logger.Info("shutting down")
}
