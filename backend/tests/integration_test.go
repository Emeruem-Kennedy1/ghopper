package tests

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Emeruem-Kennedy1/ghopper/config"
	"github.com/Emeruem-Kennedy1/ghopper/internal/auth"
	"github.com/Emeruem-Kennedy1/ghopper/internal/handlers"
	"github.com/Emeruem-Kennedy1/ghopper/internal/middleware"
	"github.com/Emeruem-Kennedy1/ghopper/internal/models"
	"github.com/Emeruem-Kennedy1/ghopper/internal/repository"
	"github.com/Emeruem-Kennedy1/ghopper/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestSetup contains all the dependencies required for integration tests
type TestSetup struct {
	Router          *gin.Engine
	DB              *gorm.DB
	UserRepo        *repository.UserRepository
	SpotifySongRepo *repository.SpotifySongRepository
	ClientManager   services.ClientManagerInterface
	SpotifyService  services.SpotifyServiceInterface
	Config          *config.Config
	Logger          *zap.Logger
}

// setupIntegrationTest creates a test environment with real repositories but mock database
func setupIntegrationTest(t *testing.T) *TestSetup {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create in-memory SQLite database
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err, "Failed to create test database")

	// Drop any existing tables first
	err = db.Migrator().DropTable(&models.User{}, &models.Song{}, &models.Playlist{})
	require.NoError(t, err, "Failed to drop existing tables")

	// Migrate the schema
	err = db.AutoMigrate(&models.User{}, &models.Song{}, &models.Playlist{})
	require.NoError(t, err, "Failed to migrate schema")

	// Create logger
	logger, _ := zap.NewDevelopment()

	// Create repositories
	userRepo := repository.NewUserRepository(db)
	spotifySongRepo := repository.NewSpotifySongRepository(db)

	// Create client manager and service
	clientManager := services.NewClientManager()
	spotifyService := &handlers.MockSpotifyService{}

	// Create a test config
	cfg := &config.Config{
		Port:                "8888",
		Env:                 "test",
		SpotifyClientID:     "test-client-id",
		SpotifyClientSecret: "test-client-secret",
		SpotifyRedirectURI:  "http://localhost:8888/callback",
		FrontendURL:         "http://localhost:3000",
		JWTSecret:           "test-jwt-secret",
	}

	// Create router
	r := gin.New()

	// Create test setup
	setup := &TestSetup{
		Router:          r,
		DB:              db,
		UserRepo:        userRepo,
		SpotifySongRepo: spotifySongRepo,
		ClientManager:   clientManager,
		SpotifyService:  spotifyService,
		Config:          cfg,
		Logger:          logger,
	}

	// Configure routes - similar to how it's done in server.go but simplified for tests
	setupTestRoutes(setup)

	return setup
}

// setupTestRoutes configures routes for integration tests
func setupTestRoutes(setup *TestSetup) {
	// Public routes
	setup.Router.GET("/ping", handlers.Ping())

	// Create a protected group with the auth middleware
	protected := setup.Router.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		// User routes
		protected.GET("/user", handlers.GetUser(setup.UserRepo))

		// Playlist routes
		protected.GET("/user/playlists", handlers.GetUserPlaylists(setup.SpotifySongRepo, setup.SpotifyService))
		protected.DELETE("/user/playlists/:playlistID", handlers.DeletePlaylist(setup.SpotifyService, setup.SpotifySongRepo))

		// User account management
		protected.DELETE("/user/account", handlers.DeleteUserAccount(setup.UserRepo, setup.SpotifySongRepo, setup.ClientManager))
	}
}

// createTestUser creates a test user and returns a valid token
func createTestUser(t *testing.T, userRepo *repository.UserRepository) (string, *models.User) {
	user := &models.User{
		ID:           "test-user-id" + t.Name(),
		DisplayName:  "Test User",
		Email:        "test" + t.Name() + "@example.com",
		SpotifyURI:   "spotify:user:test" + t.Name(),
		Country:      "US",
		ProfileImage: "https://example.com/profile.jpg",
	}

	err := userRepo.Create(user)
	require.NoError(t, err, "Failed to create test user")

	token, err := auth.GenerateToken(user)
	require.NoError(t, err, "Failed to generate token for test user")

	return token, user
}

// createTestPlaylist creates a test playlist for the given user
func createTestPlaylist(t *testing.T, repo *repository.SpotifySongRepository, userID string) *models.Playlist {
	playlist := &models.Playlist{
		ID:          "test-playlist-id",
		UserID:      userID,
		Name:        "Test Playlist",
		Description: "Test playlist description",
		URL:         "https://open.spotify.com/playlist/test-playlist-id",
		Image:       "https://example.com/playlist.jpg",
	}

	err := repo.SavePlaylist(playlist)
	require.NoError(t, err, "Failed to create test playlist")

	return playlist
}

func TestPingEndpoint(t *testing.T) {
	// Setup
	setup := setupIntegrationTest(t)

	// Create a test HTTP request
	req := httptest.NewRequest("GET", "/ping", nil)
	resp := httptest.NewRecorder()

	// Serve the request
	setup.Router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusOK, resp.Code, "Should return OK status")

	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	require.NoError(t, err, "Should return valid JSON")

	assert.Equal(t, "pong", response["message"], "Response should contain 'pong' message")
}

