package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Emeruem-Kennedy1/ghopper/internal/models"
	"github.com/Emeruem-Kennedy1/ghopper/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zmb3/spotify"
)

func setupSongHandlerTest(songRepo repository.SongRepositoryInterface) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Add a mock context middleware to simulate authenticated user
	r.Use(func(c *gin.Context) {
		c.Set("userID", "test-user-id")
		c.Next()
	})

	r.POST("/search", SearchSongByGenre(songRepo))
	return r
}

func setupAnalyzeSongsTest(songRepo *MockSongRepository, clientManager *MockClientManager, spotifyService *MockSpotifyService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Add a mock context middleware to simulate authenticated user
	r.Use(func(c *gin.Context) {
		c.Set("userID", "test-user-id")
		c.Next()
	})

	r.POST("/toptracks-analysis", AnalyzeSongsGivenGenre(songRepo, clientManager, spotifyService))
	return r
}

func setupDeletePlaylistTest(spotifyService *MockSpotifyService, spotifySongRepo *MockSpotifySongRepository) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Add a mock context middleware to simulate authenticated user
	r.Use(func(c *gin.Context) {
		c.Set("userID", "test-user-id")
		c.Next()
	})

	r.DELETE("/user/playlists/:playlistID", DeletePlaylist(spotifyService, spotifySongRepo))
	return r
}

func TestSearchSongByGenre(t *testing.T) {
	t.Run("Successful_Search", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockSongRepository)
		r := setupSongHandlerTest(mockRepo)

		// Prepare mock search request
		searchRequest := SongSearchRequest{
			Songs: []models.SongQuery{
				{Title: "Test Song", Artist: "Test Artist"},
			},
			Genre:    "rock",
			MaxDepth: 2,
		}

		// Prepare mock search results
		mockResults := []models.SearchResult{
			{
				SourceSong: models.SongNode{
					ID:    1,
					Title: "Test Song",
					Artists: []models.Artist{
						{ID: 1, Name: "Test Artist", IsMain: true},
					},
					Genres: []string{"rock"},
				},
				MatchedSong: models.SongNode{
					ID:    2,
					Title: "Another Song",
					Artists: []models.Artist{
						{ID: 2, Name: "Another Artist", IsMain: true},
					},
					Genres: []string{"rock"},
				},
				Distance: 1,
				Path: []models.SongNode{
					{
						ID:    1,
						Title: "Test Song",
					},
					{
						ID:    2,
						Title: "Another Song",
					},
				},
			},
		}

		// Setup mock expectations
		mockRepo.On("FindSongsByGenreBFS", searchRequest.Songs, searchRequest.Genre, searchRequest.MaxDepth).
			Return(mockResults, nil)

		// Convert request to JSON
		jsonRequest, _ := json.Marshal(searchRequest)

		// Create test HTTP request
		req := httptest.NewRequest("POST", "/search", bytes.NewBuffer(jsonRequest))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		// Serve the request
		r.ServeHTTP(resp, req)

		// Assert
		assert.Equal(t, http.StatusOK, resp.Code, "Should return OK status")

		// Parse the response
		var graphResponse GraphResponse
		err := json.Unmarshal(resp.Body.Bytes(), &graphResponse)
		require.NoError(t, err, "Should parse response JSON")

		// Verify graph response contents
		assert.Len(t, graphResponse.Paths, 1, "Should have one path")
		assert.Len(t, graphResponse.Nodes, 2, "Should have two nodes")
		assert.NotEmpty(t, graphResponse.AdjacencyList, "Should have adjacency list")

		// Verify mock expectations
		mockRepo.AssertExpectations(t)
	})

	t.Run("Missing_Genre", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockSongRepository)
		r := setupSongHandlerTest(mockRepo)

		// Prepare search request without genre
		searchRequest := SongSearchRequest{
			Songs: []models.SongQuery{
				{Title: "Test Song", Artist: "Test Artist"},
			},
		}

		// Convert request to JSON
		jsonRequest, _ := json.Marshal(searchRequest)

		// Create test HTTP request
		req := httptest.NewRequest("POST", "/search", bytes.NewBuffer(jsonRequest))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		// Serve the request
		r.ServeHTTP(resp, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, resp.Code, "Should return Bad Request status")

		// Parse the response
		var errorResponse map[string]string
		err := json.Unmarshal(resp.Body.Bytes(), &errorResponse)
		require.NoError(t, err, "Should parse error response JSON")

		assert.Contains(t, errorResponse["error"], "genre is required", "Error message should indicate missing genre")
	})

	t.Run("Search_Error", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockSongRepository)
		r := setupSongHandlerTest(mockRepo)

		// Prepare search request
		searchRequest := SongSearchRequest{
			Songs: []models.SongQuery{
				{Title: "Test Song", Artist: "Test Artist"},
			},
			Genre:    "rock",
			MaxDepth: 2,
		}

		// Setup mock to return an error
		mockRepo.On("FindSongsByGenreBFS", searchRequest.Songs, searchRequest.Genre, searchRequest.MaxDepth).
			Return([]models.SearchResult{}, assert.AnError)

		// Convert request to JSON
		jsonRequest, _ := json.Marshal(searchRequest)

		// Create test HTTP request
		req := httptest.NewRequest("POST", "/search", bytes.NewBuffer(jsonRequest))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		// Serve the request
		r.ServeHTTP(resp, req)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, resp.Code, "Should return Internal Server Error status")

		// Parse the response
		var errorResponse map[string]string
		err := json.Unmarshal(resp.Body.Bytes(), &errorResponse)
		require.NoError(t, err, "Should parse error response JSON")

		assert.Contains(t, errorResponse["error"], "failed to search songs", "Error message should indicate search failure")

		// Verify mock expectations
		mockRepo.AssertExpectations(t)
	})
}

