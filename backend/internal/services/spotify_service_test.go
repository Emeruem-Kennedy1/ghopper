package services

import (
	"errors"
	"testing"

	"github.com/Emeruem-Kennedy1/ghopper/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zmb3/spotify"
)

func TestGetSongURL(t *testing.T) {
	t.Run("Song_Found_In_Database", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockSpotifySongRepository)
		clientManager := NewClientManager()
		service := NewSpotifyService(clientManager, mockRepo)

		// Test data
		userID := "test-user"
		songName := "Test Song"
		artistName := "Test Artist"
		expectedURL := "https://open.spotify.com/track/1234567890"

		// Setup repository mock to return the song
		mockRepo.On("FindSongByNameAndArtist", songName, artistName).Return(
			&models.Song{
				ID:         "1234567890",
				Name:       songName,
				Artist:     artistName,
				SpotifyURL: expectedURL,
			}, nil)

		// Act
		url, err := service.GetSongURL(userID, songName, artistName)

		// Assert
		require.NoError(t, err, "GetSongURL should not return error when song is found in database")
		assert.Equal(t, expectedURL, url, "Returned URL should match expected URL")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Song_Not_Found_No_Client", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockSpotifySongRepository)
		clientManager := NewClientManager()
		service := NewSpotifyService(clientManager, mockRepo)

		// Test data
		userID := "test-user"
		songName := "Test Song"
		artistName := "Test Artist"

		// Setup repository mock to return nil (song not found)
		mockRepo.On("FindSongByNameAndArtist", songName, artistName).Return(nil, nil)

		// Act
		url, err := service.GetSongURL(userID, songName, artistName)

		// Assert
		assert.Error(t, err, "GetSongURL should return error when song is not found and no client exists")
		assert.Empty(t, url, "URL should be empty when error occurs")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Song_Not_Found_With_Client", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockSpotifySongRepository)
		mockClient := new(MockSpotifyClient)
		clientManager := NewClientManager()
		service := NewSpotifyService(clientManager, mockRepo)

		// Test data
		userID := "test-user"
		songName := "Test Song"
		artistName := "Test Artist"
		trackID := "track123"
		expectedURL := "https://open.spotify.com/track/track123"

		// Store mock client
		clientManager.StoreClient(userID, mockClient)

		// Setup repository mock to return nil (song not found)
		mockRepo.On("FindSongByNameAndArtist", songName, artistName).Return(nil, nil)

		// Setup Spotify client mock
		searchQuery := "track:Test Song artist:Test Artist"
		searchResult := &spotify.SearchResult{
			Tracks: &spotify.FullTrackPage{
				Tracks: []spotify.FullTrack{
					{
						SimpleTrack: spotify.SimpleTrack{
							ID: spotify.ID(trackID),
						},
						Album: spotify.SimpleAlbum{
							Images: []spotify.Image{
								{URL: "https://example.com/image.jpg"},
							},
						},
					},
				},
			},
		}

		mockClient.On("Search", searchQuery, int(spotify.SearchTypeTrack)).Return(searchResult, nil)

		// Setup repository mock to save the song
		mockRepo.On("SaveSong", mock.MatchedBy(func(song *models.Song) bool {
			return song.ID == trackID && song.Name == songName && song.Artist == artistName
		})).Return(nil)

		// Act
		url, err := service.GetSongURL(userID, songName, artistName)

		// Assert
		require.NoError(t, err, "GetSongURL should not return error when song is found via Spotify API")
		assert.Equal(t, expectedURL, url, "Returned URL should match expected URL")
		mockRepo.AssertExpectations(t)
		mockClient.AssertExpectations(t)
	})
}

