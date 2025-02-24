package handlers

import (
	"fmt"
	"net/http"

	"github.com/Emeruem-Kennedy1/ghopper/internal/repository"
	"github.com/Emeruem-Kennedy1/ghopper/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/zmb3/spotify"
	"go.uber.org/zap"
)

type Song struct {
	ID      string   `json:"id"`
	Artists []string `json:"artists"`
	Name    string   `json:"name"`
	Image   string   `json:"image"`
}

type PlaylistResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Image       string `json:"image"`
}

func GetUser(userRepo *repository.UserRepository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, exists := ctx.Get("userID")
		if !exists {
			zap.L().Warn("Unauthorized attempt to get user")
			ctx.JSON(401, gin.H{"error": "Unauthorized"})
			return
		}

		user, err := userRepo.GetByID(userID.(string))
		if err != nil {
			zap.L().Error("Failed to retrieve user from database",
				zap.String("userID", userID.(string)),
				zap.Error(err))
			ctx.JSON(500, gin.H{"error": "Failed to get user"})
			return
		}

		zap.L().Info("Successfully retrieved user details", zap.String("userID", userID.(string)))
		ctx.JSON(http.StatusOK, gin.H{"user": user})
	}
}

func GetUserTopArtists(clientManager *services.ClientManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, exists := ctx.Get("userID")
		if !exists {
			zap.L().Warn("Unauthorized attempt to get user's top artists")
			ctx.JSON(401, gin.H{"error": "Unauthorized"})
			return
		}

		client, exists := clientManager.GetClient(userID.(string))
		if !exists {
			zap.L().Warn("No Spotify client found for user",
				zap.String("userID", userID.(string)))
			ctx.JSON(401, gin.H{"error": "Unauthorized"})
			return
		}

		limit := 25
		timeRange := "short"

		artists, err := client.CurrentUsersTopArtistsOpt(&spotify.Options{Limit: &limit, Timerange: &timeRange})
		if err != nil {
			zap.L().Error("Failed to fetch top artists from Spotify",
				zap.String("userID", userID.(string)),
				zap.Error(err))
			ctx.JSON(500, gin.H{"error": "Failed to get user's top artists"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"artists": artists})
		zap.L().Info("Successfully retrieved user's top artists",
			zap.String("userID", userID.(string)),
			zap.Int("count", len(artists.Artists)))
	}
}

func GetUserTopTracks(clientManager *services.ClientManager, spotifyService *services.SpotifyService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, exists := ctx.Get("userID")
		if !exists {
			zap.L().Warn("Unauthorized attempt to get user's top tracks")
			ctx.JSON(401, gin.H{"error": "Unauthorized"})
			return
		}

		client, exists := clientManager.GetClient(userID.(string))
		if !exists {
			zap.L().Warn("No Spotify client found for user",
				zap.String("userID", userID.(string)))
			ctx.JSON(401, gin.H{"error": "Unauthorized"})
			return
		}

		limit := 25
		timeRange := "short"

		tracksRes, err := client.CurrentUsersTopTracksOpt(&spotify.Options{Limit: &limit, Timerange: &timeRange})

		if err != nil {
			zap.L().Error("Failed to fetch top tracks from Spotify",
				zap.String("userID", userID.(string)),
				zap.Error(err))
			ctx.JSON(500, gin.H{"error": "Failed to get user's top tracks"})
			return
		}

		var tracks []Song
		for _, track := range tracksRes.Tracks {
			artists := make([]string, 0)
			for _, artist := range track.Artists {
				artists = append(artists, artist.Name)
			}

			tracks = append(tracks, Song{
				ID:      track.ID.String(),
				Artists: artists,
				Name:    track.Name,
				Image:   track.Album.Images[0].URL,
			})
		}

		zap.L().Info("Successfully retrieved user's top tracks",
			zap.String("userID", userID.(string)),
			zap.Int("count", len(tracks)))

		ctx.JSON(http.StatusOK, gin.H{"tracks": tracks})
	}
}

func GetUserPlaylists(spotifySongRepo *repository.SpotifySongRepository, spotifyService *services.SpotifyService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, exists := ctx.Get("userID")
		if !exists {
			zap.L().Warn("Unauthorized attempt to get user's playlists")
			ctx.JSON(401, gin.H{"error": "Unauthorized"})
			return
		}

		playlists, err := spotifySongRepo.GetUserPlaylists(userID.(string))
		if err != nil {
			zap.L().Error("Failed to fetch user's playlists from database",
				zap.String("userID", userID.(string)),
				zap.Error(err))

			ctx.JSON(500, gin.H{"error": "Failed to get user's playlists"})
			return
		}
		var playlistResponse []PlaylistResponse
		// Go through each playlist and get the images

		if len(playlists) == 0 {
			zap.L().Info("No playlists found for user",
				zap.String("userID", userID.(string)))
			ctx.JSON(http.StatusOK, gin.H{"playlists": playlistResponse, "message": "No playlists found"})
			return
		}

		for _, playlist := range playlists {
			var playlistImage string
			if playlist.Image == "" {
				playlistImage, err = spotifyService.GetPlaylistImageURL(userID.(string), playlist.ID)
				spotifySongRepo.UpdatePlaylistImageURL(playlistImage, &playlist)
			} else {
				playlistImage = playlist.Image
			}

			if err != nil {
				zap.L().Error("Failed to get playlist image",
					zap.String("userID", userID.(string)),
					zap.String("playlistID", playlist.ID),
					zap.Error(err))
				ctx.JSON(500, gin.H{"error": "Failed to get playlist image"})
				return
			}
			playlistResponse = append(playlistResponse, PlaylistResponse{
				ID:          playlist.ID,
				Name:        playlist.Name,
				Description: playlist.Description,
				URL:         playlist.URL,
				Image:       playlistImage,
			})
		}

		zap.L().Info("Successfully retrieved user's playlists",
			zap.String("userID", userID.(string)),
			zap.Int("count", len(playlistResponse)))

		ctx.JSON(http.StatusOK, gin.H{"playlists": playlistResponse})
	}
}

func DeleteUserAccount(userRepo *repository.UserRepository, spotifySongRepo *repository.SpotifySongRepository, clientManager *services.ClientManager) gin.HandlerFunc {
    return func(ctx *gin.Context) {
        userID, exists := ctx.Get("userID")
        if !exists {
            ctx.JSON(401, gin.H{"error": "Unauthorized"})
            return
        }

		fmt.Println("Deleting user account")
        // Delete all playlists associated with user
        if err := spotifySongRepo.DeleteUserPlaylists(userID.(string)); err != nil {
            ctx.JSON(500, gin.H{"error": "Failed to delete user playlists"})
            return
        }

		fmt.Println("Removing user from client manager")
        // Remove the client from clientManager
        clientManager.RemoveClient(userID.(string))

        // Delete the user
		fmt.Println("Deleting user")
        if err := userRepo.Delete(userID.(string)); err != nil {
            ctx.JSON(500, gin.H{"error": "Failed to delete user account"})
            return
        }

		fmt.Println("User account successfully deleted")

        ctx.JSON(http.StatusOK, gin.H{"message": "Account successfully deleted"})
    }
}