func TestAnalyzeSongsGivenGenre(t *testing.T) {
	t.Run("Successful_Analysis", func(t *testing.T) {
		// Arrange
		mockSongRepo := new(MockSongRepository)
		mockClientManager := new(MockClientManager)
		mockSpotifyService := new(MockSpotifyService)
		r := setupAnalyzeSongsTest(mockSongRepo, mockClientManager, mockSpotifyService)

		// Mock Spotify client
		mockClient := new(MockSpotifyClient)
		mockTracks := &spotify.FullTrackPage{
			Tracks: []spotify.FullTrack{
				{
					SimpleTrack: spotify.SimpleTrack{
						Name: "Test Song",
						Artists: []spotify.SimpleArtist{
							{Name: "Test Artist"},
						},
					},
				},
			},
		}

		// Prepare mock search request
		analysisRequest := TopTracksAnalysisRequest{
			Genre: "rock",
		}

		// Setup mock expectations
		mockClientManager.On("GetClient", "test-user-id").Return(mockClient, true)
		mockClient.On("CurrentUsersTopTracksOpt", mock.Anything).Return(mockTracks, nil)

		mockSongRepo.On("FindSongsByGenreBFS", mock.Anything, "rock", 2).Return(
			[]models.SearchResult{
				{
					MatchedSong: models.SongNode{
						Title: "Matched Song",
						Artists: []models.Artist{
							{Name: "Matched Artist"},
						},
					},
				},
			}, nil)

		mockSpotifyService.On("GetSongURL", "test-user-id", "Matched Song", "Matched Artist").
			Return("https://open.spotify.com/track/123", nil)

		mockSpotifyService.On("CreatePlaylistFromSongs",
			"test-user-id",
			mock.Anything,
			mock.Anything,
			mock.Anything).
			Return("https://open.spotify.com/playlist/test-playlist", nil)

		// Convert request to JSON
		jsonRequest, _ := json.Marshal(analysisRequest)

		// Create test HTTP request
		req := httptest.NewRequest("POST", "/toptracks-analysis", bytes.NewBuffer(jsonRequest))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		// Serve the request
		r.ServeHTTP(resp, req)

		// Assert
		assert.Equal(t, http.StatusOK, resp.Code, "Should return OK status")

		// Parse the response
		var topTracksResponse TopTracksAnalysisResponse
		err := json.Unmarshal(resp.Body.Bytes(), &topTracksResponse)
		require.NoError(t, err, "Should parse response JSON")

		// Verify response contents
		assert.NotEmpty(t, topTracksResponse.Playlist, "Should have a playlist URL")
		assert.NotEmpty(t, topTracksResponse.Songs, "Should have songs")

		// Verify mock expectations
		mockClientManager.AssertExpectations(t)
		mockSongRepo.AssertExpectations(t)
		mockSpotifyService.AssertExpectations(t)
	})

	t.Run("Missing_Genre", func(t *testing.T) {
		// Arrange
		mockSongRepo := new(MockSongRepository)
		mockClientManager := new(MockClientManager)
		mockSpotifyService := new(MockSpotifyService)
		r := setupAnalyzeSongsTest(mockSongRepo, mockClientManager, mockSpotifyService)

		// Mock the GetClient method to return a mock client
		mockClient := new(MockSpotifyClient)
		mockClientManager.On("GetClient", "test-user-id").Return(mockClient, true)

		// Prepare mock search request without genre
		analysisRequest := TopTracksAnalysisRequest{}

		// Convert request to JSON
		jsonRequest, _ := json.Marshal(analysisRequest)

		// Create test HTTP request
		req := httptest.NewRequest("POST", "/toptracks-analysis", bytes.NewBuffer(jsonRequest))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		// Serve the request
		r.ServeHTTP(resp, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, resp.Code, "Should return Bad Request status")

		// Parse the response
		var errorResponse map[string]string
		err := json.Unmarshal(resp.Body.Bytes(), &errorResponse)
		require.NoError(t, err, "Should parse error response JSON")

		assert.Contains(t, errorResponse["error"], "genre is required", "Error message should indicate missing genre")

		// Verify mock expectations
		mockClientManager.AssertExpectations(t)
	})
}

