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
	"github.com/SButnyakov/luna/audio-processing/internal/services"
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

	// services
	fileService := services.NewFileService()
	configService := services.NewConfigService(ctx, cfg)

	// usecases
	processAudioUsecase := usecases.NewProcessAudioUsecase(ctx, fileService, configService)

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
