package repository

import (
	"context"
	"fmt"

	"github.com/SButnyakov/luna/audio-upload/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TrackRepository interface {
	Create(ctx context.Context, track *models.Track) error
	GetByID(ctx context.Context, id string) (*models.Track, error) 
	GetAll(ctx context.Context, limit, offset int) ([]*models.Track, error)
	Update(ctx context.Context, track *models.Track) error
	Delete(ctx context.Context, id string) error
	GetByArtistID(ctx context.Context, artistID string) ([]*models.Track, error)
	GetByStatus(ctx context.Context, status models.ProcessingStatus) ([]*models.Track, error)
	UpdateStatus(ctx context.Context, id string, status models.ProcessingStatus) error
	UpdateM3U8Playlists(ctx context.Context, id string, playlists map[int]string) error
}

type trackRepository struct {
	collection *mongo.Collection
}

func NewTrackRepository(db *mongo.Database, collectionName string) TrackRepository {
	return &trackRepository{
		collection: db.Collection(collectionName),
	}
}

func (r *trackRepository) Create(ctx context.Context, track *models.Track) error {
	result, err := r.collection.InsertOne(ctx, track)
	if err != nil {
		return fmt.Errorf("failed to create track: %w", err)
	}
	
	track.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *trackRepository) GetByID(ctx context.Context, id string) (*models.Track, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid track ID: %w", err)
	}

	var track models.Track
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&track)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("track not found")
		}
		return nil, fmt.Errorf("failed to get track: %w", err)
	}

	return &track, nil
}

func (r *trackRepository) GetAll(ctx context.Context, limit, offset int) ([]*models.Track, error) {
	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}
	if offset > 0 {
		opts.SetSkip(int64(offset))
	}
	opts.SetSort(bson.D{{Key: "uploaded_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get tracks: %w", err)
	}
	defer cursor.Close(ctx)

	var tracks []*models.Track
	for cursor.Next(ctx) {
		var track models.Track
		if err := cursor.Decode(&track); err != nil {
			return nil, fmt.Errorf("failed to decode track: %w", err)
		}
		tracks = append(tracks, &track)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return tracks, nil
}

func (r *trackRepository) Update(ctx context.Context, track *models.Track) error {
	filter := bson.M{"_id": track.ID}
	update := bson.M{"$set": track}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update track: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("track not found")
	}

	return nil
}

func (r *trackRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid track ID: %w", err)
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("failed to delete track: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("track not found")
	}

	return nil
}

func (r *trackRepository) GetByArtistID(ctx context.Context, artistID string) ([]*models.Track, error) {
	objectID, err := primitive.ObjectIDFromHex(artistID)
	if err != nil {
		return nil, fmt.Errorf("invalid artist ID: %w", err)
	}

	filter := bson.M{"artist_ids": objectID}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get tracks by artist: %w", err)
	}
	defer cursor.Close(ctx)

	var tracks []*models.Track
	for cursor.Next(ctx) {
		var track models.Track
		if err := cursor.Decode(&track); err != nil {
			return nil, fmt.Errorf("failed to decode track: %w", err)
		}
		tracks = append(tracks, &track)
	}

	return tracks, nil
}

func (r *trackRepository) GetByStatus(ctx context.Context, status models.ProcessingStatus) ([]*models.Track, error) {
	filter := bson.M{"status": status}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get tracks by status: %w", err)
	}
	defer cursor.Close(ctx)

	var tracks []*models.Track
	for cursor.Next(ctx) {
		var track models.Track
		if err := cursor.Decode(&track); err != nil {
			return nil, fmt.Errorf("failed to decode track: %w", err)
		}
		tracks = append(tracks, &track)
	}

	return tracks, nil
}

func (r *trackRepository) UpdateStatus(ctx context.Context, id string, status models.ProcessingStatus) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid track ID: %w", err)
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{"status": status}}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update track status: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("track not found")
	}

	return nil
}

func (r *trackRepository) UpdateM3U8Playlists(ctx context.Context, id string, playlists map[int]string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid track ID: %w", err)
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{"m3u8_playlists": playlists}}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update m3u8 playlists: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("track not found")
	}

	return nil
}