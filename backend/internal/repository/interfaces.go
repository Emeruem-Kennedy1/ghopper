package repository

import "github.com/Emeruem-Kennedy1/ghopper/internal/models"

// UserRepositoryInterface defines the methods we use from UserRepository
type UserRepositoryInterface interface {
	Create(user *models.User) error
	UpsertUser(user *models.User) error
	GetByID(id string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(user *models.User) error
	Delete(id string) error
}

// song repository interface

type SongRepositoryInterface interface {
	GetSongIDsByTitleAndArtist(title, artist string) ([]int, error)
	GetSongWithDetails(SongID int) (*models.SongNode, error)
	GetAllSampledSongs(songID int) ([]int, error)
	FindSongsByGenreBFS(songQueries []models.SongQuery, targetGenre string, maxDepth int) ([]models.SearchResult, error)
}

type SpotifySongRepositoryInterface interface {
	FindSongByNameAndArtist(name, artist string) (*models.Song, error)
	SaveSong(song *models.Song) error
	FindPlaylistByNameAndUser(name, userID string) (*models.Playlist, error)
	SavePlaylist(playlist *models.Playlist) error
	FindPlaylistByIDAndUser(playlistID, userID string) (*models.Playlist, error)
	DeletePlaylist(playlist *models.Playlist) error
	GetUserPlaylists(userID string) ([]models.Playlist, error)
	UpdatePlaylistImageURL(playlistImageURL string, playlist *models.Playlist) error
	DeleteUserPlaylists(userID string) error
}

// NonSpotifyUserRepositoryInterface defines the methods for the NonSpotifyUserRepository
type NonSpotifyUserRepositoryInterface interface {
	FindByID(id string) (*models.NonSpotifyUser, error)
	Create(user *models.NonSpotifyUser) error
	Verify(id, passphrase string) (bool, error)
	SavePlaylist(playlist *models.NonSpotifyPlaylist, tracks []models.NonSpotifyPlaylistTrack, seedTracks []models.NonSpotifyPlaylistSeedTrack) error
	GetUserPlaylists(userID string) ([]models.NonSpotifyPlaylist, error)
	GetPlaylistWithTracks(playlistID string) (*models.NonSpotifyPlaylistWithTracks, error)
	UpdateTrackStatus(trackID string, addedToPlaylist bool) error
	DeletePlaylist(playlistID string) error
}

// Ensure the UserRepository, SpotifySongRepository and SongRepository implement our interfaces
var _ UserRepositoryInterface = (*UserRepository)(nil)
var _ SongRepositoryInterface = (*SongRepository)(nil)
var _ SpotifySongRepositoryInterface = (*SpotifySongRepository)(nil)
var _ NonSpotifyUserRepositoryInterface = (*NonSpotifyUserRepository)(nil)
