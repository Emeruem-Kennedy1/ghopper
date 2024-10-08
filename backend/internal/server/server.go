package server

import (
	"fmt"

	"github.com/Emeruem-Kennedy1/ghopper/config"
	"github.com/Emeruem-Kennedy1/ghopper/internal/auth"
	"github.com/Emeruem-Kennedy1/ghopper/internal/handlers"
	"github.com/Emeruem-Kennedy1/ghopper/internal/middleware"
	"github.com/Emeruem-Kennedy1/ghopper/internal/repository"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router      *gin.Engine
	config      *config.Config
	spotifyAuth *auth.SpotifyAuth
	userRepo    *repository.UserRepository
}

func NewServer(cfg *config.Config, userRepo *repository.UserRepository) (*Server, error) {
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	spotifyAuth, err := auth.NewSpotifyAuth(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create spotify auth: %v", err)
	}

	s := &Server{
		router:      r,
		config:      cfg,
		spotifyAuth: spotifyAuth,
		userRepo:    userRepo,
	}

	s.setupRoutes()
	return s, nil
}

func (s *Server) setupRoutes() {
	s.router.GET("/ping", handlers.Ping())
	s.router.GET("/auth/spotify/login", handlers.SpotifyLogin(s.spotifyAuth))
	s.router.GET("/auth/spotify/callback", handlers.SpotifyCallback(s.spotifyAuth, s.userRepo))

	protected := s.router.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/user", handlers.GetUser(s.userRepo))
	}
}

func (s *Server) Run() error {
	return s.router.Run(":" + s.config.Port)
}
