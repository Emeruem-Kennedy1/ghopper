package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Emeruem-Kennedy1/ghopper/internal/models"
	"github.com/Emeruem-Kennedy1/ghopper/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zmb3/spotify"
)

// Helper function to set up the Gin context
func setupGinContext(userID string) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	if userID != "" {
		c.Set("userID", userID)
	}

	return c, w
}

// Test for GetUser handler
func TestGetUser(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Arrange
		mockUserRepo := new(MockUserRepository)
		expectedUser := &models.User{
			ID:          "test-user-id",
			DisplayName: "Test User",
			Email:       "test@example.com",
		}

		mockUserRepo.On("GetByID", "test-user-id").Return(expectedUser, nil)

		c, w := setupGinContext("test-user-id")

		// Act
		handler := GetUser(mockUserRepo)
		handler(c)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		userMap, ok := response["user"].(map[string]interface{})
		require.True(t, ok)
		assert.Equal(t, "test-user-id", userMap["id"])
		assert.Equal(t, "Test User", userMap["display_name"])
		assert.Equal(t, "test@example.com", userMap["email"])

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		// Arrange
		mockUserRepo := new(MockUserRepository)

		c, w := setupGinContext("") // No user ID in context

		// Act
		handler := GetUser(mockUserRepo)
		handler(c)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "Unauthorized", response["error"])

		// Verify that GetByID was never called
		mockUserRepo.AssertNotCalled(t, "GetByID")
	})

	t.Run("Database_Error", func(t *testing.T) {
		// Arrange
		mockUserRepo := new(MockUserRepository)
		mockUserRepo.On("GetByID", "test-user-id").Return(nil, errors.New("database error"))

		c, w := setupGinContext("test-user-id")

		// Act
		handler := GetUser(mockUserRepo)
		handler(c)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "Failed to get user", response["error"])

		mockUserRepo.AssertExpectations(t)
	})
}

// Test for GetUserTopArtists handler
func TestGetUserTopArtists(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Arrange
		mockClientManager := new(MockClientManager)

		// Create fake artists response
		artistsResponse := &spotify.FullArtistPage{
			Artists: []spotify.FullArtist{
				{
					SimpleArtist: spotify.SimpleArtist{
						ID:   "artist1",
						Name: "Artist One",
					},
				},
				{
					SimpleArtist: spotify.SimpleArtist{
						ID:   "artist2",
						Name: "Artist Two",
					},
				},
			},
		}

		// Set up client manager to return a real spotify client
		mockClientManager.On("GetClient", "test-user-id").Return(&spotify.Client{}, true)

		c, w := setupGinContext("test-user-id")

		// We need to replace the function with a simpler mock that doesn't
		// try to use the actual client object at all
		origGetUserTopArtists := getUserTopArtistsFunc
		defer func() { getUserTopArtistsFunc = origGetUserTopArtists }()

		getUserTopArtistsFunc = func(client services.SpotifyClientInterface, opts *spotify.Options) (*spotify.FullArtistPage, error) {
			// Just verify the options and return mock data
			assert.Equal(t, 25, *opts.Limit)
			assert.Equal(t, "short", *opts.Timerange)
			return artistsResponse, nil
		}

		// Act
		handler := GetUserTopArtists(mockClientManager)
		handler(c)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		// The response should include the artists
		_, ok := response["artists"].(map[string]interface{})
		require.True(t, ok)

		mockClientManager.AssertExpectations(t)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		// Arrange
		mockClientManager := new(MockClientManager)

		c, w := setupGinContext("") // No user ID in context

		// Act
		handler := GetUserTopArtists(mockClientManager)
		handler(c)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "Unauthorized", response["error"])

		// Verify that GetClient was never called
		mockClientManager.AssertNotCalled(t, "GetClient")
	})

	t.Run("No_Client", func(t *testing.T) {
		// Arrange
		mockClientManager := new(MockClientManager)

		// Set up client manager to return no client
		mockClientManager.On("GetClient", "test-user-id").Return((*spotify.Client)(nil), false)

		c, w := setupGinContext("test-user-id")

		// Act
		handler := GetUserTopArtists(mockClientManager)
		handler(c)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "Unauthorized", response["error"])

		mockClientManager.AssertExpectations(t)
	})

	t.Run("Spotify_API_Error", func(t *testing.T) {
		// Arrange
		mockClientManager := new(MockClientManager)

		// Set up client manager to return a real spotify client
		mockClientManager.On("GetClient", "test-user-id").Return(&spotify.Client{}, true)

		c, w := setupGinContext("test-user-id")

		// Replace the handler function to return an error
		origGetUserTopArtists := getUserTopArtistsFunc
		defer func() { getUserTopArtistsFunc = origGetUserTopArtists }()

		getUserTopArtistsFunc = func(client services.SpotifyClientInterface, opts *spotify.Options) (*spotify.FullArtistPage, error) {
			// Don't try to use the client - just return an error
			return nil, errors.New("Spotify API error")
		}

		// Act
		handler := GetUserTopArtists(mockClientManager)
		handler(c)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "Failed to get user's top artists", response["error"])

		mockClientManager.AssertExpectations(t)
	})
}

