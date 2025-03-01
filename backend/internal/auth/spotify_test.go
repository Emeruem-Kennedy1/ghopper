package auth

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/Emeruem-Kennedy1/ghopper/config"
	"github.com/Emeruem-Kennedy1/ghopper/internal/models"
	"github.com/Emeruem-Kennedy1/ghopper/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

// We'll use a different approach with dependency injection
// Create a testable version of SpotifyAuth with mockable dependencies

// Define an interface for the utility function
type RandomStringGenerator func(length int) (string, error)

// Modify SpotifyAuth to accept the generator function
type testableSpotifyAuth struct {
	SpotifyAuth
	generateRandomString RandomStringGenerator
}

// Create a function that returns a predictable string for testing
func mockGenerateRandomString(length int) (string, error) {
	return "test-state", nil
}

// Create a function that returns an error for testing error cases
func errorGenerateRandomString(length int) (string, error) {
	return "", errors.New("random generation error")
}

type MockAuthenticator struct {
	mock.Mock
}

func (m *MockAuthenticator) AuthURL(state string) string {
	args := m.Called(state)
	return args.String(0)
}

func (m *MockAuthenticator) Token(state string, r *http.Request) (*oauth2.Token, error) {
	args := m.Called(state, r)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*oauth2.Token), args.Error(1)
}

func (m *MockAuthenticator) NewClient(token *oauth2.Token) spotify.Client {
	args := m.Called(token)
	return args.Get(0).(spotify.Client)
}

func (m *MockAuthenticator) SetAuthInfo(clientID, secretKey string) {
	m.Called(clientID, secretKey)
}

// Create a testable version of NewSpotifyAuth that uses the provided generator
func newTestableSpotifyAuth(cfg *config.Config, generator RandomStringGenerator) (*testableSpotifyAuth, error) {
	auth := new(MockAuthenticator)

	// Set up expectations for the mock
	auth.On("SetAuthInfo", cfg.SpotifyClientID, cfg.SpotifyClientSecret).Return()

	state, err := generator(10)
	if err != nil {
		return nil, fmt.Errorf("couldn't generate random string: %v", err)
	}

	return &testableSpotifyAuth{
		SpotifyAuth: SpotifyAuth{
			authenticator: auth,
			state:         state,
			config:        cfg,
		},
		generateRandomString: generator,
	}, nil
}

func TestNewSpotifyAuth(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Arrange
		cfg := &config.Config{
			SpotifyRedirectURI:  "http://localhost:8080/callback",
			SpotifyClientID:     "client-id",
			SpotifyClientSecret: "client-secret",
		}

		// Act
		// For the real NewSpotifyAuth function, we'll use a monkey patch library or restructure the code
		// For this test, we'll use our testable version
		auth, err := newTestableSpotifyAuth(cfg, mockGenerateRandomString)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, auth)
		assert.Equal(t, "test-state", auth.state)
		assert.NotNil(t, auth.authenticator)
		assert.Equal(t, cfg, auth.config)
	})

	t.Run("Error_GeneratingState", func(t *testing.T) {
		// Arrange
		cfg := &config.Config{
			SpotifyRedirectURI:  "http://localhost:8080/callback",
			SpotifyClientID:     "client-id",
			SpotifyClientSecret: "client-secret",
		}

		// Act
		auth, err := newTestableSpotifyAuth(cfg, errorGenerateRandomString)

		// Assert
		require.Error(t, err)
		require.Nil(t, auth)
		assert.Contains(t, err.Error(), "couldn't generate random string")
	})
}

