package main

import (
	"log"

	"github.com/Emeruem-Kennedy1/ghopper/config"
	"github.com/Emeruem-Kennedy1/ghopper/internal/database"
	"github.com/Emeruem-Kennedy1/ghopper/internal/repository"
	"github.com/Emeruem-Kennedy1/ghopper/internal/server"
)

func main() {
	// load env vars
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// init db
	dbs, err := database.InitDBConnections(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize databases: %v", err)
	}

	userRepo := repository.NewUserRepository(dbs.AppDB)
	songRepo := repository.NewSongRepository(dbs.SamplesDB)
	spotifySongRepo := repository.NewSpotifySongRepository(dbs.AppDB)

	// init and start server
	s, err := server.NewServer(cfg, userRepo, songRepo, spotifySongRepo)

	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}
	if err := s.Run(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
