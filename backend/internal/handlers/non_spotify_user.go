package handlers

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Emeruem-Kennedy1/ghopper/internal/models"
	"github.com/Emeruem-Kennedy1/ghopper/internal/repository"
	"github.com/Emeruem-Kennedy1/ghopper/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// RegisterNonSpotifyUserRequest contains data to register a new non-Spotify user
type RegisterNonSpotifyUserRequest struct {
	UserID string `json:"user_id" binding:"required"`
}

// RegisterNonSpotifyUserResponse contains the response for a user registration
type RegisterNonSpotifyUserResponse struct {
	UserID     string `json:"user_id"`
	Passphrase string `json:"passphrase"`
}

// VerifyNonSpotifyUserRequest contains credentials for verification
type VerifyNonSpotifyUserRequest struct {
	UserID     string `json:"user_id" binding:"required"`
	Passphrase string `json:"passphrase" binding:"required"`
}

// NonSpotifyPlaylistRequest contains data to generate a playlist
type NonSpotifyPlaylistRequest struct {
	SeedTracks []struct {
		Title  string `json:"title" binding:"required"`
		Artist string `json:"artist" binding:"required"`
	} `json:"seed_tracks" binding:"required,min=1,max=20"`
	Genre string `json:"genre" binding:"required"`
}

// UpdateTrackStatusRequest contains data to update a track's status
type UpdateTrackStatusRequest struct {
	AddedToPlaylist bool `json:"added_to_playlist"`
}

// RegisterNonSpotifyUser handles registration of new non-Spotify users
func RegisterNonSpotifyUser(userRepo *repository.NonSpotifyUserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterNonSpotifyUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}

		// Check if user already exists
		existingUser, err := userRepo.FindByID(req.UserID)
		if err != nil {
			zap.L().Error("Failed to check for existing user", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process registration"})
			return
		}

		if existingUser != nil {
			c.JSON(http.StatusConflict, gin.H{"error": "User ID already exists"})
			return
		}

		// Generate a random passphrase (4 words for easy memorization)
		passphrase, err := utils.GenerateRandomWords(4)
		if err != nil {
			zap.L().Error("Failed to generate passphrase", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate passphrase"})
			return
		}

		// Create the new user
		newUser := &models.NonSpotifyUser{
			ID:         req.UserID,
			Passphrase: passphrase,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		if err := userRepo.Create(newUser); err != nil {
			zap.L().Error("Failed to create user", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		c.JSON(http.StatusCreated, RegisterNonSpotifyUserResponse{
			UserID:     req.UserID,
			Passphrase: passphrase,
		})
	}
}

// VerifyNonSpotifyUser verifies a non-Spotify user's credentials
func VerifyNonSpotifyUser(userRepo *repository.NonSpotifyUserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req VerifyNonSpotifyUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}

		isValid, err := userRepo.Verify(req.UserID, req.Passphrase)
		if err != nil {
			zap.L().Error("Failed to verify user", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify user"})
			return
		}

		if !isValid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID or passphrase"})
			return
		}

		// Set user ID in context for subsequent requests
		c.Set("userID", req.UserID)
		c.Set("isNonSpotifyUser", true)

		c.JSON(http.StatusOK, gin.H{"message": "User verified successfully"})
	}
}

// GenerateNonSpotifyPlaylist creates a playlist based on seed tracks
func GenerateNonSpotifyPlaylist(
	userRepo *repository.NonSpotifyUserRepository,
	songRepo repository.SongRepositoryInterface,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists || userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		var req NonSpotifyPlaylistRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}

		// Convert seed tracks to song queries
		songQueries := make([]models.SongQuery, len(req.SeedTracks))
		for i, track := range req.SeedTracks {
			songQueries[i] = models.SongQuery{
				Title:  track.Title,
				Artist: track.Artist,
			}
		}

		// Search for songs by genre using the sample song repository
		maxDepth := 2 // Adjust as needed
		searchResults, err := songRepo.FindSongsByGenreBFS(songQueries, req.Genre, maxDepth)
		if err != nil {
			zap.L().Error("Failed to search for songs", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate playlist"})
			return
		}

		if len(searchResults) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "No matching songs found for the given genre and seed tracks"})
			return
		}

		// Create a map to avoid duplicates
		uniqueTracks := make(map[string]models.NonSpotifyPlaylistTrack)

		// Helper function to check if a track matches any seed track
		isASeedTrack := func(title, artist string) bool {
			for _, seed := range req.SeedTracks {
				if strings.EqualFold(seed.Title, title) && strings.EqualFold(seed.Artist, artist) {
					return true
				}
			}
			return false
		}

		for _, result := range searchResults {

			if isASeedTrack(result.MatchedSong.Title, result.MatchedSong.Artists[0].Name) {
				continue
			}

			// Use artist and title as the key
			key := fmt.Sprintf("%s-%s", result.MatchedSong.Title, result.MatchedSong.Artists[0].Name)
			if _, exists := uniqueTracks[key]; !exists {
				uniqueTracks[key] = models.NonSpotifyPlaylistTrack{
					ID:              uuid.New().String(),
					Title:           result.MatchedSong.Title,
					Artist:          result.MatchedSong.Artists[0].Name,
					AddedToPlaylist: false,
					CreatedAt:       time.Now(),
					UpdatedAt:       time.Now(),
				}
			}
		}

		// Convert map to slice
		tracks := make([]models.NonSpotifyPlaylistTrack, 0, len(uniqueTracks))
		for _, track := range uniqueTracks {
			tracks = append(tracks, track)
		}

		// Create seed tracks for the playlist
		seedTracks := make([]models.NonSpotifyPlaylistSeedTrack, len(req.SeedTracks))
		for i, seed := range req.SeedTracks {
			seedTracks[i] = models.NonSpotifyPlaylistSeedTrack{
				ID:        uuid.New().String(),
				Title:     seed.Title,
				Artist:    seed.Artist,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
		}

		// Create the playlist
		now := time.Now()
		dateStr := now.Format("2006-01-02-15-04")
		playlist := &models.NonSpotifyPlaylist{
			ID:          uuid.New().String(),
			UserID:      userID.(string),
			Name:        fmt.Sprintf("%s-%s-playlist", req.Genre, dateStr),
			Genre:       req.Genre,
			Description: fmt.Sprintf("Playlist generated for the %s genre", req.Genre),
			ImageURL:    generateImageURL(req.Genre), // Implement this function to generate image URLs
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		// Save everything to the database
		if err := userRepo.SavePlaylist(playlist, tracks, seedTracks); err != nil {
			zap.L().Error("Failed to save playlist", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save playlist"})
			return
		}

		// Return the created playlist with tracks
		c.JSON(http.StatusCreated, gin.H{
			"playlist": models.NonSpotifyPlaylistWithTracks{
				NonSpotifyPlaylist: *playlist,
				Tracks:             tracks,
				SeedTracks:         seedTracks,
			},
		})
	}
}

// GetNonSpotifyUserPlaylists retrieves all playlists for a user
func GetNonSpotifyUserPlaylists(userRepo *repository.NonSpotifyUserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists || userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		playlists, err := userRepo.GetUserPlaylists(userID.(string))
		if err != nil {
			zap.L().Error("Failed to get user playlists", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get playlists"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"playlists": playlists})
	}
}

// GetNonSpotifyPlaylistDetails retrieves a playlist with its tracks
func GetNonSpotifyPlaylistDetails(userRepo *repository.NonSpotifyUserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists || userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		playlistID := c.Param("playlistID")
		if playlistID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Playlist ID is required"})
			return
		}

		playlist, err := userRepo.GetPlaylistWithTracks(playlistID)
		if err != nil {
			zap.L().Error("Failed to get playlist details", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get playlist details"})
			return
		}

		if playlist == nil || playlist.UserID != userID.(string) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Playlist not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"playlist": playlist})
	}
}

