package main

import (
	"log"

	"github.com/Emeruem-Kennedy1/ghopper/config"
	"github.com/Emeruem-Kennedy1/ghopper/internal/database"
	"github.com/Emeruem-Kennedy1/ghopper/internal/logging"
	"github.com/Emeruem-Kennedy1/ghopper/internal/repository"
	"github.com/Emeruem-Kennedy1/ghopper/internal/server"
	"go.uber.org/zap"
)

func main() {
	// init logger
	logger, err := logging.NewLogger()
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Sync()

	zap.ReplaceGlobals(logger)

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
	nonSpotifyUserRepo := repository.NewNonSpotifyUserRepository(dbs.AppDB)

	// init and start server
	s, err := server.NewServer(cfg, userRepo, songRepo, spotifySongRepo, nonSpotifyUserRepo, logger)

	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}
	if err := s.Run(); err != nil {
		// log.Fatalf("Failed to start server: %v", err)
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
