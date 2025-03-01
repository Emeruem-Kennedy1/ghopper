package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Emeruem-Kennedy1/ghopper/internal/auth"
	"github.com/Emeruem-Kennedy1/ghopper/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupAuthTest() *gin.Engine {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)
	
	// Create a new Gin router with the auth middleware
	r := gin.New()
	r.Use(AuthMiddleware())
	
	// Add a test handler that returns the userID from context
	r.GET("/protected", func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "userID not found in context"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"userId": userID})
	})
	
	return r
}

func TestAuthMiddleware(t *testing.T) {
	t.Run("Valid_Token", func(t *testing.T) {
		// Setup
		r := setupAuthTest()
		
		// Create a valid token
		user := &models.User{ID: "test-user-id"}
		token, err := auth.GenerateToken(user)
		require.NoError(t, err, "Setup: GenerateToken failed")
		
		// Create a test HTTP request with the token
		req := httptest.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp := httptest.NewRecorder()
		
		// Serve the request
		r.ServeHTTP(resp, req)
		
		// Assert
		assert.Equal(t, http.StatusOK, resp.Code, "Should return OK status for valid token")
		assert.Contains(t, resp.Body.String(), "test-user-id", "Response should contain the user ID")
	})
	
	t.Run("Missing_Authorization_Header", func(t *testing.T) {
		// Setup
		r := setupAuthTest()
		
		// Create a test HTTP request without Authorization header
		req := httptest.NewRequest("GET", "/protected", nil)
		resp := httptest.NewRecorder()
		
		// Serve the request
		r.ServeHTTP(resp, req)
		
		// Assert
		assert.Equal(t, http.StatusUnauthorized, resp.Code, "Should return Unauthorized status for missing header")
		assert.Contains(t, resp.Body.String(), "Authriztion header is required", "Response should mention missing header")
	})
	
	t.Run("Invalid_Token_Format", func(t *testing.T) {
		// Setup
		r := setupAuthTest()
		
		// Create a test HTTP request with incorrectly formatted token
		req := httptest.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "InvalidFormat") // Not "Bearer token"
		resp := httptest.NewRecorder()
		
		// Serve the request
		r.ServeHTTP(resp, req)
		
		// Assert
		assert.Equal(t, http.StatusUnauthorized, resp.Code, "Should return Unauthorized status for invalid format")
		assert.Contains(t, resp.Body.String(), "Invalid token format", "Response should mention invalid format")
	})
	
	t.Run("Invalid_Token", func(t *testing.T) {
		// Setup
		r := setupAuthTest()
		
		// Create a test HTTP request with invalid token
		req := httptest.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer invalidtoken123")
		resp := httptest.NewRecorder()
		
		// Serve the request
		r.ServeHTTP(resp, req)
		
		// Assert
		assert.Equal(t, http.StatusUnauthorized, resp.Code, "Should return Unauthorized status for invalid token")
		assert.Contains(t, resp.Body.String(), "Invalid token", "Response should mention invalid token")
	})
}