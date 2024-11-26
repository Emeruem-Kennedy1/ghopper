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

	// init and start server
	s, err := server.NewServer(cfg, userRepo, songRepo)

	// songQuery := models.SongQuery{
	// 	Title:  "Why You Wanna Trip on Me",
	// 	Artist: "Michael Jackson",
	// }

	// results, any_err := songRepo.FindSongsByGenreBFS([]models.SongQuery{songQuery}, "Rap", 2)
	// if any_err != nil {
	// 	log.Fatalf("failed to search songs: %v", err)
	// }
	// fmt.Printf("Results: %v\n", results)
	// for _, result := range results {
	// 	fmt.Printf("Found song %s at distance %d\n", result.MatchedSong.Title, result.Distance)
	// }

	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}
	if err := s.Run(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
