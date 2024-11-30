package services

import (
	"fmt"

	"github.com/Emeruem-Kennedy1/ghopper/internal/models"
	"github.com/Emeruem-Kennedy1/ghopper/internal/repository"
	"github.com/zmb3/spotify"
)

type SpotifyService struct {
	clientManager   *ClientManager
	spotifySongRepo *repository.SpotifySongRepository
}

func NewSpotifyService(clientManager *ClientManager, spotifySongRepo *repository.SpotifySongRepository) *SpotifyService {
	return &SpotifyService{
		clientManager:   clientManager,
		spotifySongRepo: spotifySongRepo,
	}
}

func (s *SpotifyService) GetSongURL(userID, name, artist string) (string, error) {
	// First check if we have it in our database
	song, err := s.spotifySongRepo.FindSongByNameAndArtist(name, artist)
	if err != nil {
		return "", fmt.Errorf("failed to find song by name and artist: %v", err)
	}

	if song != nil {
		return song.SpotifyURL, nil
	}

	// Get client from ClientManager
	client, exists := s.clientManager.GetClient(userID)
	if !exists {
		return "", fmt.Errorf("no spotify client found for user %s", userID)
	}

	// Search on Spotify
	query := fmt.Sprintf("track:%s artist:%s", name, artist)
	results, err := client.Search(query, spotify.SearchTypeTrack)
	fmt.Println(results.Tracks.Tracks[0].ID, results.Tracks.Tracks[0].Name, results.Tracks.Tracks[0].Artists[0].Name, results.Tracks.Tracks[0].URI)
	if err != nil {
		return "", fmt.Errorf("failed to search Spotify: %v", err)
	}

	if len(results.Tracks.Tracks) == 0 {
		return "", fmt.Errorf("no tracks found for %s by %s", name, artist)
	}

	// Get the first result
	track := results.Tracks.Tracks[0]

	// Create new song record
	newSong := &models.Song{
		ID:         track.ID.String(),
		Name:       name,
		Artist:     artist,
		SpotifyURL: string(track.URI),
	}

	// Save to database
	if err := s.spotifySongRepo.SaveSong(newSong); err != nil {
		return "", fmt.Errorf("failed to save song: %v", err)
	}

	return newSong.SpotifyURL, nil
}
