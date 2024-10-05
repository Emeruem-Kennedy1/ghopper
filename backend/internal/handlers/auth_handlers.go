package handlers

import (
	"net/http"

	"github.com/Emeruem-Kennedy1/ghopper/internal/auth"
	"github.com/Emeruem-Kennedy1/ghopper/internal/repository"
	"github.com/gin-gonic/gin"
)

func SpotifyLogin(spotifyAuth *auth.SpotifyAuth) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		url := spotifyAuth.AuthURL()
		ctx.Redirect(http.StatusTemporaryRedirect, url)
	}
}

func SpotifyCallback(spotufyAuth *auth.SpotifyAuth, userRepo *repository.UserRepository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		client, err := spotufyAuth.CallBack(ctx.Request)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		spotifyUser, err := spotufyAuth.GetUserInfo(client)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
			return
		}

		// create or update user in the database
		user, token, err := auth.CreateOrUpdateUserFromSpotifyData(userRepo, *spotifyUser)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create or update user"})
			return
		}

		// TODO: Create jwt token and send it back to the client

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Successfully authenticated",
			"user":    user,
			"token":   token,
		})
	}
}
