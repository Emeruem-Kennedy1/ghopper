package middleware

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/Emeruem-Kennedy1/ghopper/internal/repository"
	"github.com/gin-gonic/gin"
)

// NonSpotifyAuthMiddleware authenticates non-Spotify users using their ID and passphrase
func NonSpotifyAuthMiddleware(userRepo repository.NonSpotifyUserRepositoryInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Check if it's Basic auth
		if !strings.HasPrefix(authHeader, "Basic ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Basic authentication is required"})
			c.Abort()
			return
		}

		// Extract credentials
		credentials := strings.TrimPrefix(authHeader, "Basic ")
		decodedBytes, err := base64.StdEncoding.DecodeString(credentials)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials format"})
			c.Abort()
			return
		}
		decoded := string(decodedBytes)

		// Split into ID and passphrase
		parts := strings.SplitN(decoded, ":", 2)
		if len(parts) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials format"})
			c.Abort()
			return
		}

		userID := parts[0]
		passphrase := parts[1]

		// Verify credentials
		isValid, err := userRepo.Verify(userID, passphrase)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Authentication error"})
			c.Abort()
			return
		}

		if !isValid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			c.Abort()
			return
		}

		// Set user ID in context for handlers to use
		c.Set("userID", userID)
		c.Set("isNonSpotifyUser", true)
		c.Next()
	}
}
