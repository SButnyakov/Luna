package services

import (
	"context"
	"sync"

	"github.com/SButnyakov/luna/audio-processing/config"
	"github.com/SButnyakov/luna/audio-processing/internal/lib/logger"
)

type configService struct {
	config *config.Config
	logger *logger.Logger
	mu     sync.RWMutex
}

func NewConfigService(ctx context.Context, config *config.Config) *configService {
	return &configService{
		config: config,
		logger: logger.FromContext(ctx),
	}
}

func (s *configService) GetAudioConfig() config.AudioConfig {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.config.Audio
}
