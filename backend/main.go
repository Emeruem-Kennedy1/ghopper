package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": os.Getenv("TEST"),
	})
}

func main() {
	envFile := ".env.development"
	if os.Getenv("GO_ENV") == "production" {
		envFile = ".env.production"
	}

	// Try to load .env.production file
	if err := godotenv.Load(envFile); err != nil {
		print(err)
		log.Printf("No %s file found, using environment variables", envFile)
	}

	// Set Gin to production mode for better performance
	if os.Getenv("GO_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	r.GET("/ping", ping)

	port := os.Getenv("PORT")
	if port == "" {
		port = "9797"
	}
	r.Run(":" + port)
}
