package handlers

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Ping() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		zap.L().Info("Ping request received")
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	}
}
