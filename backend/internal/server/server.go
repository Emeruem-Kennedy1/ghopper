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
	"go.uber.org/zap"
)

type Server struct {
	router             *gin.Engine
	config             *config.Config
	spotifyAuth        *auth.SpotifyAuth
	userRepo           repository.UserRepositoryInterface
	songRepo           *repository.SongRepository
	spotifySongRepo    repository.SpotifySongRepositoryInterface
	nonSpotifyUserRepo *repository.NonSpotifyUserRepository
	cleintManager      services.ClientManagerInterface
	spotifyService     services.SpotifyServiceInterface
	logger             *zap.Logger
}

func NewServer(
	cfg *config.Config,
	userRepo *repository.UserRepository,
	songRepo *repository.SongRepository,
	spotifySongRepo *repository.SpotifySongRepository,
	nonSpotifyUserRepo *repository.NonSpotifyUserRepository,
	logger *zap.Logger,
) (*Server, error) {
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	spotifyAuth, err := auth.NewSpotifyAuth(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create spotify auth: %v", err)
	}

	clientManager := services.NewClientManager()

	spotifyService := services.NewSpotifyService(clientManager, spotifySongRepo)

	s := &Server{
		router:             r,
		config:             cfg,
		spotifyAuth:        spotifyAuth,
		userRepo:           userRepo,
		songRepo:           songRepo,
		spotifySongRepo:    spotifySongRepo,
		nonSpotifyUserRepo: nonSpotifyUserRepo,
		cleintManager:      clientManager,
		spotifyService:     spotifyService,
		logger:             logger,
	}
	gin.Logger()

	s.setupRoutes()
	return s, nil
}

func (s *Server) setupRoutes() {

	s.router.GET("/auth/spotify/login", handlers.SpotifyLogin(s.spotifyAuth))
	s.router.GET("/auth/spotify/callback", handlers.SpotifyCallback(s.spotifyAuth, s.userRepo, s.cleintManager, s.config))

	// non-Spotify users Public routes
	s.router.POST("/auth/non-spotify/register", handlers.RegisterNonSpotifyUser(s.nonSpotifyUserRepo))
	s.router.POST("/auth/non-spotify/verify", handlers.VerifyNonSpotifyUser(s.nonSpotifyUserRepo))

	protected := s.router.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/ping", handlers.Ping())
		protected.GET("/user", handlers.GetUser(s.userRepo))
		protected.GET("/user/top-artists", handlers.GetUserTopArtists(s.cleintManager))
		protected.GET("/user/top-tracks", handlers.GetUserTopTracks(s.cleintManager, s.spotifyService))
		protected.POST("/search", handlers.SearchSongByGenre(s.songRepo))
		protected.POST("/toptracks-analysis", handlers.AnalyzeSongsGivenGenre(s.songRepo, s.cleintManager, s.spotifyService))
		protected.GET("/user/playlists", handlers.GetUserPlaylists(s.spotifySongRepo, s.spotifyService))
		protected.DELETE("/user/playlists/:playlistID", handlers.DeletePlaylist(s.spotifyService, s.spotifySongRepo))
		protected.DELETE("/user/account", handlers.DeleteUserAccount(s.userRepo, s.spotifySongRepo, s.cleintManager))
	}

	nonSpotifyProtected := s.router.Group("/api/non-spotify")
	nonSpotifyProtected.Use(middleware.NonSpotifyAuthMiddleware(s.nonSpotifyUserRepo))
	{
		// Routes for non-Spotify users
		nonSpotifyProtected.POST("/playlists", handlers.GenerateNonSpotifyPlaylist(s.nonSpotifyUserRepo, s.songRepo))
		nonSpotifyProtected.GET("/playlists", handlers.GetNonSpotifyUserPlaylists(s.nonSpotifyUserRepo))
		nonSpotifyProtected.GET("/playlists/:playlistID", handlers.GetNonSpotifyPlaylistDetails(s.nonSpotifyUserRepo))
		nonSpotifyProtected.PATCH("/tracks/:trackID", handlers.UpdateNonSpotifyTrackStatus(s.nonSpotifyUserRepo))
		nonSpotifyProtected.DELETE("/playlists/:playlistID", handlers.DeleteNonSpotifyPlaylist(s.nonSpotifyUserRepo))
	}
}

func (s *Server) Run() error {
	return s.router.Run(":" + s.config.Port)
}
