package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/Emeruem-Kennedy1/ghopper/config"
	"github.com/Emeruem-Kennedy1/ghopper/internal/auth"
	"github.com/Emeruem-Kennedy1/ghopper/internal/models"
	"github.com/Emeruem-Kennedy1/ghopper/internal/repository"
	"github.com/Emeruem-Kennedy1/ghopper/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zmb3/spotify"
	"go.uber.org/zap"
)

// Setup JWT key for testing - we need to make sure this is used by the tests
func init() {
	// Set a fixed JWT key for testing
	auth.SetJWTKey([]byte("test-jwt-key"))
}

func mockCreateOrUpdateUserFromSpotifyData() func() {
	originalFunc := auth.CreateOrUpdateUserFromSpotifyDataFunc

	// Replace with our mock version
	auth.CreateOrUpdateUserFromSpotifyDataFunc = func(userRepo repository.UserRepositoryInterface, spotifyUser spotify.PrivateUser) (*models.User, string, error) {
		user := &models.User{
			ID:          spotifyUser.ID,
			DisplayName: spotifyUser.DisplayName,
			Email:       spotifyUser.Email,
			Country:     spotifyUser.Country,
			SpotifyURI:  string(spotifyUser.URI),
		}

		if len(spotifyUser.Images) > 0 {
			user.ProfileImage = spotifyUser.Images[0].URL
		}

		// Create a test token that will always be the same
		token := "test-jwt-token-for-" + user.ID

		return user, token, nil
	}

	// Return a function to restore the original
	return func() {
		auth.CreateOrUpdateUserFromSpotifyDataFunc = originalFunc
	}
}

// Setup function to create a test environment
func setupTest() (*gin.Engine, *MockSpotifyAuth, *MockUserRepository, services.ClientManagerInterface, *config.Config) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a new Gin router
	r := gin.New()

	// Create mocks
	mockSpotifyAuth := new(MockSpotifyAuth)
	mockUserRepo := new(MockUserRepository)

	// Create a client manager
	clientManager := services.NewClientManager()

	// Create a basic config
	cfg := &config.Config{
		FrontendURL: "http://localhost:3000",
	}

	return r, mockSpotifyAuth, mockUserRepo, clientManager, cfg
}

func TestSpotifyLogin(t *testing.T) {
	// Setup
	r, mockSpotifyAuth, _, _, _ := setupTest()

	// Configure mock
	expectedURL := "https://accounts.spotify.com/authorize?some=params"
	mockSpotifyAuth.On("AuthURL").Return(expectedURL)

	// Add the handler to router
	r.GET("/login", SpotifyLogin(mockSpotifyAuth))

	// Create a test HTTP request
	req := httptest.NewRequest("GET", "/login", nil)
	resp := httptest.NewRecorder()

	// Serve the request
	r.ServeHTTP(resp, req)

	// Assert expectations
	mockSpotifyAuth.AssertExpectations(t)
	assert.Equal(t, http.StatusTemporaryRedirect, resp.Code, "Should return temporary redirect status")
	assert.Equal(t, expectedURL, resp.Header().Get("Location"), "Should redirect to Spotify auth URL")
}

func TestSpotifyCallback(t *testing.T) {
	// Initialize the zap logger for testing
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)

	t.Run("Successful_Callback", func(t *testing.T) {
		// Mock the CreateOrUpdateUserFromSpotifyData function
		restore := mockCreateOrUpdateUserFromSpotifyData()
		defer restore()

		// Setup
		r, mockSpotifyAuth, mockUserRepo, clientManager, cfg := setupTest()

		// Create mocks for Spotify client and user
		spotifyClient := &spotify.Client{}
		spotifyUser := &spotify.PrivateUser{
			User: spotify.User{
				ID:          "spotify-user-id",
				DisplayName: "Spotify User",
				Images:      []spotify.Image{{URL: "https://example.com/profile.jpg"}},
				URI:         "spotify:user:spotify-user-id",
			},
			Email:   "user@example.com",
			Country: "US",
		}

		// Configure mocks
		mockSpotifyAuth.On("CallBack", mock.Anything).Return(spotifyClient, nil)
		mockSpotifyAuth.On("GetUserInfo", spotifyClient).Return(spotifyUser, nil)

		// Add the handler to router
		r.GET("/callback", SpotifyCallback(mockSpotifyAuth, mockUserRepo, clientManager, cfg))

		// Create a test HTTP request
		req := httptest.NewRequest("GET", "/callback", nil)
		resp := httptest.NewRecorder()

		// Serve the request
		r.ServeHTTP(resp, req)

		// Assert expectations
		mockSpotifyAuth.AssertExpectations(t)

		// Verify HTTP status
		assert.Equal(t, http.StatusTemporaryRedirect, resp.Code, "Should return temporary redirect status")

		// Verify we have a redirect URL containing the data
		redirectURL := resp.Header().Get("Location")
		assert.Contains(t, redirectURL, cfg.FrontendURL, "Redirect URL should contain frontend URL")
		assert.Contains(t, redirectURL, "data=", "Redirect URL should contain encoded data")

		// Parse the URL to get the data parameter
		parsedURL, err := url.Parse(redirectURL)
		require.NoError(t, err, "Should parse redirect URL")

		// Get the data parameter
		dataParam := parsedURL.Query().Get("data")
		require.NotEmpty(t, dataParam, "Data parameter should not be empty")

		// Decode the data
		decodedData, err := base64.URLEncoding.DecodeString(dataParam)
		require.NoError(t, err, "Should decode base64 data")

		// Parse the JSON
		var responseData map[string]interface{}
		err = json.Unmarshal(decodedData, &responseData)
		require.NoError(t, err, "Should parse JSON data")

		// Verify the JSON data
		assert.Equal(t, "Successfully authenticated", responseData["message"], "Message should indicate success")
		assert.Equal(t, "test-jwt-token-for-spotify-user-id", responseData["token"], "Should include a JWT token")
		assert.NotNil(t, responseData["user"], "Should include user data")
	})

	t.Run("Callback_Failure_Invalid_Client", func(t *testing.T) {
		// Setup
		r, mockSpotifyAuth, mockUserRepo, clientManager, cfg := setupTest()

		// Configure mock to return an error
		expectedError := fmt.Errorf("callback error")
		mockSpotifyAuth.On("CallBack", mock.Anything).Return((*spotify.Client)(nil), expectedError)

		// Add the handler to router
		r.GET("/callback/error", SpotifyCallback(mockSpotifyAuth, mockUserRepo, clientManager, cfg))

		// Create a test HTTP request
		req := httptest.NewRequest("GET", "/callback/error", nil)
		resp := httptest.NewRecorder()

		// Serve the request
		r.ServeHTTP(resp, req)

		// Assert expectations
		mockSpotifyAuth.AssertExpectations(t)
		assert.Equal(t, http.StatusBadRequest, resp.Code, "Should return bad request status")

		// Parse the response body
		var responseBody map[string]string
		err := json.Unmarshal(resp.Body.Bytes(), &responseBody)
		require.NoError(t, err, "Should parse JSON response")

		// Verify the error message
		assert.Equal(t, expectedError.Error(), responseBody["error"], "Response should contain error message")
	})
}
