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
	if err != nil {
		return "", fmt.Errorf("failed to search Spotify: %v", err)
	}

	// if we don't have any results return empty string
	if len(results.Tracks.Tracks) == 0 {
		return "", nil
	}

	// Get the first result
	track := results.Tracks.Tracks[0]
	baseUrl := "https://open.spotify.com/track/"

	// Create new song record
	newSong := &models.Song{
		ID:         track.ID.String(),
		Name:       name,
		Artist:     artist,
		SpotifyURL: baseUrl + track.ID.String(),
		ImageURL:   track.Album.Images[0].URL,
	}

	// Save to database
	if err := s.spotifySongRepo.SaveSong(newSong); err != nil {
		return "", fmt.Errorf("failed to save song: %v", err)
	}

	return newSong.SpotifyURL, nil
}

func (s *SpotifyService) CreatePlaylistFromSongs(userID string, songSpotifyIDs []spotify.ID, playlistName string, playlistDescription string) (string, error) {
	client, exists := s.clientManager.GetClient(userID)
	if !exists {
		return "", fmt.Errorf("no spotify client found for user %s", userID)
	}

	user, err := client.CurrentUser()
	if err != nil {
		return "", fmt.Errorf("failed to get current user: %v", err)
	}

	// check if playlist already exists in my database
	playlist, _ := s.spotifySongRepo.FindPlaylistByNameAndUser(playlistName, userID)

	if playlist != nil {
		return playlist.URL, nil
	}

	// Create a new playlist 
	newPlaylist, err := client.CreatePlaylistForUser(user.ID, playlistName, playlistDescription, false)
	if err != nil {
		return "", fmt.Errorf("failed to create playlist: %v", err)
	}

	// Add songs to the playlist
	_, err = client.AddTracksToPlaylist(newPlaylist.ID, songSpotifyIDs...)
	if err != nil {
		return "", fmt.Errorf("failed to add tracks to playlist: %v", err)
	}

	// Save to database
	newPlaylistRecord := &models.Playlist{
		ID:          newPlaylist.ID.String(),
		UserID:      userID,
		Name:        playlistName,
		Description: playlistDescription,
		URL:         newPlaylist.ExternalURLs["spotify"],
	}

	if err := s.spotifySongRepo.SavePlaylist(newPlaylistRecord); err != nil {
		return "", fmt.Errorf("failed to save playlist: %v", err)
	}

	return newPlaylist.ExternalURLs["spotify"], nil
}
