package mongo

import (
	"context"

	"github.com/SButnyakov/luna/audio-upload/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoStorage struct {
	client *mongo.Client
}

func Connect(ctx context.Context, cfg config.MongoConfig) (*mongo.Database, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.URI))
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return client.Database(cfg.Database), nil
}

func (s *mongoStorage) Disconnect(ctx context.Context) error {
	return s.client.Disconnect(ctx)
}
