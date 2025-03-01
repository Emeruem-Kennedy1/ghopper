package handlers

import (
	"net/http"

	"github.com/Emeruem-Kennedy1/ghopper/internal/auth"
	"github.com/Emeruem-Kennedy1/ghopper/internal/models"
	"github.com/Emeruem-Kennedy1/ghopper/internal/repository"
	"github.com/Emeruem-Kennedy1/ghopper/internal/services"
	"github.com/stretchr/testify/mock"
	"github.com/zmb3/spotify"
)

// ! Mock for SpotifyAuth
type MockSpotifyAuth struct {
	mock.Mock
}

// Ensure the mock implements the interface
var _ auth.SpotifyAuthInterface = (*MockSpotifyAuth)(nil)

func (m *MockSpotifyAuth) AuthURL() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockSpotifyAuth) CallBack(r *http.Request) (*spotify.Client, error) {
	args := m.Called(r)
	// Don't try to convert to *spotify.Client, just return the raw value
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	// This is a real *spotify.Client or a manually created one
	return args.Get(0).(*spotify.Client), args.Error(1)
}

func (m *MockSpotifyAuth) GetUserInfo(client *spotify.Client) (*spotify.PrivateUser, error) {
	args := m.Called(client)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*spotify.PrivateUser), args.Error(1)
}

func (m *MockSpotifyAuth) GetAuthenticator() auth.AuthenticatorInterface {
	args := m.Called()
	return args.Get(0).(auth.AuthenticatorInterface)
}

// ! Mock for UserRepository
type MockUserRepository struct {
	mock.Mock
}

// Ensure the mock implements the interface
var _ repository.UserRepositoryInterface = (*MockUserRepository)(nil)

func (m *MockUserRepository) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) UpsertUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(id string) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// ! Mock Spotify client for testing
type MockSpotifyClient struct {
	mock.Mock
}

// Ensure the mock implements the interface
var _ services.SpotifyClientInterface = (*MockSpotifyClient)(nil)

func (m *MockSpotifyClient) Search(query string, t spotify.SearchType) (*spotify.SearchResult, error) {
	args := m.Called(query, t)
	return args.Get(0).(*spotify.SearchResult), args.Error(1)
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

func (m *MockSpotifyClient) CurrentUsersTopTracksOpt(opt *spotify.Options) (*spotify.FullTrackPage, error) {
	args := m.Called(opt)
	return args.Get(0).(*spotify.FullTrackPage), args.Error(1)
}

func (m *MockSpotifyClient) CurrentUsersTopArtistsOpt(opts *spotify.Options) (*spotify.FullArtistPage, error) {
	args := m.Called(opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*spotify.FullArtistPage), args.Error(1)
}

// ! MockSongRepository is a mock implementation of the SongRepository
type MockSongRepository struct {
	mock.Mock
}

func (m *MockSongRepository) FindSongsByGenreBFS(songQueries []models.SongQuery, targetGenre string, maxDepth int) ([]models.SearchResult, error) {
	args := m.Called(songQueries, targetGenre, maxDepth)
	return args.Get(0).([]models.SearchResult), args.Error(1)
}

func (m *MockSongRepository) GetAllSampledSongs(limit int) ([]int, error) {
	args := m.Called(limit)
	return args.Get(0).([]int), args.Error(1)
}

func (m *MockSongRepository) GetSongWithDetails(SongID int) (*models.SongNode, error) {
	args := m.Called(SongID)
	return args.Get(0).(*models.SongNode), args.Error(1)
}
func (m *MockSongRepository) GetSongIDsByTitleAndArtist(title, artist string) ([]int, error) {
	args := m.Called(title, artist)
	return args.Get(0).([]int), args.Error(1)
}

// ! MockSpotifyService for testing
type MockSpotifyService struct {
	mock.Mock
}

func (m *MockSpotifyService) GetSongURL(userID, name, artist string) (string, error) {
	args := m.Called(userID, name, artist)
	return args.String(0), args.Error(1)
}

func (m *MockSpotifyService) CreatePlaylistFromSongs(userID string, songSpotifyIDs []spotify.ID, playlistName, playlistDescription string) (string, error) {
	args := m.Called(userID, songSpotifyIDs, playlistName, playlistDescription)
	return args.String(0), args.Error(1)
}

func (m *MockSpotifyService) DeletePlaylist(userID, playlistID string) error {
	args := m.Called(userID, playlistID)
	return args.Error(0)
}

func (m *MockSpotifyService) GetPlaylistImageURL(userID, playlistID string) (string, error) {
	args := m.Called(userID, playlistID)
	return args.String(0), args.Error(1)
}

// ! MockClientManager for testing
type MockClientManager struct {
	mock.Mock
}

func (m *MockClientManager) GetClient(userID string) (services.SpotifyClientInterface, bool) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, false
	}
	return args.Get(0).(services.SpotifyClientInterface), args.Bool(1)
}

func (m *MockClientManager) StoreClient(userID string, client services.SpotifyClientInterface) {
	m.Called(userID, client)
}
func (m *MockClientManager) DeleteClient(userID string) {
	m.Called(userID)
}
func (m *MockClientManager) RemoveClient(userID string) {
	m.Called(userID)
}

// ! MockSpotifySongRepository for testing
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
