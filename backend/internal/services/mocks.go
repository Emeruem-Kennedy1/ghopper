package services

import (
	"fmt"

	"github.com/Emeruem-Kennedy1/ghopper/internal/models"
	"github.com/stretchr/testify/mock"
	"github.com/zmb3/spotify"
)

// ! MockSpotifySongRepository mocks the SpotifySongRepository
type MockSpotifySongRepository struct {
	mock.Mock
}

func (m *MockSpotifySongRepository) FindSongByNameAndArtist(name, artist string) (*models.Song, error) {
	args := m.Called(name, artist)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Song), args.Error(1)
}

func (m *MockSpotifySongRepository) SaveSong(song *models.Song) error {
	args := m.Called(song)
	return args.Error(0)
}

func (m *MockSpotifySongRepository) FindPlaylistByNameAndUser(name, userID string) (*models.Playlist, error) {
	args := m.Called(name, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Playlist), args.Error(1)
}

func (m *MockSpotifySongRepository) SavePlaylist(playlist *models.Playlist) error {
	args := m.Called(playlist)
	return args.Error(0)
}

func (m *MockSpotifySongRepository) FindPlaylistByIDAndUser(playlistID, userID string) (*models.Playlist, error) {
	args := m.Called(playlistID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Playlist), args.Error(1)
}

func (m *MockSpotifySongRepository) DeletePlaylist(playlist *models.Playlist) error {
	args := m.Called(playlist)
	return args.Error(0)
}

func (m *MockSpotifySongRepository) GetUserPlaylists(userID string) ([]models.Playlist, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.Playlist), args.Error(1)
}

func (m *MockSpotifySongRepository) UpdatePlaylistImageURL(playlistImageURL string, playlist *models.Playlist) error {
	args := m.Called(playlistImageURL, playlist)
	return args.Error(0)
}

func (m *MockSpotifySongRepository) DeleteUserPlaylists(userID string) error {
	args := m.Called(userID)
	return args.Error(0)
}

// ! Mock Spotify client for testing
type MockSpotifyClient struct {
	mock.Mock
}

func (m *MockSpotifyClient) Search(query string, searchType spotify.SearchType) (*spotify.SearchResult, error) {
	// Convert searchType to int explicitly for matching
	args := m.Called(query, int(searchType))

	if result, ok := args.Get(0).(*spotify.SearchResult); ok {
		return result, args.Error(1)
	}

	return nil, fmt.Errorf("mock: failed to convert search result")
}

func (m *MockSpotifyClient) CurrentUser() (*spotify.PrivateUser, error) {
	args := m.Called()
	return args.Get(0).(*spotify.PrivateUser), args.Error(1)
}

func (m *MockSpotifyClient) CreatePlaylistForUser(userID, name, description string, public bool) (*spotify.FullPlaylist, error) {
	args := m.Called(userID, name, description, public)
	return args.Get(0).(*spotify.FullPlaylist), args.Error(1)
}

func (m *MockSpotifyClient) GetPlaylistTracks(playlistID spotify.ID) (*spotify.PlaylistTrackPage, error) {
	args := m.Called(playlistID)
	return args.Get(0).(*spotify.PlaylistTrackPage), args.Error(1)
}

func (m *MockSpotifyClient) AddTracksToPlaylist(playlistID spotify.ID, trackIDs ...spotify.ID) (string, error) {
	args := m.Called(playlistID, trackIDs)
	return args.String(0), args.Error(1)
}

func (m *MockSpotifyClient) GetPlaylist(playlistID spotify.ID) (*spotify.FullPlaylist, error) {
	args := m.Called(playlistID)
	return args.Get(0).(*spotify.FullPlaylist), args.Error(1)
}

func (m *MockSpotifyClient) UnfollowPlaylist(userID, playlistID spotify.ID) error {
	args := m.Called(userID, playlistID)
	return args.Error(0)
}

func (m *MockSpotifyClient) CurrentUsersTopArtistsOpt(opt *spotify.Options) (*spotify.FullArtistPage, error) {
	args := m.Called(opt)
	return args.Get(0).(*spotify.FullArtistPage), args.Error(1)
}

func (m *MockSpotifyClient) CurrentUsersTopTracksOpt(opt *spotify.Options) (*spotify.FullTrackPage, error) {
	args := m.Called(opt)
	return args.Get(0).(*spotify.FullTrackPage), args.Error(1)
}
