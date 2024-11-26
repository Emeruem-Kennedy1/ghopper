package server

import (
	"fmt"

	"github.com/Emeruem-Kennedy1/ghopper/config"
	"github.com/Emeruem-Kennedy1/ghopper/internal/auth"
	"github.com/Emeruem-Kennedy1/ghopper/internal/handlers"
	"github.com/Emeruem-Kennedy1/ghopper/internal/middleware"
	"github.com/Emeruem-Kennedy1/ghopper/internal/repository"
	"github.com/Emeruem-Kennedy1/ghopper/internal/services"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router        *gin.Engine
	config        *config.Config
	spotifyAuth   *auth.SpotifyAuth
	userRepo      *repository.UserRepository
	songRepo      *repository.SongRepository
	cleintManager *services.ClientManager
}

func NewServer(cfg *config.Config, userRepo *repository.UserRepository, songRepo *repository.SongRepository) (*Server, error) {
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	spotifyAuth, err := auth.NewSpotifyAuth(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create spotify auth: %v", err)
	}

	s := &Server{
		router:        r,
		config:        cfg,
		spotifyAuth:   spotifyAuth,
		userRepo:      userRepo,
		songRepo:      songRepo,
		cleintManager: services.NewClientManager(),
	}

	s.setupRoutes()
	return s, nil
}

func (s *Server) setupRoutes() {

	s.router.GET("/auth/spotify/login", handlers.SpotifyLogin(s.spotifyAuth))
	s.router.GET("/auth/spotify/callback", handlers.SpotifyCallback(s.spotifyAuth, s.userRepo, s.cleintManager))

	protected := s.router.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/ping", handlers.Ping())
		protected.GET("/user", handlers.GetUser(s.userRepo))
		protected.GET("/user/top-artists", handlers.GetUserTopArtists(s.cleintManager))
		protected.GET("/user/top-tracks", handlers.GetUserTopTracks(s.cleintManager))
		protected.POST("/search", handlers.SearchSongByGenre(s.songRepo))
	}
}

func (s *Server) Run() error {
	return s.router.Run(":" + s.config.Port)
}