func TestCreatePlaylistFromSongs(t *testing.T) {
	t.Run("Create_New_Playlist", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockSpotifySongRepository)
		mockClient := new(MockSpotifyClient)
		clientManager := NewClientManager()
		service := NewSpotifyService(clientManager, mockRepo)

		// Test data
		userID := "test-user"
		playlistName := "Test Playlist"
		playlistDesc := "Test Description"
		songIDs := []spotify.ID{"song1", "song2", "song3"}

		// Store mock client
		clientManager.StoreClient(userID, mockClient)

		// Setup client mock for CurrentUser
		spotifyUser := &spotify.PrivateUser{User: spotify.User{ID: "spotify-user-id"}}
		mockClient.On("CurrentUser").Return(spotifyUser, nil)

		// Setup repository mock for FindPlaylistByNameAndUser
		mockRepo.On("FindPlaylistByNameAndUser", playlistName, userID).Return(nil, nil)

		// Setup client mock for CreatePlaylistForUser
		expectedURL := "https://open.spotify.com/playlist/playlist123"
		playlist := &spotify.FullPlaylist{
			SimplePlaylist: spotify.SimplePlaylist{
				ID: "playlist123",
				ExternalURLs: map[string]string{
					"spotify": expectedURL,
				},
			},
		}
		mockClient.On("CreatePlaylistForUser", spotifyUser.ID, playlistName, playlistDesc, false).Return(playlist, nil)

		// Setup repository mock for SavePlaylist
		mockRepo.On("SavePlaylist", mock.MatchedBy(func(p *models.Playlist) bool {
			return p.ID == "playlist123" && p.Name == playlistName && p.URL == expectedURL
		})).Return(nil)

		// Setup client mock for GetPlaylistTracks
		mockClient.On("GetPlaylistTracks", spotify.ID("playlist123")).Return(&spotify.PlaylistTrackPage{
			Tracks: []spotify.PlaylistTrack{},
		}, nil)

		// Setup client mock for AddTracksToPlaylist
		mockClient.On("AddTracksToPlaylist", spotify.ID("playlist123"), songIDs).Return("snapshot123", nil)

		// Act
		url, err := service.CreatePlaylistFromSongs(userID, songIDs, playlistName, playlistDesc)

		// Assert
		require.NoError(t, err, "CreatePlaylistFromSongs should not return error")
		assert.Equal(t, expectedURL, url, "Returned URL should match expected URL")
		mockRepo.AssertExpectations(t)
		mockClient.AssertExpectations(t)
	})

	t.Run("Use_Existing_Playlist", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockSpotifySongRepository)
		mockClient := new(MockSpotifyClient)
		clientManager := NewClientManager()
		service := NewSpotifyService(clientManager, mockRepo)

		// Test data
		userID := "test-user"
		playlistName := "Test Playlist"
		playlistDesc := "Test Description"
		playlistID := "existing-playlist"
		songIDs := []spotify.ID{"song1", "song2", "song3"}
		expectedURL := "https://open.spotify.com/playlist/existing-playlist"

		// Store mock client
		clientManager.StoreClient(userID, mockClient)

		spotifyUser := &spotify.PrivateUser{User: spotify.User{ID: "spotify-user-id"}}
		mockClient.On("CurrentUser").Return(spotifyUser, nil)

		// Setup repository mock for FindPlaylistByNameAndUser to return existing playlist
		existingPlaylist := &models.Playlist{
			ID:          playlistID,
			UserID:      userID,
			Name:        playlistName,
			Description: playlistDesc,
			URL:         expectedURL,
		}
		mockRepo.On("FindPlaylistByNameAndUser", playlistName, userID).Return(existingPlaylist, nil)

		// Setup client mock for GetPlaylistTracks
		mockClient.On("GetPlaylistTracks", spotify.ID(playlistID)).Return(&spotify.PlaylistTrackPage{
			Tracks: []spotify.PlaylistTrack{},
		}, nil)

		// Setup client mock for AddTracksToPlaylist
		mockClient.On("AddTracksToPlaylist", spotify.ID(playlistID), songIDs).Return("snapshot123", nil)

		// Act
		url, err := service.CreatePlaylistFromSongs(userID, songIDs, playlistName, playlistDesc)

		// Assert
		require.NoError(t, err, "CreatePlaylistFromSongs should not return error for existing playlist")
		assert.Equal(t, expectedURL, url, "Returned URL should match existing playlist URL")
		mockRepo.AssertExpectations(t)
		mockClient.AssertExpectations(t)
	})

	t.Run("No_Client_Error", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockSpotifySongRepository)
		clientManager := NewClientManager()
		service := NewSpotifyService(clientManager, mockRepo)

		// Test data
		userID := "test-user"
		playlistName := "Test Playlist"
		playlistDesc := "Test Description"
		songIDs := []spotify.ID{"song1", "song2", "song3"}

		// Act - no client registered for this user
		url, err := service.CreatePlaylistFromSongs(userID, songIDs, playlistName, playlistDesc)

		// Assert
		assert.Error(t, err, "CreatePlaylistFromSongs should return error when no client exists")
		assert.Contains(t, err.Error(), "no spotify client found", "Error message should indicate client not found")
		assert.Empty(t, url, "URL should be empty when error occurs")
	})
}