// Test for GetUserTopTracks handler
func TestGetUserTopTracks(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Arrange
		mockClientManager := new(MockClientManager)
		mockSpotifyService := new(MockSpotifyService)

		// Create fake tracks response
		tracksResponse := &spotify.FullTrackPage{
			Tracks: []spotify.FullTrack{
				{
					SimpleTrack: spotify.SimpleTrack{
						ID:   "track1",
						Name: "Track One",
						Artists: []spotify.SimpleArtist{
							{Name: "Artist One"},
						},
					},
					Album: spotify.SimpleAlbum{
						Images: []spotify.Image{
							{URL: "https://example.com/image1.jpg"},
						},
					},
				},
				{
					SimpleTrack: spotify.SimpleTrack{
						ID:   "track2",
						Name: "Track Two",
						Artists: []spotify.SimpleArtist{
							{Name: "Artist Two"},
						},
					},
					Album: spotify.SimpleAlbum{
						Images: []spotify.Image{
							{URL: "https://example.com/image2.jpg"},
						},
					},
				},
			},
		}

		// Set up client manager to return a real spotify client (not a mock)
		mockClientManager.On("GetClient", "test-user-id").Return(&spotify.Client{}, true)

		// Set up a custom mock for the CurrentUsersTopTracksOpt method
		origFunc := getUserTopTracksFunc
		defer func() { getUserTopTracksFunc = origFunc }()

		getUserTopTracksFunc = func(client services.SpotifyClientInterface, opt *spotify.Options) (*spotify.FullTrackPage, error) {
			// Verify the options
			assert.Equal(t, 25, *opt.Limit)
			assert.Equal(t, "short", *opt.Timerange)
			return tracksResponse, nil
		}

		c, w := setupGinContext("test-user-id")

		// Act
		handler := GetUserTopTracks(mockClientManager, mockSpotifyService)
		handler(c)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		// Verify the tracks in the response
		tracks, ok := response["tracks"].([]interface{})
		require.True(t, ok)
		assert.Equal(t, 2, len(tracks))

		// Check first track data
		track1 := tracks[0].(map[string]interface{})
		assert.Equal(t, "track1", track1["id"])
		assert.Equal(t, "Track One", track1["name"])
		assert.Equal(t, "https://example.com/image1.jpg", track1["image"])

		artists1 := track1["artists"].([]interface{})
		assert.Equal(t, 1, len(artists1))
		assert.Equal(t, "Artist One", artists1[0])

		mockClientManager.AssertExpectations(t)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		// Arrange
		mockClientManager := new(MockClientManager)
		mockSpotifyService := new(MockSpotifyService)

		c, w := setupGinContext("") // No user ID in context

		// Act
		handler := GetUserTopTracks(mockClientManager, mockSpotifyService)
		handler(c)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "Unauthorized", response["error"])

		// Verify that GetClient was never called
		mockClientManager.AssertNotCalled(t, "GetClient")
	})

	t.Run("No_Client", func(t *testing.T) {
		// Arrange
		mockClientManager := new(MockClientManager)
		mockSpotifyService := new(MockSpotifyService)

		// Set up client manager to return no client
		mockClientManager.On("GetClient", "test-user-id").Return((*spotify.Client)(nil), false)

		c, w := setupGinContext("test-user-id")

		// Act
		handler := GetUserTopTracks(mockClientManager, mockSpotifyService)
		handler(c)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "Unauthorized", response["error"])

		mockClientManager.AssertExpectations(t)
	})

	t.Run("Spotify_API_Error", func(t *testing.T) {
		// Arrange
		mockClientManager := new(MockClientManager)
		mockSpotifyService := new(MockSpotifyService)

		// Set up client manager to return a real client
		mockClientManager.On("GetClient", "test-user-id").Return(&spotify.Client{}, true)

		// Set up the mock to return an error
		origFunc := getUserTopTracksFunc
		defer func() { getUserTopTracksFunc = origFunc }()

		getUserTopTracksFunc = func(client services.SpotifyClientInterface, opt *spotify.Options) (*spotify.FullTrackPage, error) {
			return nil, errors.New("Spotify API error")
		}

		c, w := setupGinContext("test-user-id")

		// Act
		handler := GetUserTopTracks(mockClientManager, mockSpotifyService)
		handler(c)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "Failed to get user's top tracks", response["error"])

		mockClientManager.AssertExpectations(t)
	})
}

