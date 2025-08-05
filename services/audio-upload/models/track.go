package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProcessingStatus string

const (
	StatusPending    ProcessingStatus = "pending"
	StatusProcessing ProcessingStatus = "processing"
	StatusCompleted  ProcessingStatus = "completed"
	StatusFailed     ProcessingStatus = "failed"
)

type Track struct {
	ID           primitive.ObjectID    `bson:"_id,omitempty" json:"id"`
	Title        string                `bson:"title" json:"title"`
	ArtistIDs    []primitive.ObjectID  `bson:"artist_ids" json:"artist_ids"`
	Duration     int                   `bson:"duration" json:"duration"` // in seconds
	TrackURL     string                `bson:"track_url" json:"track_url"`
	UploadedAt   time.Time             `bson:"uploaded_at" json:"uploaded_at"`
	ReleaseDate  *time.Time            `bson:"release_date,omitempty" json:"release_date,omitempty"`
	Status       ProcessingStatus      `bson:"status" json:"status"`
	M3U8Playlists map[int]string       `bson:"m3u8_playlists" json:"m3u8_playlists"` // bitrate -> m3u8 URL
}