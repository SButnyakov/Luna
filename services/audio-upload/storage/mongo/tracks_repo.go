package mongo

import (
	"context"
	"fmt"

	"github.com/SButnyakov/luna/audio-upload/logger"
	"github.com/SButnyakov/luna/audio-upload/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type tracksRepository struct {
	c      *mongo.Collection
	logger *logger.Logger
}

func NewTracksRepo(ctx context.Context, db *mongo.Database, collectionName string) *tracksRepository {
	return &tracksRepository{
		c:      db.Collection(collectionName),
		logger: logger.FromContext(ctx),
	}
}

func (r *tracksRepository) InsertOne(ctx context.Context, track models.Track) (string, error) {
	const op = "mongo.tracks_repo.InsertOne"

	res, err := r.c.InsertOne(ctx, track)
	if err != nil {
		return "", fmt.Errorf("%s: failed to insert track into database: %w", op, err)
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (r *tracksRepository) UpdateByID(ctx context.Context, id string, fields map[string]interface{}) error {
	const op = "mongo.tracks_repo.UpdateByID"

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("%s: invalid id format: %w", op, err)
	}

	filter := bson.M{"_id": objID}
	update := bson.M{"$set": fields}

	res, err := r.c.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return fmt.Errorf("%s: no document found with id %s", op, id)
	}

	return nil
}