func TestDeletePlaylist(t *testing.T) {
	t.Run("Successful_Delete", func(t *testing.T) {
		// Arrange
		mockSpotifyService := new(MockSpotifyService)
		mockSpotifySongRepo := new(MockSpotifySongRepository)
		r := setupDeletePlaylistTest(mockSpotifyService, mockSpotifySongRepo)

		// Prepare test playlist
		playlistID := "test-playlist-id"
		playlist := &models.Playlist{
			ID:     playlistID,
			UserID: "test-user-id",
		}

		// Setup mock expectations
		mockSpotifySongRepo.On("FindPlaylistByIDAndUser", playlistID, "test-user-id").Return(playlist, nil)
		mockSpotifySongRepo.On("DeletePlaylist", playlist).Return(nil)
		mockSpotifyService.On("DeletePlaylist", "test-user-id", playlistID).Return(nil)

		// Create test HTTP request
		req := httptest.NewRequest("DELETE", "/user/playlists/"+playlistID, nil)
		resp := httptest.NewRecorder()

		// Serve the request
		r.ServeHTTP(resp, req)

		// Assert
		assert.Equal(t, http.StatusOK, resp.Code, "Should return OK status")

		// Parse the response
		var response map[string]string
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		require.NoError(t, err, "Should parse response JSON")

		assert.Contains(t, response["message"], "playlist deleted", "Should have success message")

		// Verify mock expectations
		mockSpotifySongRepo.AssertExpectations(t)
		mockSpotifyService.AssertExpectations(t)
	})

	t.Run("Playlist_Not_Found", func(t *testing.T) {
		// Arrange
		mockSpotifyService := new(MockSpotifyService)
		mockSpotifySongRepo := new(MockSpotifySongRepository)
		r := setupDeletePlaylistTest(mockSpotifyService, mockSpotifySongRepo)

		// Prepare test playlist
		playlistID := "test-playlist-id"

		// Setup mock expectations
		mockSpotifySongRepo.On("FindPlaylistByIDAndUser", playlistID, "test-user-id").Return(nil, nil)

		// Create test HTTP request
		req := httptest.NewRequest("DELETE", "/user/playlists/"+playlistID, nil)
		resp := httptest.NewRecorder()

		// Serve the request
		r.ServeHTTP(resp, req)

		// Assert
		assert.Equal(t, http.StatusNotFound, resp.Code, "Should return Not Found status")

		// Parse the response
		var response map[string]string
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		require.NoError(t, err, "Should parse response JSON")

		assert.Contains(t, response["error"], "playlist not found", "Should have not found error")

		// Verify mock expectations
		mockSpotifySongRepo.AssertExpectations(t)
	})
}