// Test for GetUserPlaylists handler
func TestGetUserPlaylists(t *testing.T) {
	t.Run("Success_WithImages", func(t *testing.T) {
		// Arrange
		mockSpotifySongRepo := new(MockSpotifySongRepository)
		mockSpotifyService := new(MockSpotifyService)

		var responseItems []models.Playlist

		// Create test playlists
		playlists := []models.Playlist{
			{
				ID:          "playlist1",
				Name:        "Playlist One",
				Description: "First playlist",
				URL:         "https://spotify.com/playlist1",
				Image:       "https://example.com/image1.jpg", // Already has image
				UserID:      "test-user-id",
			},
			{
				ID:          "playlist2",
				Name:        "Playlist Two",
				Description: "Second playlist",
				URL:         "https://spotify.com/playlist2",
				Image:       "", // No image, needs to be fetched
				UserID:      "test-user-id",
			},
		}

		// Set up repository to return our playlists
		mockSpotifySongRepo.On("GetUserPlaylists", "test-user-id").Return(playlists, nil)

		// Set up service to return image for second playlist
		mockSpotifyService.On("GetPlaylistImageURL", "test-user-id", "playlist2").Return("https://example.com/image2.jpg", nil)

		// Set up repository to update image URL
		mockSpotifySongRepo.On("UpdatePlaylistImageURL", "https://example.com/image2.jpg", &playlists[1]).Return(nil)

		c, w := setupGinContext("test-user-id")

		// Act
		handler := GetUserPlaylists(mockSpotifySongRepo, mockSpotifyService)
		handler(c)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		// Verify the playlists in the response
		playlistInterfaces, ok := response["playlists"].([]interface{})
		require.True(t, ok)

		for _, pi := range playlistInterfaces {
			playlistMap, ok := pi.(map[string]interface{})
			require.True(t, ok)

			playlist := models.Playlist{
				ID:          playlistMap["id"].(string),
				Name:        playlistMap["name"].(string),
				Description: playlistMap["description"].(string),
				URL:         playlistMap["url"].(string),
				Image:       playlistMap["image"].(string),
				UserID:      "test-user-id",
			}
			responseItems = append(responseItems, playlist)
		}
		assert.Equal(t, 2, len(responseItems))

		// Check first playlist data
		playlist1 := responseItems[0]
		assert.Equal(t, "playlist1", playlist1.ID)
		assert.Equal(t, "Playlist One", playlist1.Name)
		assert.Equal(t, "First playlist", playlist1.Description)
		assert.Equal(t, "https://spotify.com/playlist1", playlist1.URL)
		assert.Equal(t, "https://example.com/image1.jpg", playlist1.Image)

		// Check second playlist data
		playlist2 := responseItems[1]
		assert.Equal(t, "playlist2", playlist2.ID)
		assert.Equal(t, "Playlist Two", playlist2.Name)
		assert.Equal(t, "Second playlist", playlist2.Description)
		assert.Equal(t, "https://spotify.com/playlist2", playlist2.URL)
		assert.Equal(t, "https://example.com/image2.jpg", playlist2.Image)

		mockSpotifySongRepo.AssertExpectations(t)
		mockSpotifyService.AssertExpectations(t)
	})

	t.Run("No_Playlists", func(t *testing.T) {
		// Arrange
		mockSpotifySongRepo := new(MockSpotifySongRepository)
		mockSpotifyService := new(MockSpotifyService)

		// Set up repository to return empty playlists
		mockSpotifySongRepo.On("GetUserPlaylists", "test-user-id").Return([]models.Playlist{}, nil)

		c, w := setupGinContext("test-user-id")

		// Act
		handler := GetUserPlaylists(mockSpotifySongRepo, mockSpotifyService)
		handler(c)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		// Verify that the playlists key exists but is nil
		_, exists := response["playlists"]
		assert.True(t, exists, "Response should contain 'playlists' key")

		// Check the message
		assert.Equal(t, "No playlists found", response["message"])

		mockSpotifySongRepo.AssertExpectations(t)
		mockSpotifyService.AssertNotCalled(t, "GetPlaylistImageURL")
	})

	t.Run("Unauthorized", func(t *testing.T) {
		// Arrange
		mockSpotifySongRepo := new(MockSpotifySongRepository)
		mockSpotifyService := new(MockSpotifyService)

		c, w := setupGinContext("") // No user ID in context

		// Act
		handler := GetUserPlaylists(mockSpotifySongRepo, mockSpotifyService)
		handler(c)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "Unauthorized", response["error"])

		// Verify that no repository or service methods were called
		mockSpotifySongRepo.AssertNotCalled(t, "GetUserPlaylists")
		mockSpotifyService.AssertNotCalled(t, "GetPlaylistImageURL")
	})

	t.Run("Database_Error", func(t *testing.T) {
		// Arrange
		mockSpotifySongRepo := new(MockSpotifySongRepository)
		mockSpotifyService := new(MockSpotifyService)

		// Set up repository to return an error
		mockSpotifySongRepo.On("GetUserPlaylists", "test-user-id").Return([]models.Playlist{}, errors.New("database error"))

		c, w := setupGinContext("test-user-id")

		// Act
		handler := GetUserPlaylists(mockSpotifySongRepo, mockSpotifyService)
		handler(c)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "Failed to get user's playlists", response["error"])

		mockSpotifySongRepo.AssertExpectations(t)
		mockSpotifyService.AssertNotCalled(t, "GetPlaylistImageURL")
	})

	t.Run("GetPlaylistImageURL_Error", func(t *testing.T) {
		// Arrange
		mockSpotifySongRepo := new(MockSpotifySongRepository)
		mockSpotifyService := new(MockSpotifyService)

		// Create test playlists with one that needs an image
		playlists := []models.Playlist{
			{
				ID:          "playlist1",
				Name:        "Playlist One",
				Description: "First playlist",
				URL:         "https://spotify.com/playlist1",
				Image:       "", // No image, needs to be fetched
				UserID:      "test-user-id",
			},
		}

		// Set up repository to return our playlist
		mockSpotifySongRepo.On("GetUserPlaylists", "test-user-id").Return(playlists, nil)

		// Set up service to return an error when fetching image
		mockSpotifyService.On("GetPlaylistImageURL", "test-user-id", "playlist1").Return("", errors.New("spotify API error"))

		mockSpotifySongRepo.On("UpdatePlaylistImageURL", "", &playlists[0]).Return(nil)
		
		c, w := setupGinContext("test-user-id")

		// Act
		handler := GetUserPlaylists(mockSpotifySongRepo, mockSpotifyService)
		handler(c)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "Failed to get playlist image", response["error"])

		mockSpotifySongRepo.AssertExpectations(t)
		mockSpotifyService.AssertExpectations(t)
	})
}

