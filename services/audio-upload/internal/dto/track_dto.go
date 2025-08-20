package dto

import (
	"io"
	"time"
)

type UploadTrackDTO struct {
	Title       string         `json:"title"`
	ArtistIDs   []string       `json:"artist_ids"`
	Genres      []string       `json:"genres"`
	Duration    int            `json:"duration"`
	ReleaseDate time.Time      `json:"release_date,omitempty"`
	File        io.Reader      `json:"-"`
	FileName    string         `json:"file_name"`
	Covers      map[int]string `json:"covers"`
}

type UpdateTrackStatusDTO struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type UpdateTrackM3U8PlaylistsDTO struct {
	ID        string         `json:"id"`
	Playlists map[int]string `json:"playlists"`
}

type GetTracksDTO struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}
