package converter

import (
	"fmt"
	"os"

	"github.com/dhowden/tag"
)

// Metadata はMP4/M4Aファイルのメタデータを表す
type Metadata struct {
	Title       string
	Artist      string
	Album       string
	AlbumArtist string
	Composer    string
	Genre       string
	Year        int
	Track       int
	TrackTotal  int
	Disc        int
	DiscTotal   int
	Picture     *tag.Picture
}

// ReadMetadata はMP4ファイルからメタデータを読み取る
func ReadMetadata(filePath string) (*Metadata, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	m, err := tag.ReadFrom(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata: %w", err)
	}

	track, trackTotal := m.Track()
	disc, discTotal := m.Disc()

	metadata := &Metadata{
		Title:       m.Title(),
		Artist:      m.Artist(),
		Album:       m.Album(),
		AlbumArtist: m.AlbumArtist(),
		Composer:    m.Composer(),
		Genre:       m.Genre(),
		Year:        m.Year(),
		Track:       track,
		TrackTotal:  trackTotal,
		Disc:        disc,
		DiscTotal:   discTotal,
		Picture:     m.Picture(),
	}

	return metadata, nil
}
