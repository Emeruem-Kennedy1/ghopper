package handlers

import (
	"net/http"

	"github.com/Emeruem-Kennedy1/ghopper/internal/repository"
	"github.com/gin-gonic/gin"
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
