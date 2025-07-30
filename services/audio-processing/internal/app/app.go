package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/SButnyakov/luna/audio-processing/config"
	"github.com/gin-gonic/gin"
)

func Run() {
	cfg := config.MustLoad()

	router := gin.Default()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := router.Run(fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port)); err != nil {
			log.Fatalf("failed to run server: %v", err)
		}
	}()
}