func TestGetUserEndpoint(t *testing.T) {
	// Setup
	setup := setupIntegrationTest(t)

	// Create a test user and get a token
	token, user := createTestUser(t, setup.UserRepo)

	// Create a test HTTP request
	req := httptest.NewRequest("GET", "/api/user", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp := httptest.NewRecorder()

	// Serve the request
	setup.Router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusOK, resp.Code, "Should return OK status")

	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	require.NoError(t, err, "Should return valid JSON")

	// Check user data in response
	userData, ok := response["user"].(map[string]interface{})
	require.True(t, ok, "Response should contain user object")
	assert.Equal(t, user.ID, userData["id"], "User ID should match")
	assert.Equal(t, user.DisplayName, userData["display_name"], "Display name should match")
	assert.Equal(t, user.Email, userData["email"], "Email should match")
}

func TestGetUserPlaylists(t *testing.T) {
	// Setup
	setup := setupIntegrationTest(t)

	// Create a test user and get a token
	token, user := createTestUser(t, setup.UserRepo)

	// Create test playlists
	playlist1 := createTestPlaylist(t, setup.SpotifySongRepo, user.ID)
	playlist2 := &models.Playlist{
		ID:          "test-playlist-id-2",
		UserID:      user.ID,
		Name:        "Test Playlist 2",
		Description: "Another test playlist",
		URL:         "https://open.spotify.com/playlist/test-playlist-id-2",
		Image:       "https://example.com/playlist2.jpg",
	}
	err := setup.SpotifySongRepo.SavePlaylist(playlist2)
	require.NoError(t, err, "Failed to create second test playlist")

	// Create a test HTTP request
	req := httptest.NewRequest("GET", "/api/user/playlists", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp := httptest.NewRecorder()

	// Serve the request
	setup.Router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusOK, resp.Code, "Should return OK status")

	var response map[string]interface{}
	err = json.Unmarshal(resp.Body.Bytes(), &response)
	require.NoError(t, err, "Should return valid JSON")

	// Check playlists in response
	playlists, ok := response["playlists"].([]interface{})
	require.True(t, ok, "Response should contain playlists array")
	assert.Len(t, playlists, 2, "Should return 2 playlists")

	// Verify first playlist details
	playlist1Data := playlists[0].(map[string]interface{})
	assert.Equal(t, playlist1.ID, playlist1Data["id"], "Playlist ID should match")
	assert.Equal(t, playlist1.Name, playlist1Data["name"], "Playlist name should match")
}

func TestDeleteUserAccount(t *testing.T) {
	// Setup
	setup := setupIntegrationTest(t)

	// Create a test user and get a token
	token, user := createTestUser(t, setup.UserRepo)

	// Create a test playlist for the user
	createTestPlaylist(t, setup.SpotifySongRepo, user.ID)

	// Create a test HTTP request
	req := httptest.NewRequest("DELETE", "/api/user/account", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp := httptest.NewRecorder()

	// Serve the request
	setup.Router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusOK, resp.Code, "Should return OK status")

	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	require.NoError(t, err, "Should return valid JSON")

	assert.Contains(t, response["message"], "successfully", "Response should indicate success")

	// Verify user was deleted
	_, err = setup.UserRepo.GetByID(user.ID)
	assert.Error(t, err, "User should be deleted from database")

	// Verify playlists were deleted
	playlists, err := setup.SpotifySongRepo.GetUserPlaylists(user.ID)
	require.NoError(t, err, "Should not error when checking playlists")
	assert.Empty(t, playlists, "All playlists should be deleted")
}

func TestDeletePlaylist(t *testing.T) {
	// Setup
	setup := setupIntegrationTest(t)

	// Create a test user and get a token
	token, user := createTestUser(t, setup.UserRepo)

	// Create a test playlist for the user
	playlistName := "Test Playlist " + t.Name()
	playlist := &models.Playlist{
		ID:          "test-playlist-" + t.Name(),
		UserID:      user.ID,
		Name:        playlistName,
		Description: "Test playlist description",
		URL:         "https://open.spotify.com/playlist/test-playlist-" + t.Name(),
		Image:       "https://example.com/playlist.jpg",
	}

	err := setup.SpotifySongRepo.SavePlaylist(playlist)
	require.NoError(t, err, "Failed to create test playlist")

	// Configure the mock Spotify service
	mockSpotifyService := setup.SpotifyService.(*handlers.MockSpotifyService)
	mockSpotifyService.On("DeletePlaylist", user.ID, playlist.ID).Return(nil)

	// Create a test HTTP request to delete the playlist
	req := httptest.NewRequest("DELETE", "/api/user/playlists/"+playlist.ID, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp := httptest.NewRecorder()

	// Serve the request
	setup.Router.ServeHTTP(resp, req)

	// Print response body for debugging
	bodyBytes, _ := io.ReadAll(resp.Body)
	t.Logf("Response Status: %d", resp.Code)
	t.Logf("Response Body: %s", string(bodyBytes))

	// Assert
	assert.Equal(t, http.StatusOK, resp.Code, "Should return OK status")

	// Verify mock expectations
	mockSpotifyService.AssertExpectations(t)

	// Verify playlist was deleted from the database
	foundPlaylist, err := setup.SpotifySongRepo.FindPlaylistByIDAndUser(playlist.ID, user.ID)
	assert.Nil(t, foundPlaylist, "Playlist should be deleted from database")
	assert.NoError(t, err, "Should not error when checking for deleted playlist")
}