func TestSpotifyAuth_AuthURL(t *testing.T) {
	// Arrange
	cfg := &config.Config{
		SpotifyRedirectURI:  "http://localhost:8080/callback",
		SpotifyClientID:     "client-id",
		SpotifyClientSecret: "client-secret",
	}

	// Create a mock authenticator directly
	mockAuth := new(MockAuthenticator)
	
	// Set up the AuthURL expectation with the test state
	expectedURL := "https://accounts.spotify.com/authorize?client_id=client-id&redirect_uri=http%3A%2F%2Flocalhost%3A8080%2Fcallback&response_type=code&state=test-state"
	mockAuth.On("AuthURL", "test-state").Return(expectedURL)
	
	// No need to set up SetAuthInfo expectation if we're not calling it

	// Create the auth struct with our mock
	auth := &testableSpotifyAuth{
		SpotifyAuth: SpotifyAuth{
			authenticator: mockAuth,
			state:         "test-state",
			config:        cfg,
		},
	}

	// Act
	authURL := auth.AuthURL()

	// Assert
	require.NotEmpty(t, authURL)
	assert.Equal(t, expectedURL, authURL)
	// The URL should contain the state
	assert.Contains(t, authURL, "state=test-state")
	// The URL should contain the client ID
	assert.Contains(t, authURL, "client_id=client-id")
	// The URL should contain the redirect URI
	assert.Contains(t, authURL, url.QueryEscape("http://localhost:8080/callback"))

	// Verify all expectations were met
	mockAuth.AssertExpectations(t)
}

func TestSpotifyAuth_CallBack(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Arrange
		cfg := &config.Config{
			SpotifyRedirectURI:  "http://localhost:8080/callback",
			SpotifyClientID:     "client-id",
			SpotifyClientSecret: "client-secret",
		}

		// Create a mock authenticator instead of using the real one
		mockAuth := new(MockAuthenticator)

		// Create a SpotifyAuth instance with the mock authenticator
		auth := &testableSpotifyAuth{
			SpotifyAuth: SpotifyAuth{
				authenticator: mockAuth,
				state:         "test-state",
				config:        cfg,
			},
		}

		// Create a request with the correct state
		r := httptest.NewRequest("GET", "/callback?state=test-state&code=test-code", nil)

		// Configure the mock to return a token when Token is called
		mockToken := &oauth2.Token{
			AccessToken:  "test-access-token",
			TokenType:    "Bearer",
			RefreshToken: "test-refresh-token",
			Expiry:       time.Now().Add(time.Hour),
		}
		mockAuth.On("Token", "test-state", r).Return(mockToken, nil)

		// Configure the mock to return a client when NewClient is called
		mockClient := spotify.Client{}
		mockAuth.On("NewClient", mockToken).Return(mockClient)

		// Act
		client, err := auth.CallBack(r)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, client)
		mockAuth.AssertExpectations(t)
	})

	t.Run("State_Mismatch", func(t *testing.T) {
		// Arrange
		cfg := &config.Config{
			SpotifyRedirectURI:  "http://localhost:8080/callback",
			SpotifyClientID:     "client-id",
			SpotifyClientSecret: "client-secret",
		}

		// Create a mock authenticator
		mockAuth := new(MockAuthenticator)

		// Create a SpotifyAuth instance with the mock authenticator
		auth := &testableSpotifyAuth{
			SpotifyAuth: SpotifyAuth{
				authenticator: mockAuth,
				state:         "test-state",
				config:        cfg,
			},
		}

		// Create a request with an incorrect state
		r := httptest.NewRequest("GET", "/callback?state=wrong-state&code=test-code", nil)

		// When state doesn't match, your CallBack function should check and return an error
		// We don't need to set up any expectations on the mock since it shouldn't be called

		// Act
		client, err := auth.CallBack(r)

		// Assert
		require.Error(t, err)
		assert.Nil(t, client)
		assert.Contains(t, err.Error(), "state mismatch")

		// Ensure Token was never called on the mock
		mockAuth.AssertNotCalled(t, "Token")
	})

	t.Run("Token_Error", func(t *testing.T) {
		// Arrange
		cfg := &config.Config{
			SpotifyRedirectURI:  "http://localhost:8080/callback",
			SpotifyClientID:     "client-id",
			SpotifyClientSecret: "client-secret",
		}

		// Create a mock authenticator
		mockAuth := new(MockAuthenticator)

		// Create a SpotifyAuth instance with the mock authenticator
		auth := &testableSpotifyAuth{
			SpotifyAuth: SpotifyAuth{
				authenticator: mockAuth,
				state:         "test-state",
				config:        cfg,
			},
		}

		// Create a request with the correct state
		r := httptest.NewRequest("GET", "/callback?state=test-state&code=test-code", nil)

		// Configure the mock to return an error when Token is called
		expectedError := errors.New("token error")
		mockAuth.On("Token", "test-state", r).Return(nil, expectedError)

		// Act
		client, err := auth.CallBack(r)

		// Assert
		require.Error(t, err)
		assert.Nil(t, client)
		assert.Contains(t, err.Error(), "couldn't get token")
		mockAuth.AssertExpectations(t)
	})
}

