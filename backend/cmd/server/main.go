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
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	userRepo := repository.NewUserRepository(db)

	// init and start server
	s, err := server.NewServer(cfg, userRepo)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}
	if err := s.Run(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
