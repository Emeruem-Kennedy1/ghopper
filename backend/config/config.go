package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                string
	Env                 string
	SpotifyClientID     string
	SpotifyClientSecret string
	SpotifyRedirectURI  string
}

func getEnv(key, fallack string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallack
}

func Load() (*Config, error) {
	envFile := ".env.development"
	if os.Getenv("GO_ENV") == "production" {
		envFile = ".env.production"
	}

	if err := godotenv.Load(envFile); err != nil {
		log.Printf("No %s file found, using environment variables", envFile)
	}

	return &Config{
		Port:                getEnv("PORT", "9797"),
		Env:                 getEnv("GO_ENV", "development"),
		SpotifyClientID:     getEnv("SPOTIFY_CLIENT_ID", ""),
		SpotifyClientSecret: getEnv("SPOTIFY_CLIENT_SECRET", ""),
		SpotifyRedirectURI:  getEnv("SPOTIFY_REDIRECT_URI", ""),
	}, nil
}
