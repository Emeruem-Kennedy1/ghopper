package repository

import (
	"testing"
	"time"

	"github.com/Emeruem-Kennedy1/ghopper/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupSpotifySongTestDB creates an in-memory SQLite database for testing
func setupSpotifySongTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err, "Failed to open in-memory database")

	// Drop any existing tables first
	err = db.Migrator().DropTable(&models.User{}, &models.Song{}, &models.Playlist{})
	require.NoError(t, err, "Failed to drop existing tables")

	// Migrate the schema for Song and Playlist models
	err = db.AutoMigrate(&models.Song{}, &models.Playlist{})
	require.NoError(t, err, "Failed to migrate Song and Playlist models")

	return db
}

func TestSpotifySongRepository_FindSongByNameAndArtist(t *testing.T) {
	// Setup test database
	db := setupSpotifySongTestDB(t)
	repo := NewSpotifySongRepository(db)

	// Create a test song
	testSong := &models.Song{
		ID:        "test-song-id",
		Name:      "Test Song",
		Artist:    "Test Artist",
		CreatedAt: time.Now(),
	}
	err := db.Create(testSong).Error
	require.NoError(t, err, "Setup: Should create test song")

	t.Run("Find_Existing_Song", func(t *testing.T) {
		// Act
		song, err := repo.FindSongByNameAndArtist(testSong.Name, testSong.Artist)

		// Assert
		require.NoError(t, err, "Should not return error for existing song")
		assert.NotNil(t, song, "Retrieved song should not be nil")
		assert.Equal(t, testSong.ID, song.ID, "Song ID should match")
		assert.Equal(t, testSong.Name, song.Name, "Song name should match")
		assert.Equal(t, testSong.Artist, song.Artist, "Artist name should match")
	})

	t.Run("Find_NonExistent_Song", func(t *testing.T) {
		// Act
		song, err := repo.FindSongByNameAndArtist("Non-existent Song", "Non-existent Artist")

		// Assert
		assert.NoError(t, err, "Should not return error for non-existent song")
		assert.Nil(t, song, "Song should be nil")
	})
}

func TestSpotifySongRepository_SaveSong(t *testing.T) {
	// Setup test database
	db := setupSpotifySongTestDB(t)
	repo := NewSpotifySongRepository(db)

	t.Run("Save_New_Song", func(t *testing.T) {
		// Arrange
		song := &models.Song{
			ID:     "new-song-id",
			Name:   "New Song",
			Artist: "New Artist",
		}

		// Act
		err := repo.SaveSong(song)

		// Assert
		require.NoError(t, err, "Should successfully save a new song")

		// Verify song was saved
		var savedSong models.Song
		result := db.First(&savedSong, "id = ?", song.ID)
		assert.NoError(t, result.Error, "Should find the saved song")
		assert.Equal(t, song.ID, savedSong.ID, "Saved song ID should match")
		assert.Equal(t, song.Name, savedSong.Name, "Saved song name should match")
	})

	t.Run("Save_Existing_Song", func(t *testing.T) {
		// Arrange - Create a song first
		existingSong := &models.Song{
			ID:     "existing-song-id",
			Name:   "Existing Song",
			Artist: "Existing Artist",
		}
		err := db.Create(existingSong).Error
		require.NoError(t, err, "Setup: Should create existing song")

		// Try to save the same song again
		songToSave := &models.Song{
			ID:     "existing-song-id",
			Name:   "Updated Song",   // Changed name
			Artist: "Updated Artist", // Changed artist
		}

		// Act
		err = repo.SaveSong(songToSave)

		// Assert
		require.NoError(t, err, "Should not return error when saving existing song")

		// Verify song was not updated (as per your implementation)
		var savedSong models.Song
		result := db.First(&savedSong, "id = ?", existingSong.ID)
		assert.NoError(t, result.Error, "Should find the existing song")
		assert.Equal(t, existingSong.Name, savedSong.Name, "Song name should remain unchanged")
		assert.Equal(t, existingSong.Artist, savedSong.Artist, "Artist name should remain unchanged")
	})
}

