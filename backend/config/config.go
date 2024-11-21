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
	DBHost              string
	DBPort              string
	DBUser              string
	DBName              string
	DBPassword          string
	SamplesDBUser       string
	SamplesDBPassword   string
	SamplesDBPort       string
	SamplesDBHost       string
	SamplesDBName       string
	JWTSecret           string
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
		DBHost:              getEnv("MYSQL_HOST", ""),
		DBPort:              getEnv("MYSQL_PORT", ""),
		DBUser:              getEnv("MYSQL_USER", ""),
		DBName:              getEnv("MYSQL_DATABASE", ""),
		DBPassword:          getEnv("MYSQL_PASSWORD", ""),
		JWTSecret:           getEnv("JWT_SECRET", ""),
		SamplesDBUser:       getEnv("SAMPLES_DB_USER", ""),
		SamplesDBPassword:   getEnv("SAMPLES_DB_PASSWORD", ""),
		SamplesDBPort:       getEnv("SAMPLES_DB_PORT", ""),
		SamplesDBHost:       getEnv("SAMPLES_DB_HOST", ""),
		SamplesDBName:       getEnv("SAMPLES_DB_NAME", ""),
	}, nil
}
