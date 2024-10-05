package middleware

import (
	"net/http"
	"strings"

	"github.com/Emeruem-Kennedy1/ghopper/internal/auth"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")

		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authriztion header is required"})
		}

		bearertoken := strings.Split(authHeader, " ")
		if len(bearertoken) != 2 {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			ctx.Abort()
			return
		}

		userID, err := auth.ValidateToken(bearertoken[1])
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			ctx.Abort()
			return
		}

		ctx.Set("userID", userID)
		ctx.Next()
	}
}
