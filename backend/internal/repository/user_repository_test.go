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

// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err, "Failed to open in-memory database")

	// Migrate the schema for User model
	err = db.AutoMigrate(&models.User{})
	require.NoError(t, err, "Failed to migrate User model")

	return db
}

func TestUserRepository_Create(t *testing.T) {
	// Setup test database
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	t.Run("Create_Valid_User", func(t *testing.T) {
		// Arrange
		user := &models.User{
			ID:           "test-id-1",
			DisplayName:  "Test User",
			Email:        "test1@example.com",
			SpotifyURI:   "spotify:user:test1",
			Country:      "US",
			ProfileImage: "https://example.com/profile.jpg",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		// Act
		err := repo.Create(user)

		// Assert
		require.NoError(t, err, "Should successfully create a valid user")

		// Verify user was created
		var savedUser models.User
		result := db.First(&savedUser, "id = ?", user.ID)
		assert.NoError(t, result.Error, "Should find the created user")
		assert.Equal(t, user.ID, savedUser.ID, "Saved user ID should match")
		assert.Equal(t, user.Email, savedUser.Email, "Saved user email should match")
	})

	t.Run("Create_Duplicate_ID", func(t *testing.T) {
		// Arrange - Create a user
		user1 := &models.User{
			ID:          "duplicate-id",
			DisplayName: "First User",
			Email:       "first@example.com",
			SpotifyURI:  "spotify:user:first",
		}
		err := repo.Create(user1)
		require.NoError(t, err, "Setup: Should create first user")

		// Try to create another user with the same ID
		user2 := &models.User{
			ID:          "duplicate-id", // Same ID
			DisplayName: "Second User",
			Email:       "second@example.com",
			SpotifyURI:  "spotify:user:second",
		}

		// Act
		err = repo.Create(user2)

		// Assert
		assert.Error(t, err, "Should return error when creating user with duplicate ID")
	})
}

func TestUserRepository_UpsertUser(t *testing.T) {
	// Setup test database
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	t.Run("Insert_New_User", func(t *testing.T) {
		// Arrange
		user := &models.User{
			ID:          "new-user-id",
			DisplayName: "New User",
			Email:       "new@example.com",
			SpotifyURI:  "spotify:user:new",
		}

		// Act
		err := repo.UpsertUser(user)

		// Assert
		require.NoError(t, err, "Should successfully insert a new user")

		// Verify user was created
		var savedUser models.User
		result := db.First(&savedUser, "id = ?", user.ID)
		assert.NoError(t, result.Error, "Should find the inserted user")
		assert.Equal(t, user.DisplayName, savedUser.DisplayName, "Display name should match")
	})

	t.Run("Update_Existing_User", func(t *testing.T) {
		// Arrange - Create a user first
		userId := "existing-user-id"
		existingUser := &models.User{
			ID:          userId,
			DisplayName: "Existing User",
			Email:       "existing@example.com",
			SpotifyURI:  "spotify:user:existing",
			Country:     "US",
		}
		err := repo.Create(existingUser)
		require.NoError(t, err, "Setup: Should create existing user")

		// Create updated user with same ID
		updatedUser := &models.User{
			ID:          userId, // Same ID
			DisplayName: "Updated User",
			Email:       "updated@example.com",
			SpotifyURI:  "spotify:user:updated",
			Country:     "CA", // Changed country
		}

		// Act
		err = repo.UpsertUser(updatedUser)

		// Assert
		require.NoError(t, err, "Should successfully update existing user")

		// Verify user was updated
		var savedUser models.User
		result := db.First(&savedUser, "id = ?", userId)
		assert.NoError(t, result.Error, "Should find the updated user")
		assert.Equal(t, updatedUser.DisplayName, savedUser.DisplayName, "Display name should be updated")
		assert.Equal(t, updatedUser.Email, savedUser.Email, "Email should be updated")
		assert.Equal(t, updatedUser.Country, savedUser.Country, "Country should be updated")
	})
}

func TestUserRepository_GetByID(t *testing.T) {
	// Setup test database
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// Create a test user
	testUser := &models.User{
		ID:          "get-by-id-test",
		DisplayName: "Get By ID Test",
		Email:       "getbyid@example.com",
		SpotifyURI:  "spotify:user:getbyid",
	}
	err := repo.Create(testUser)
	require.NoError(t, err, "Setup: Should create test user")

	t.Run("Get_Existing_User", func(t *testing.T) {
		// Act
		user, err := repo.GetByID(testUser.ID)

		// Assert
		require.NoError(t, err, "Should not return error for existing user")
		assert.NotNil(t, user, "Retrieved user should not be nil")
		assert.Equal(t, testUser.ID, user.ID, "User ID should match")
		assert.Equal(t, testUser.Email, user.Email, "Email should match")
	})

	t.Run("Get_NonExistent_User", func(t *testing.T) {
		// Act
		user, err := repo.GetByID("non-existent-id")

		// Assert
		assert.Error(t, err, "Should return error for non-existent user")
		assert.Nil(t, user, "User should be nil")
	})
}

func TestUserRepository_Delete(t *testing.T) {
	// Setup test database
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// Create a test user
	testUser := &models.User{
		ID:          "delete-test",
		DisplayName: "Delete Test",
		Email:       "delete@example.com",
		SpotifyURI:  "spotify:user:delete",
	}
	err := repo.Create(testUser)
	require.NoError(t, err, "Setup: Should create test user")

	t.Run("Delete_Existing_User", func(t *testing.T) {
		// Act
		err := repo.Delete(testUser.ID)

		// Assert
		require.NoError(t, err, "Should not return error when deleting existing user")

		// Verify user was deleted
		var user models.User
		result := db.First(&user, "id = ?", testUser.ID)
		assert.Error(t, result.Error, "Should not find deleted user")
		assert.Equal(t, gorm.ErrRecordNotFound, result.Error, "Should return record not found error")
	})

	t.Run("Delete_NonExistent_User", func(t *testing.T) {
		// Act
		err := repo.Delete("non-existent-id")

		// Assert
		// Note: GORM's Delete doesn't return an error if the record doesn't exist
		assert.NoError(t, err, "Should not return error when deleting non-existent user")
	})
}