func TestDeletePlaylist(t *testing.T) {
	t.Run("Delete_Success", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockSpotifySongRepository)
		mockClient := new(MockSpotifyClient)
		clientManager := NewClientManager()
		service := NewSpotifyService(clientManager, mockRepo)

		// Test data
		userID := "test-user"
		playlistID := "playlist123"

		// Store mock client
		clientManager.StoreClient(userID, mockClient)

		// Setup client mock for GetPlaylist
		mockClient.On("GetPlaylist", spotify.ID(playlistID)).Return(&spotify.FullPlaylist{
			SimplePlaylist: spotify.SimplePlaylist{
				ID: spotify.ID(playlistID),
			},
		}, nil)

		// Setup client mock for UnfollowPlaylist
		mockClient.On("UnfollowPlaylist", spotify.ID(userID), spotify.ID(playlistID)).Return(nil)

		// Act
		err := service.DeletePlaylist(userID, playlistID)

		// Assert
		assert.NoError(t, err, "DeletePlaylist should not return error on successful deletion")
		mockClient.AssertExpectations(t)
	})

	t.Run("No_Client_Error", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockSpotifySongRepository)
		clientManager := NewClientManager()
		service := NewSpotifyService(clientManager, mockRepo)

		// Test data
		userID := "test-user"
		playlistID := "playlist123"

		// Act - No client registered for this user
		err := service.DeletePlaylist(userID, playlistID)

		// Assert
		assert.Error(t, err, "DeletePlaylist should return error when no client exists")
		assert.Contains(t, err.Error(), "no spotify client found", "Error message should indicate client not found")
	})

	t.Run("Get_Playlist_Error", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockSpotifySongRepository)
		mockClient := new(MockSpotifyClient)
		clientManager := NewClientManager()
		service := NewSpotifyService(clientManager, mockRepo)

		// Test data
		userID := "test-user"
		playlistID := "playlist123"
		expectedError := errors.New("playlist not found")

		// Store mock client
		clientManager.StoreClient(userID, mockClient)

		// Setup client mock for GetPlaylist to return error
		mockClient.On("GetPlaylist", spotify.ID(playlistID)).Return(&spotify.FullPlaylist{}, expectedError)

		// Act
		err := service.DeletePlaylist(userID, playlistID)

		// Assert
		assert.Error(t, err, "DeletePlaylist should return error when GetPlaylist fails")
		assert.Contains(t, err.Error(), "failed to get playlist", "Error message should indicate GetPlaylist failure")
		mockClient.AssertExpectations(t)
	})
}

func TestGetPlaylistImageURL(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockSpotifySongRepository)
		mockClient := new(MockSpotifyClient)
		clientManager := NewClientManager()
		service := NewSpotifyService(clientManager, mockRepo)

		// Test data
		userID := "test-user"
		playlistID := "playlist123"
		expectedURL := "https://example.com/image.jpg"

		// Store mock client
		clientManager.StoreClient(userID, mockClient)

		// Setup client mock for GetPlaylist
		mockClient.On("GetPlaylist", spotify.ID(playlistID)).Return(&spotify.FullPlaylist{
			SimplePlaylist: spotify.SimplePlaylist{
				Images: []spotify.Image{
					{URL: expectedURL},
				},
			},
		}, nil)

		// Act
		url, err := service.GetPlaylistImageURL(userID, playlistID)

		// Assert
		require.NoError(t, err, "GetPlaylistImageURL should not return error")
		assert.Equal(t, expectedURL, url, "Returned URL should match expected image URL")
		mockClient.AssertExpectations(t)
	})

	t.Run("No_Client_Error", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockSpotifySongRepository)
		clientManager := NewClientManager()
		service := NewSpotifyService(clientManager, mockRepo)

		// Test data
		userID := "test-user"
		playlistID := "playlist123"

		// Act - No client registered for this user
		url, err := service.GetPlaylistImageURL(userID, playlistID)

		// Assert
		assert.Error(t, err, "GetPlaylistImageURL should return error when no client exists")
		assert.Contains(t, err.Error(), "no spotify client found", "Error message should indicate client not found")
		assert.Empty(t, url, "URL should be empty when error occurs")
	})
}

func TestRemoveDuplicates(t *testing.T) {
	// Test with duplicates
	t.Run("With_Duplicates", func(t *testing.T) {
		// Arrange
		ids := []spotify.ID{"song1", "song2", "song1", "song3", "song2", "song4"}
		expected := []spotify.ID{"song1", "song2", "song3", "song4"}

		// Act
		result := removeDuplicates(ids)

		// Assert
		assert.ElementsMatch(t, expected, result, "Should remove duplicate IDs")
		assert.Len(t, result, 4, "Result should have 4 unique elements")
	})

	// Test without duplicates
	t.Run("Without_Duplicates", func(t *testing.T) {
		// Arrange
		ids := []spotify.ID{"song1", "song2", "song3", "song4"}

		// Act
		result := removeDuplicates(ids)

		// Assert
		assert.ElementsMatch(t, ids, result, "Should return same elements when no duplicates")
		assert.Len(t, result, 4, "Result should have 4 elements")
	})

	// Test with empty slice
	t.Run("Empty_Slice", func(t *testing.T) {
		// Arrange
		var ids []spotify.ID

		// Act
		result := removeDuplicates(ids)

		// Assert
		assert.Empty(t, result, "Should return empty slice for empty input")
	})
}