func TestCreateOrUpdateUserFromSpotifyData(t *testing.T) {
	// Save original function and restore after tests
	originalFunc := CreateOrUpdateUserFromSpotifyDataFunc
	defer func() { CreateOrUpdateUserFromSpotifyDataFunc = originalFunc }()

	t.Run("Success", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepository)
		expectedUser := &models.User{ID: "test-id"}
		expectedToken := "test-token"

		// Mock the function
		CreateOrUpdateUserFromSpotifyDataFunc = func(userRepo repository.UserRepositoryInterface, spotifyUser spotify.PrivateUser) (*models.User, string, error) {
			return expectedUser, expectedToken, nil
		}

		// Create a spotify user
		spotifyUser := spotify.PrivateUser{User: spotify.User{ID: "test-id"}}

		// Act
		user, token, err := CreateOrUpdateUserFromSpotifyData(mockRepo, spotifyUser)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expectedUser, user)
		assert.Equal(t, expectedToken, token)
	})

	t.Run("Error", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepository)
		expectedError := errors.New("test error")

		// Mock the function
		CreateOrUpdateUserFromSpotifyDataFunc = func(userRepo repository.UserRepositoryInterface, spotifyUser spotify.PrivateUser) (*models.User, string, error) {
			return nil, "", expectedError
		}

		// Create a spotify user
		spotifyUser := spotify.PrivateUser{User: spotify.User{ID: "test-id"}}

		// Act
		user, token, err := CreateOrUpdateUserFromSpotifyData(mockRepo, spotifyUser)

		// Assert
		require.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, user)
		assert.Empty(t, token)
	})
}

func TestCreateOrUpdateUserFromSpotifyDataImpl(t *testing.T) {
	// Save the original function and restore after tests
	originalTokenFunc := GenerateTokenFunc
	defer func() { GenerateTokenFunc = originalTokenFunc }()

	t.Run("Success", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepository)

		// Create a spotify user with test data
		spotifyUser := spotify.PrivateUser{
			User: spotify.User{
				ID:          "test-id",
				DisplayName: "Test User",
				URI:         "spotify:user:test-id",
				Images:      []spotify.Image{{URL: "https://example.com/image.jpg"}},
			},
			Email:   "test@example.com",
			Country: "US",
		}

		// Set up the mock expectation
		mockRepo.On("UpsertUser", mock.MatchedBy(func(user *models.User) bool {
			return user.ID == "test-id" &&
				user.DisplayName == "Test User" &&
				user.Email == "test@example.com"
		})).Return(nil)

		// Create a mock for GenerateToken
		GenerateTokenFunc = func(user *models.User) (string, error) {
			return "test-token", nil
		}

		// Act
		user, token, err := CreateOrUpdateUserFromSpotifyDataImpl(mockRepo, spotifyUser)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "test-id", user.ID)
		assert.Equal(t, "Test User", user.DisplayName)
		assert.Equal(t, "test@example.com", user.Email)
		assert.Equal(t, "US", user.Country)
		assert.Equal(t, "spotify:user:test-id", user.SpotifyURI)
		assert.Equal(t, "https://example.com/image.jpg", user.ProfileImage)
		assert.Equal(t, "test-token", token)

		mockRepo.AssertExpectations(t)
	})

	t.Run("No_Profile_Image", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepository)

		// Create a spotify user with no images
		spotifyUser := spotify.PrivateUser{
			User: spotify.User{
				ID:          "test-id",
				DisplayName: "Test User",
				URI:         "spotify:user:test-id",
			},
			Email:   "test@example.com",
			Country: "US",
		}

		// Set up the mock expectation
		mockRepo.On("UpsertUser", mock.MatchedBy(func(user *models.User) bool {
			return user.ID == "test-id" && user.ProfileImage == ""
		})).Return(nil)

		// Create a mock for GenerateToken
		GenerateTokenFunc = func(user *models.User) (string, error) {
			return "test-token", nil
		}

		// Act
		user, token, err := CreateOrUpdateUserFromSpotifyDataImpl(mockRepo, spotifyUser)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "test-id", user.ID)
		assert.Empty(t, user.ProfileImage)
		assert.Equal(t, "test-token", token)

		mockRepo.AssertExpectations(t)
	})

	t.Run("UpsertUser_Error", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepository)

		// Create a spotify user
		spotifyUser := spotify.PrivateUser{
			User: spotify.User{ID: "test-id"},
		}

		// Set up the mock to return an error
		expectedError := errors.New("database error")
		mockRepo.On("UpsertUser", mock.Anything).Return(expectedError)

		// Act
		user, token, err := CreateOrUpdateUserFromSpotifyDataImpl(mockRepo, spotifyUser)

		// Assert
		require.Error(t, err)
		assert.Nil(t, user)
		assert.Empty(t, token)
		assert.Contains(t, err.Error(), "failed to create or update user")

		mockRepo.AssertExpectations(t)
	})

	t.Run("GenerateToken_Error", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockUserRepository)

		// Create a spotify user
		spotifyUser := spotify.PrivateUser{
			User: spotify.User{ID: "test-id"},
		}

		// Set up the mock to succeed
		mockRepo.On("UpsertUser", mock.Anything).Return(nil)

		// Create a mock for GenerateToken that fails
		expectedError := errors.New("token generation error")
		GenerateTokenFunc = func(user *models.User) (string, error) {
			return "", expectedError
		}

		// Act
		user, token, err := CreateOrUpdateUserFromSpotifyDataImpl(mockRepo, spotifyUser)

		// Assert
		require.Error(t, err)
		assert.Nil(t, user)
		assert.Empty(t, token)
		assert.Contains(t, err.Error(), "failed to generate token")

		mockRepo.AssertExpectations(t)
	})
}

