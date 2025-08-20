package repository

import (
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/SButnyakov/luna/audio-processing/config"
	"github.com/SButnyakov/luna/audio-processing/internal/lib/logger"
	"github.com/SButnyakov/luna/audio-processing/internal/models"
)

type fileStorage interface {
	Download(context.Context, string) (io.ReadCloser, error)
	Upload(context.Context, string, io.Reader) error
	UploadMultiple(context.Context, map[string]io.Reader) error
}

type fileRepository struct {
	storage fileStorage
	logger  *logger.Logger
	cfg     config.S3AudioConfig
}

func NewFileRepository(ctx context.Context, storage fileStorage, cfg config.S3AudioConfig) *fileRepository { // TODO: config service
	return &fileRepository{
		storage: storage,
		logger:  logger.FromContext(ctx),
	}
}

func (r *fileRepository) GetFile(ctx context.Context, key string) (io.ReadCloser, error) {
	const op = "repository.fileRepository.GetFile"
	r.logger.DebugWithOp(op, "downloading file", slog.String("key", key))
	file, err := r.storage.Download(ctx, key)
	if err != nil {
		r.logger.ErrorWithOp(op, "failed to download file", err, slog.String("key", key))
		return nil, err
	}
	return file, nil
}

func (r *fileRepository) SaveFile(ctx context.Context, key string, body io.Reader) error {
	const op = "repository.fileRepository.SaveFile"
	r.logger.DebugWithOp(op, "uploading file", slog.String("key", key))
	if err := r.storage.Upload(ctx, key, body); err != nil {
		r.logger.ErrorWithOp(op, "failed to upload file", err, slog.String("key", key))
		return err
	}
	return nil
}

func (r *fileRepository) SaveHLSPlaylist(ctx context.Context, hls models.HLSPlaylist) error {
	const op = "repository.fileRepository.SaveHLSPlaylist"
	r.logger.Info("uploading HLS playlist",
		slog.String("audioID", hls.AudioID),
		slog.Int("bitrate", hls.Bitrate),
		"segments_count", len(hls.Segments),
	)

	objects := make(map[string]io.Reader)

	for _, segment := range hls.Segments {
		key := fmt.Sprintf("%s/%s/%d/segment_%03d.ts", r.cfg.HLSKey, hls.AudioID, hls.Bitrate, segment.Index)
		r.logger.DebugWithOp(op, "preparing segment for upload", slog.String("key", key), slog.Int("index", segment.Index))
		objects[key] = segment.Data
	}

	playlistKey := fmt.Sprintf("%s/%s/%d/playlist.m3u8", r.cfg.HLSKey, hls.AudioID, hls.Bitrate)
	objects[playlistKey] = hls.Playlist

	if err := r.storage.UploadMultiple(ctx, objects); err != nil {
		r.logger.ErrorWithOp(op, "failed to upload HLS playlist and segments", err,
			slog.String("audioID", hls.AudioID), slog.Int("bitrate", hls.Bitrate))

	}

	r.logger.Info("HLS playlist saved successfully", slog.String("audioID", hls.AudioID), slog.Int("bitrate", hls.Bitrate))
	return nil
}

func (r *fileRepository) SaveHLSPlaylistBatch(ctx context.Context, playlists []models.HLSPlaylist) error {
	const op = "repository.fileRepository.SaveHLSPlaylistBatch"
	r.logger.Info("uploading batch of HLS playlists", slog.Int("count", len(playlists)))

	for _, hls := range playlists {
		if err := r.SaveHLSPlaylist(ctx, hls); err != nil {
			r.logger.ErrorWithOp(op, "failed to save HLS playlist", err,
				slog.String("audioID", hls.AudioID), slog.Int("bitrate", hls.Bitrate))
			return err
		}
	}

	r.logger.Info("all HLS playlists saved successfully")
	return nil
}
