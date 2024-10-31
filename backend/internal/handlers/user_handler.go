package handlers

import (
	"net/http"

	"github.com/Emeruem-Kennedy1/ghopper/internal/repository"
	"github.com/Emeruem-Kennedy1/ghopper/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/zmb3/spotify"
)

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
