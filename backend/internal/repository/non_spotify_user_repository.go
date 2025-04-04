package repository

import (
	"errors"

	"github.com/Emeruem-Kennedy1/ghopper/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// NonSpotifyUserRepository handles database operations for non-Spotify users
type NonSpotifyUserRepository struct {
	db *gorm.DB
}

// NewNonSpotifyUserRepository creates a new repository instance
func NewNonSpotifyUserRepository(db *gorm.DB) *NonSpotifyUserRepository {
	return &NonSpotifyUserRepository{db: db}
}

// FindByID retrieves a non-Spotify user by ID
func (r *NonSpotifyUserRepository) FindByID(id string) (*models.NonSpotifyUser, error) {
	var user models.NonSpotifyUser
	result := r.db.Where("id = ?", id).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

// Create adds a new non-Spotify user
func (r *NonSpotifyUserRepository) Create(user *models.NonSpotifyUser) error {
	return r.db.Create(user).Error
}

// Verify checks if the passphrase matches for a given user ID
func (r *NonSpotifyUserRepository) Verify(id, passphrase string) (bool, error) {
	var user models.NonSpotifyUser
	result := r.db.Where("id = ?", id).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, result.Error
	}
	return user.Passphrase == passphrase, nil
}

// SavePlaylist creates a new playlist with tracks and seed tracks
func (r *NonSpotifyUserRepository) SavePlaylist(
	playlist *models.NonSpotifyPlaylist,
	tracks []models.NonSpotifyPlaylistTrack,
	seedTracks []models.NonSpotifyPlaylistSeedTrack,
) error {
	// Generate IDs if not provided
	if playlist.ID == "" {
		playlist.ID = uuid.New().String()
	}

	// Use a transaction to ensure all operations succeed or fail together
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Save the playlist
		if err := tx.Create(playlist).Error; err != nil {
			return err
		}

		// Save the tracks
		for i := range tracks {
			tracks[i].PlaylistID = playlist.ID
			if tracks[i].ID == "" {
				tracks[i].ID = uuid.New().String()
			}
			if err := tx.Create(&tracks[i]).Error; err != nil {
				return err
			}
		}

		// Save the seed tracks
		for i := range seedTracks {
			seedTracks[i].PlaylistID = playlist.ID
			if seedTracks[i].ID == "" {
				seedTracks[i].ID = uuid.New().String()
			}
			if err := tx.Create(&seedTracks[i]).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// GetUserPlaylists retrieves all playlists for a user
func (r *NonSpotifyUserRepository) GetUserPlaylists(userID string) ([]models.NonSpotifyPlaylist, error) {
	var playlists []models.NonSpotifyPlaylist
	result := r.db.Where("user_id = ?", userID).Find(&playlists)
	return playlists, result.Error
}

// GetPlaylistWithTracks retrieves a playlist with its tracks and seed tracks
func (r *NonSpotifyUserRepository) GetPlaylistWithTracks(playlistID string) (*models.NonSpotifyPlaylistWithTracks, error) {
	var playlist models.NonSpotifyPlaylist
	result := r.db.Where("id = ?", playlistID).First(&playlist)
	if result.Error != nil {
		return nil, result.Error
	}

	var tracks []models.NonSpotifyPlaylistTrack
	result = r.db.Where("playlist_id = ?", playlistID).Find(&tracks)
	if result.Error != nil {
		return nil, result.Error
	}

	var seedTracks []models.NonSpotifyPlaylistSeedTrack
	result = r.db.Where("playlist_id = ?", playlistID).Find(&seedTracks)
	if result.Error != nil {
		return nil, result.Error
	}

	return &models.NonSpotifyPlaylistWithTracks{
		NonSpotifyPlaylist: playlist,
		Tracks:             tracks,
		SeedTracks:         seedTracks,
	}, nil
}

// UpdateTrackStatus updates the "added to playlist" status for a track
func (r *NonSpotifyUserRepository) UpdateTrackStatus(trackID string, addedToPlaylist bool) error {
	return r.db.Model(&models.NonSpotifyPlaylistTrack{}).
		Where("id = ?", trackID).
		Update("added_to_playlist", addedToPlaylist).
		Error
}

// DeletePlaylist deletes a playlist with its tracks and seed tracks
func (r *NonSpotifyUserRepository) DeletePlaylist(playlistID string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Delete the tracks
		if err := tx.Where("playlist_id = ?", playlistID).Delete(&models.NonSpotifyPlaylistTrack{}).Error; err != nil {
			return err
		}

		// Delete the seed tracks
		if err := tx.Where("playlist_id = ?", playlistID).Delete(&models.NonSpotifyPlaylistSeedTrack{}).Error; err != nil {
			return err
		}

		// Delete the playlist
		return tx.Where("id = ?", playlistID).Delete(&models.NonSpotifyPlaylist{}).Error
	})
}
