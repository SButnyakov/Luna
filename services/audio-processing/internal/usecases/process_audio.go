package usecases

import (
	"context"
	"io"

	"github.com/SButnyakov/luna/audio-processing/config"
	"github.com/SButnyakov/luna/audio-processing/internal/dto"
	"github.com/SButnyakov/luna/audio-processing/internal/lib/logger"
	"github.com/SButnyakov/luna/audio-processing/internal/models"
)

type audioFileProvider interface {
	GetFile(ctx context.Context, key string) (io.ReadCloser, error)
	SaveHLSPlaylistBatch(ctx context.Context, playlists []models.HLSPlaylist) error
}

type audioConfigProvider interface {
	GetAudioConfig() config.AudioConfig
}

type processAudioUsecase struct {
	fileProvider   audioFileProvider
	configProvider audioConfigProvider
	logger         *logger.Logger
}

func NewProcessAudioUsecase(
	ctx context.Context,
	fileProvider audioFileProvider,
	configProvider audioConfigProvider,
) *processAudioUsecase {
	return &processAudioUsecase{
		fileProvider:   fileProvider,
		configProvider: configProvider,
		logger:         logger.FromContext(ctx),
	}
}

func (u *processAudioUsecase) Execute(audio dto.ProcessAudioDTO) (map[int]string, error) {
	// TODO: получить аудио

	// TODO: получить битрейт

	// TODO: ffmpeg - сконвертировать в нужные битрейты

	// TODO: ffmpeg - сконвертировать в плейлисты

	// TODO: загрузить плейлисты

	return nil, nil
}