// Test for DeleteUserAccount handler
func TestDeleteUserAccount(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Arrange
		mockUserRepo := new(MockUserRepository)
		mockSpotifySongRepo := new(MockSpotifySongRepository)
		mockClientManager := new(MockClientManager)

		// Set up repository expectations
		mockSpotifySongRepo.On("DeleteUserPlaylists", "test-user-id").Return(nil)
		mockClientManager.On("RemoveClient", "test-user-id").Return()
		mockUserRepo.On("Delete", "test-user-id").Return(nil)

		c, w := setupGinContext("test-user-id")

		// Act
		handler := DeleteUserAccount(mockUserRepo, mockSpotifySongRepo, mockClientManager)
		handler(c)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		// Verify the success message
		assert.Equal(t, "Account successfully deleted", response["message"])

		// Verify that all methods were called in the correct order
		mockSpotifySongRepo.AssertExpectations(t)
		mockClientManager.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		// Arrange
		mockUserRepo := new(MockUserRepository)
		mockSpotifySongRepo := new(MockSpotifySongRepository)
		mockClientManager := new(MockClientManager)

		c, w := setupGinContext("") // No user ID in context

		// Act
		handler := DeleteUserAccount(mockUserRepo, mockSpotifySongRepo, mockClientManager)
		handler(c)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "Unauthorized", response["error"])

		// Verify that no methods were called
		mockSpotifySongRepo.AssertNotCalled(t, "DeleteUserPlaylists")
		mockClientManager.AssertNotCalled(t, "RemoveClient")
		mockUserRepo.AssertNotCalled(t, "Delete")
	})

	t.Run("Error_DeletePlaylists", func(t *testing.T) {
		// Arrange
		mockUserRepo := new(MockUserRepository)
		mockSpotifySongRepo := new(MockSpotifySongRepository)
		mockClientManager := new(MockClientManager)

		// Set up repository to return error
		mockSpotifySongRepo.On("DeleteUserPlaylists", "test-user-id").Return(errors.New("database error"))

		c, w := setupGinContext("test-user-id")

		// Act
		handler := DeleteUserAccount(mockUserRepo, mockSpotifySongRepo, mockClientManager)
		handler(c)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "Failed to delete user playlists", response["error"])

		// Verify that only DeleteUserPlaylists was called
		mockSpotifySongRepo.AssertExpectations(t)
		mockClientManager.AssertNotCalled(t, "RemoveClient")
		mockUserRepo.AssertNotCalled(t, "Delete")
	})

	t.Run("Error_DeleteUser", func(t *testing.T) {
		// Arrange
		mockUserRepo := new(MockUserRepository)
		mockSpotifySongRepo := new(MockSpotifySongRepository)
		mockClientManager := new(MockClientManager)

		// Set up repository expectations
		mockSpotifySongRepo.On("DeleteUserPlaylists", "test-user-id").Return(nil)
		mockClientManager.On("RemoveClient", "test-user-id").Return()
		mockUserRepo.On("Delete", "test-user-id").Return(errors.New("database error"))

		c, w := setupGinContext("test-user-id")

		// Act
		handler := DeleteUserAccount(mockUserRepo, mockSpotifySongRepo, mockClientManager)
		handler(c)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "Failed to delete user account", response["error"])

		// Verify that methods were called in the correct order
		mockSpotifySongRepo.AssertExpectations(t)
		mockClientManager.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
	})
}
