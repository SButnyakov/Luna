package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/SButnyakov/luna/audio-upload/internal/dto"
	"github.com/SButnyakov/luna/audio-upload/internal/lib/converters"
	"github.com/SButnyakov/luna/audio-upload/internal/logger"
	"github.com/SButnyakov/luna/audio-upload/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type trackInfoProvider interface {
	InsertOne(context.Context, models.Track) (string, error)
}

type uploadTrackUsecase struct {
	infoProvider trackInfoProvider
	logger       *logger.Logger
}

func NewUploadTrackUsecase(ctx context.Context, trackInfoProvider trackInfoProvider) *uploadTrackUsecase {
	return &uploadTrackUsecase{
		infoProvider: trackInfoProvider,
		logger:       logger.FromContext(ctx),
	}
}

func (u *uploadTrackUsecase) Execute(ctx context.Context, track dto.UploadTrackDTO) (string, error) { // TODO: add file
	const op = "usecases.upload_track.Execute"

	// TODO: process file into mp3 and add meta info

	trackModel, err := uploadTrackDTOToModel(track)
	if err != nil {
		return "", fmt.Errorf("%s: failed to convert dto to model: %w", op, err)
	}

	loadedID, err := u.infoProvider.InsertOne(ctx, *trackModel) // TODO: retry
	if err != nil {
		return "", fmt.Errorf("%s: failed to create new track: %w", op, err)
	}

	// TODO: send event into kafka

	return loadedID, nil
}

func uploadTrackDTOToModel(track dto.UploadTrackDTO) (*models.Track, error) {
	const op = "usecases.upload_track.uploadTrackDTOToModel"

	artists, err := converters.StringArrToObjectID(track.ArtistIDs)
	if err != nil {
		return nil, fmt.Errorf("%s: invalid artists: %w", op, err)
	}

	genres, err := converters.StringArrToObjectID(track.Genres)
	if err != nil {
		return nil, fmt.Errorf("%s: invalid genres: %w", op, err)
	}

	trackID := primitive.NewObjectID()

	return &models.Track{
		ID:          trackID,
		Title:       track.Title,
		ArtistIDs:   artists,
		Genres:      genres,
		ReleaseDate: track.ReleaseDate,
		UploadedAt:  time.Now().UTC(),
		Covers:      track.Covers,
		Status:      models.StatusProcessing,
	}, nil
}