// UpdateNonSpotifyTrackStatus updates a track's "added to playlist" status
func UpdateNonSpotifyTrackStatus(userRepo *repository.NonSpotifyUserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists || userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		trackID := c.Param("trackID")
		if trackID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Track ID is required"})
			return
		}

		var req UpdateTrackStatusRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}

		// Update the track status
		if err := userRepo.UpdateTrackStatus(trackID, req.AddedToPlaylist); err != nil {
			zap.L().Error("Failed to update track status", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update track status"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Track status updated successfully"})
	}
}

// DeleteNonSpotifyPlaylist deletes a playlist
func DeleteNonSpotifyPlaylist(userRepo *repository.NonSpotifyUserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists || userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		playlistID := c.Param("playlistID")
		if playlistID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Playlist ID is required"})
			return
		}

		// Verify the playlist belongs to the user
		playlist, err := userRepo.GetPlaylistWithTracks(playlistID)
		if err != nil {
			zap.L().Error("Failed to get playlist for deletion", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete playlist"})
			return
		}

		if playlist == nil || playlist.UserID != userID.(string) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Playlist not found"})
			return
		}

		// Delete the playlist
		if err := userRepo.DeletePlaylist(playlistID); err != nil {
			zap.L().Error("Failed to delete playlist", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete playlist"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Playlist deleted successfully"})
	}
}

// Helper function to generate image URLs based on genre
func generateImageURL(genre string) string {
	normalizedGenre := strings.ToLower(strings.ReplaceAll(genre, " ", "-"))

	counts := map[string]int{
		"pop":        5,
		"jazz":       6,
		"hip-hop":    6,
		"country":    6,
		"electronic": 5,
		"reggae":     5,
		"soundtrack": 5,
		"world":      6,
		"classical":  5,
		"soul":       6,
	}

	genre_count, exists := counts[normalizedGenre]
	if !exists || genre_count <= 0 {
		genre_count = 1
	}

	max := big.NewInt(int64(genre_count))
	n, err := rand.Int(rand.Reader, max)

	if err != nil {
		zap.L().Error("Failed to generate random number", zap.Error(err))
		return "default.jpg" // TODO: add a default if random generation fails
	}
	number := n.Int64() + 1

	// Return a placeholder URL that includes the genre
	return fmt.Sprintf("%s_%s.jpg", normalizedGenre, strconv.FormatInt(number, 10))
}
