package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Emeruem-Kennedy1/ghopper/internal/auth"
	"github.com/Emeruem-Kennedy1/ghopper/internal/repository"
	"github.com/Emeruem-Kennedy1/ghopper/internal/services"
	"github.com/gin-gonic/gin"
)

func SpotifyLogin(spotifyAuth *auth.SpotifyAuth) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		url := spotifyAuth.AuthURL()
		ctx.Redirect(http.StatusTemporaryRedirect, url)
	}
}

func SpotifyCallback(spotufyAuth *auth.SpotifyAuth, userRepo *repository.UserRepository, clientManager *services.ClientManager) gin.HandlerFunc {
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

		clientManager.StoreClient(spotifyUser.ID, client)

		// create or update user in the database
		user, token, err := auth.CreateOrUpdateUserFromSpotifyData(userRepo, *spotifyUser)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create or update user"})
			return
		}

		// Create a map with the data you want to send
		data := map[string]interface{}{
			"message": "Successfully authenticated",
			"user":    user,
			"token":   token,
		}

		// Convert the data to JSON
		jsonData, err := json.Marshal(data)
		if err != nil {
			redirectWithError(ctx, "Failed to process user data")
			return
		}

		// Encode the JSON data to base64 to safely include it in a URL
		encodedData := base64.URLEncoding.EncodeToString(jsonData)

		// Redirect to your frontend URL with the encoded data
		frontendURL := "http://localhost:51920/auth-callback" // TODO: Put this into the config
		redirectURL := fmt.Sprintf("%s?data=%s", frontendURL, url.QueryEscape(encodedData))
		ctx.Redirect(http.StatusTemporaryRedirect, redirectURL)
	}
}

func redirectWithError(ctx *gin.Context, errorMessage string) {
	frontendURL := "http://localhost:51920/auth-callback"
	redirectURL := fmt.Sprintf("%s?error=%s", frontendURL, url.QueryEscape(errorMessage))
	ctx.Redirect(http.StatusTemporaryRedirect, redirectURL)
}
