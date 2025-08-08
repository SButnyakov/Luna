package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/SButnyakov/luna/audio-upload/internal/dto"
	"github.com/go-chi/render"
)

type trackUploadUsecase interface {
	Execute(context.Context, dto.UploadTrackDTO) (string, error)
}

func UploadTrack(ctx context.Context, uploader trackUploadUsecase, timeout time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.handlers.upload_track.UploadTrack"
		timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		var req dto.UploadTrackDTO
		err := render.DefaultDecoder(r, &req)
		if err != nil {
			// TODO
		}

		id, err := uploader.Execute(timeoutCtx, req)
		if err != nil {
			// TODO
		}

		_ = id
		// TODO
	}
}
