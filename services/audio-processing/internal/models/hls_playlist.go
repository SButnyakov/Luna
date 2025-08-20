package models

import "io"

type HLSPlaylist struct {
	AudioID  string
	Bitrate  int
	Segments []HLSSegment
	Playlist io.Reader
}

type HLSSegment struct {
	Index    int
	Duration float64
	Data     io.Reader
}