func TestSpotifySongRepository_FindPlaylistByNameAndUser(t *testing.T) {
	// Setup test database
	db := setupSpotifySongTestDB(t)
	repo := NewSpotifySongRepository(db)

	// Create a test playlist
	testPlaylist := &models.Playlist{
		ID:     "test-playlist-id",
		Name:   "Test Playlist",
		UserID: "test-user-id",
	}
	err := db.Create(testPlaylist).Error
	require.NoError(t, err, "Setup: Should create test playlist")

	t.Run("Find_Existing_Playlist", func(t *testing.T) {
		// Act
		playlist, err := repo.FindPlaylistByNameAndUser(testPlaylist.Name, testPlaylist.UserID)

		// Assert
		require.NoError(t, err, "Should not return error for existing playlist")
		assert.NotNil(t, playlist, "Retrieved playlist should not be nil")
		assert.Equal(t, testPlaylist.ID, playlist.ID, "Playlist ID should match")
		assert.Equal(t, testPlaylist.Name, playlist.Name, "Playlist name should match")
	})

	t.Run("Find_NonExistent_Playlist", func(t *testing.T) {
		// Act
		playlist, err := repo.FindPlaylistByNameAndUser("Non-existent Playlist", "non-existent-user")

		// Assert
		assert.NoError(t, err, "Should not return error for non-existent playlist")
		assert.Nil(t, playlist, "Playlist should be nil")
	})
}

func TestSpotifySongRepository_SavePlaylist(t *testing.T) {
	// Setup test database
	db := setupSpotifySongTestDB(t)
	repo := NewSpotifySongRepository(db)

	t.Run("Save_New_Playlist", func(t *testing.T) {
		// Arrange
		playlist := &models.Playlist{
			ID:     "new-playlist-id",
			Name:   "New Playlist",
			UserID: "user-id",
			Image:  "https://example.com/image.jpg",
		}

		// Act
		err := repo.SavePlaylist(playlist)

		// Assert
		require.NoError(t, err, "Should successfully save a new playlist")

		// Verify playlist was saved
		var savedPlaylist models.Playlist
		result := db.First(&savedPlaylist, "id = ?", playlist.ID)
		assert.NoError(t, result.Error, "Should find the saved playlist")
		assert.Equal(t, playlist.ID, savedPlaylist.ID, "Saved playlist ID should match")
		assert.Equal(t, playlist.Name, savedPlaylist.Name, "Saved playlist name should match")
		assert.Equal(t, playlist.Image, savedPlaylist.Image, "Saved playlist image should match")
	})
}

func TestSpotifySongRepository_FindPlaylistByIDAndUser(t *testing.T) {
	// Setup test database
	db := setupSpotifySongTestDB(t)
	repo := NewSpotifySongRepository(db)

	// Create a test playlist
	testPlaylist := &models.Playlist{
		ID:     "test-playlist-id",
		Name:   "Test Playlist",
		UserID: "test-user-id",
	}
	err := db.Create(testPlaylist).Error
	require.NoError(t, err, "Setup: Should create test playlist")

	t.Run("Find_Existing_Playlist_By_ID", func(t *testing.T) {
		// Act
		playlist, err := repo.FindPlaylistByIDAndUser(testPlaylist.ID, testPlaylist.UserID)

		// Assert
		require.NoError(t, err, "Should not return error for existing playlist")
		assert.NotNil(t, playlist, "Retrieved playlist should not be nil")
		assert.Equal(t, testPlaylist.ID, playlist.ID, "Playlist ID should match")
		assert.Equal(t, testPlaylist.Name, playlist.Name, "Playlist name should match")
	})

	t.Run("Find_NonExistent_Playlist_By_ID", func(t *testing.T) {
		// Act
		playlist, err := repo.FindPlaylistByIDAndUser("non-existent-id", "test-user-id")

		// Assert
		assert.NoError(t, err, "Should not return error for non-existent playlist")
		assert.Nil(t, playlist, "Playlist should be nil")
	})

	t.Run("Find_Playlist_Wrong_User", func(t *testing.T) {
		// Act
		playlist, err := repo.FindPlaylistByIDAndUser(testPlaylist.ID, "wrong-user-id")

		// Assert
		assert.NoError(t, err, "Should not return error for playlist with wrong user")
		assert.Nil(t, playlist, "Playlist should be nil")
	})
}

func TestSpotifySongRepository_DeletePlaylist(t *testing.T) {
	// Setup test database
	db := setupSpotifySongTestDB(t)
	repo := NewSpotifySongRepository(db)

	// Create a test playlist
	testPlaylist := &models.Playlist{
		ID:     "delete-playlist-id",
		Name:   "Delete Playlist",
		UserID: "test-user-id",
	}
	err := db.Create(testPlaylist).Error
	require.NoError(t, err, "Setup: Should create test playlist")

	t.Run("Delete_Existing_Playlist", func(t *testing.T) {
		// Act
		err := repo.DeletePlaylist(testPlaylist)

		// Assert
		require.NoError(t, err, "Should not return error when deleting existing playlist")

		// Verify playlist was deleted
		var playlist models.Playlist
		result := db.First(&playlist, "id = ?", testPlaylist.ID)
		assert.Error(t, result.Error, "Should not find deleted playlist")
		assert.Equal(t, gorm.ErrRecordNotFound, result.Error, "Should return record not found error")
	})
}