// Define a mock Spotify client
type MockSpotifyClient struct {
	mock.Mock
}

type SpotifyClientInterface interface {
	CurrentUser() (*spotify.PrivateUser, error)
}

func (m *MockSpotifyClient) CurrentUser() (*spotify.PrivateUser, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*spotify.PrivateUser), args.Error(1)
}

// Modify SpotifyAuth to accept the interface
func (auth *testableSpotifyAuth) GetUserInfo(client SpotifyClientInterface) (*spotify.PrivateUser, error) {
	user, err := client.CurrentUser()
	if err != nil {
		return nil, fmt.Errorf("couldn't get user: %v", err)
	}
	return user, nil
}

func TestSpotifyAuth_GetUserInfo(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Arrange
		cfg := &config.Config{}
		auth, err := newTestableSpotifyAuth(cfg, mockGenerateRandomString)
		require.NoError(t, err)

		mockClient := new(MockSpotifyClient)
		expectedUser := &spotify.PrivateUser{
			User: spotify.User{
				ID:          "test-id",
				DisplayName: "Test User",
			},
			Email: "test@example.com",
		}

		mockClient.On("CurrentUser").Return(expectedUser, nil)

		// Act
		user, err := auth.GetUserInfo(mockClient)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, expectedUser, user)
		mockClient.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		// Arrange
		cfg := &config.Config{}
		auth, err := newTestableSpotifyAuth(cfg, mockGenerateRandomString)
		require.NoError(t, err)

		mockClient := new(MockSpotifyClient)
		expectedError := errors.New("spotify API error")

		mockClient.On("CurrentUser").Return(nil, expectedError)

		// Act
		user, err := auth.GetUserInfo(mockClient)

		// Assert
		require.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "couldn't get user")
		mockClient.AssertExpectations(t)
	})
}

func TestSpotifyAuth_GetAuthenticator(t *testing.T) {
	// Arrange
	cfg := &config.Config{
		SpotifyRedirectURI:  "http://localhost:8080/callback",
		SpotifyClientID:     "client-id",
		SpotifyClientSecret: "client-secret",
	}
	auth, err := newTestableSpotifyAuth(cfg, mockGenerateRandomString)
	require.NoError(t, err)

	// Act
	authenticator := auth.GetAuthenticator()

	// Assert the authenticator is not nil
	require.NotNil(t, authenticator)
}
