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
	// Check if song already exists
	var existingSong models.Song
	result := r.db.Where("id = ?", song.ID).First(&existingSong)

	if result.Error == gorm.ErrRecordNotFound {
		result = r.db.Create(song)
		if result.Error != nil {
			return fmt.Errorf("error saving song: %v", result.Error)
		}
		return nil
	} else if result.Error != nil {
		return fmt.Errorf("error checking for existing song: %v", result.Error)
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

func (r *SpotifySongRepository) FindPlaylistByIDAndUser(playlistID, userID string) (*models.Playlist, error) {
	var playlist models.Playlist
	result := r.db.Where("id = ? AND user_id = ?", playlistID, userID).First(&playlist)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &playlist, nil
}

func (r *SpotifySongRepository) DeletePlaylist (playlist *models.Playlist) error {
	result := r.db.Delete(playlist)
	if result.Error != nil {
		return fmt.Errorf("error deleting playlist: %v", result.Error)
	}
	return nil
}

func (r *SpotifySongRepository) GetUserPlaylists(userID string) ([]models.Playlist, error) {
	var playlists []models.Playlist
	result := r.db.Where("user_id = ?", userID).Find(&playlists)
	if result.Error != nil {
		return nil, fmt.Errorf("error getting user playlists: %v", result.Error)
	}
	return playlists, nil
}

func (r *SpotifySongRepository) UpdatePlaylistImageURL(playlistImageURL string, playlist *models.Playlist) error {
	result := r.db.Model(playlist).Update("image", playlistImageURL)
	if result.Error != nil {
		return fmt.Errorf("error updating playlist image: %v", result.Error)
	}
	return nil
}

func (r *SpotifySongRepository) DeleteUserPlaylists(userID string) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.Playlist{}).Error
}