func TestSpotifySongRepository_GetUserPlaylists(t *testing.T) {
	// Setup test database
	db := setupSpotifySongTestDB(t)
	repo := NewSpotifySongRepository(db)

	// Create test playlists
	userID := "test-user-id"
	testPlaylists := []models.Playlist{
		{
			ID:     "playlist-1",
			Name:   "Playlist 1",
			UserID: userID,
		},
		{
			ID:     "playlist-2",
			Name:   "Playlist 2",
			UserID: userID,
		},
		{
			ID:     "playlist-3",
			Name:   "Playlist 3",
			UserID: "different-user-id", // Different user
		},
	}

	for _, playlist := range testPlaylists {
		err := db.Create(&playlist).Error
		require.NoError(t, err, "Setup: Should create test playlist")
	}

	t.Run("Get_User_Playlists", func(t *testing.T) {
		// Act
		playlists, err := repo.GetUserPlaylists(userID)

		// Assert
		require.NoError(t, err, "Should not return error when getting user playlists")
		assert.Equal(t, 2, len(playlists), "Should return 2 playlists for the user")

		// Verify playlists belong to the user
		for _, playlist := range playlists {
			assert.Equal(t, userID, playlist.UserID, "Playlist should belong to the correct user")
		}
	})

	t.Run("Get_Playlists_Non_Existent_User", func(t *testing.T) {
		// Act
		playlists, err := repo.GetUserPlaylists("non-existent-user")

		// Assert
		require.NoError(t, err, "Should not return error for non-existent user")
		assert.Empty(t, playlists, "Should return empty slice for non-existent user")
	})
}

func TestSpotifySongRepository_UpdatePlaylistImageURL(t *testing.T) {
	// Setup test database
	db := setupSpotifySongTestDB(t)
	repo := NewSpotifySongRepository(db)

	// Create a test playlist
	testPlaylist := &models.Playlist{
		ID:     "image-playlist-id",
		Name:   "Image Playlist",
		UserID: "test-user-id",
		Image:  "https://example.com/old-image.jpg",
	}
	err := db.Create(testPlaylist).Error
	require.NoError(t, err, "Setup: Should create test playlist")

	t.Run("Update_Playlist_Image", func(t *testing.T) {
		// Arrange
		newImageURL := "https://example.com/new-image.jpg"

		// Act
		err := repo.UpdatePlaylistImageURL(newImageURL, testPlaylist)

		// Assert
		require.NoError(t, err, "Should not return error when updating playlist image")

		// Verify image was updated
		var updatedPlaylist models.Playlist
		result := db.First(&updatedPlaylist, "id = ?", testPlaylist.ID)
		assert.NoError(t, result.Error, "Should find the updated playlist")
		assert.Equal(t, newImageURL, updatedPlaylist.Image, "Playlist image should be updated")
	})
}

func TestSpotifySongRepository_DeleteUserPlaylists(t *testing.T) {
	// Setup test database
	db := setupSpotifySongTestDB(t)
	repo := NewSpotifySongRepository(db)

	// Create test playlists
	userID := "delete-user-id"
	otherUserID := "other-user-id"
	testPlaylists := []models.Playlist{
		{
			ID:     "u1-playlist-1",
			Name:   "User 1 Playlist 1",
			UserID: userID,
		},
		{
			ID:     "u1-playlist-2",
			Name:   "User 1 Playlist 2",
			UserID: userID,
		},
		{
			ID:     "u2-playlist-1",
			Name:   "User 2 Playlist 1",
			UserID: otherUserID,
		},
	}

	for _, playlist := range testPlaylists {
		err := db.Create(&playlist).Error
		require.NoError(t, err, "Setup: Should create test playlist")
	}

	t.Run("Delete_User_Playlists", func(t *testing.T) {
		// Act
		err := repo.DeleteUserPlaylists(userID)

		// Assert
		require.NoError(t, err, "Should not return error when deleting user playlists")

		// Verify user's playlists were deleted
		var count int64
		db.Model(&models.Playlist{}).Where("user_id = ?", userID).Count(&count)
		assert.Equal(t, int64(0), count, "Should have deleted all playlists for the user")

		// Verify other user's playlists were not deleted
		db.Model(&models.Playlist{}).Where("user_id = ?", otherUserID).Count(&count)
		assert.Equal(t, int64(1), count, "Should not have deleted playlists for other users")
	})
}
