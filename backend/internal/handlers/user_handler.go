package handlers

import (
	"net/http"

	"github.com/Emeruem-Kennedy1/ghopper/internal/repository"
	"github.com/Emeruem-Kennedy1/ghopper/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/zmb3/spotify"
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
			ctx.JSON(401, gin.H{"error": "Unauthorized"})
			return
		}

		user, err := userRepo.GetByID(userID.(string))
		if err != nil {
			ctx.JSON(500, gin.H{"error": "Failed to get user"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"user": user})
	}
}

func GetUserTopArtists(clientManager *services.ClientManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, exists := ctx.Get("userID")
		if !exists {
			ctx.JSON(401, gin.H{"error": "Unauthorized"})
			return
		}

		client, exists := clientManager.GetClient(userID.(string))
		if !exists {
			ctx.JSON(401, gin.H{"error": "Unauthorized"})
			return
		}

		limit := 10
		timeRange := "long"

		artists, err := client.CurrentUsersTopArtistsOpt(&spotify.Options{Limit: &limit, Timerange: &timeRange})
		if err != nil {
			ctx.JSON(500, gin.H{"error": "Failed to get user's top artists"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"artists": artists})
	}
}

func GetUserTopTracks(clientManager *services.ClientManager, spotifyService *services.SpotifyService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, exists := ctx.Get("userID")
		if !exists {
			ctx.JSON(401, gin.H{"error": "Unauthorized"})
			return
		}

		client, exists := clientManager.GetClient(userID.(string))
		if !exists {
			ctx.JSON(401, gin.H{"error": "Unauthorized"})
			return
		}

		limit := 10
		timeRange := "long"

		tracksRes, err := client.CurrentUsersTopTracksOpt(&spotify.Options{Limit: &limit, Timerange: &timeRange})

		if err != nil {
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

		ctx.JSON(http.StatusOK, gin.H{"tracks": tracks})
	}
}

func GetUserPlaylists(spotifySongRepo *repository.SpotifySongRepository, spotifyService *services.SpotifyService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, exists := ctx.Get("userID")
		if !exists {
			ctx.JSON(401, gin.H{"error": "Unauthorized"})
			return
		}

		playlists, err := spotifySongRepo.GetUserPlaylists(userID.(string))
		if err != nil {
			ctx.JSON(500, gin.H{"error": "Failed to get user's playlists"})
			return
		}
		var playlistResponse []PlaylistResponse
		// Go through each playlist and get the images

		if len(playlists) == 0 {
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

		ctx.JSON(http.StatusOK, gin.H{"playlists": playlistResponse})
	}
}
