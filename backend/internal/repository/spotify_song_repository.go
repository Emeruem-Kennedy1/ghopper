package repository

import (
	"errors"
	"fmt"

	"github.com/Emeruem-Kennedy1/ghopper/internal/models"
	"gorm.io/gorm"
)

type SpotifySongRepository struct {
	db *gorm.DB
}

func NewSpotifySongRepository(db *gorm.DB) *SpotifySongRepository {
	return &SpotifySongRepository{db: db}
}

func (r *SpotifySongRepository) FindSongByNameAndArtist(name, artist string) (*models.Song, error) {
	var song models.Song
	result := r.db.Where("song_name = ? AND artist_name = ?", name, artist).First(&song)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &song, nil
}

func (r *SpotifySongRepository) SaveSong(song *models.Song) error {
	result := r.db.Create(song)
	if result.Error != nil {
		return fmt.Errorf("error saving song: %v", result.Error)
	}
	return nil
}

func (r *SpotifySongRepository) FindPlaylistByNameAndUser(name, userID string) (*models.Playlist, error) {
	var playlist models.Playlist
	result := r.db.Where("playlist_name = ? AND user_id = ?", name, userID).First(&playlist)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &playlist, nil
}

func (r *SpotifySongRepository) SavePlaylist(playlist *models.Playlist) error {
	result := r.db.Create(playlist)
	if result.Error != nil {
		return fmt.Errorf("error saving playlist: %v", result.Error)
	}
	return nil
}