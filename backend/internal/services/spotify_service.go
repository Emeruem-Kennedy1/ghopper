package services

import (
	"fmt"

	"github.com/Emeruem-Kennedy1/ghopper/internal/models"
	"github.com/Emeruem-Kennedy1/ghopper/internal/repository"
	"github.com/zmb3/spotify"
)

type SpotifyService struct {
	clientManager   *ClientManager
	spotifySongRepo repository.SpotifySongRepositoryInterface
}

// removeDuplicates removes duplicate spotify IDs from the slice
func removeDuplicates(ids []spotify.ID) []spotify.ID {
	seen := make(map[spotify.ID]struct{}, len(ids))
	result := make([]spotify.ID, 0, len(ids))

	for _, id := range ids {
		if _, ok := seen[id]; !ok {
			seen[id] = struct{}{}
			result = append(result, id)
		}
	}

	return result
}

func NewSpotifyService(clientManager *ClientManager, spotifySongRepo repository.SpotifySongRepositoryInterface) *SpotifyService {
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

	// remove duplicate song ids first
	songSpotifyIDs = removeDuplicates(songSpotifyIDs)

	// check if playlist already exists in database
	playlist, err := s.spotifySongRepo.FindPlaylistByNameAndUser(playlistName, userID)

	if err != nil {
		return "", fmt.Errorf("failed to find playlist by name and user: %v", err)
	}

	var playlistID spotify.ID
	var playlistURL string

	if playlist != nil {
		// Playlist exists, get its Spotify ID
		playlistID = spotify.ID(playlist.ID)
		playlistURL = playlist.URL
	} else {
		// Create a new playlist
		newPlaylist, err := client.CreatePlaylistForUser(user.ID, playlistName, playlistDescription, false)
		if err != nil {
			return "", fmt.Errorf("failed to create playlist: %v", err)
		}
		playlistID = newPlaylist.ID
		playlistURL = newPlaylist.ExternalURLs["spotify"]

		// Save new playlist to database
		newPlaylistRecord := &models.Playlist{
			ID:          newPlaylist.ID.String(),
			UserID:      userID,
			Name:        playlistName,
			Description: playlistDescription,
			URL:         playlistURL,
		}

		if err := s.spotifySongRepo.SavePlaylist(newPlaylistRecord); err != nil {
			return "", fmt.Errorf("failed to save playlist: %v", err)
		}
	}

	// Get existing tracks in the playlist
	existingTracks, err := client.GetPlaylistTracks(playlistID)
	if err != nil {
		return "", fmt.Errorf("failed to get playlist tracks: %v", err)
	}

	// Create a map of existing track IDs for efficient lookup
	existingTrackIDs := make(map[spotify.ID]struct{})
	for _, track := range existingTracks.Tracks {
		existingTrackIDs[track.Track.ID] = struct{}{}
	}

	// Filter out songs that are already in the playlist
	var newSongs []spotify.ID
	for _, songID := range songSpotifyIDs {
		if _, exists := existingTrackIDs[songID]; !exists {
			newSongs = append(newSongs, songID)
		}
	}

	// Add songs in batches of 100
	for i := 0; i < len(newSongs); i += 100 {
		end := i + 100
		if end > len(newSongs) {
			end = len(newSongs)
		}

		_, err = client.AddTracksToPlaylist(playlistID, newSongs[i:end]...)
		if err != nil {
			return "", fmt.Errorf("failed to add tracks to playlist (batch starting at %d): %v", i, err)
		}
	}

	return playlistURL, nil
}

func (s *SpotifyService) DeletePlaylist(userID, playlistID string) error {
	client, exists := s.clientManager.GetClient(userID)
	if !exists {
		return fmt.Errorf("no spotify client found for user %s", userID)
	}

	playlist, err := client.GetPlaylist(spotify.ID(playlistID))
	if err != nil {
		return fmt.Errorf("failed to get playlist: %v", err)
	}

	if err := client.UnfollowPlaylist(spotify.ID(userID), playlist.ID); err != nil {
		return fmt.Errorf("failed to delete playlist: %v", err)
	}

	return nil
}

func (s *SpotifyService) GetPlaylistImageURL(userID, playlistID string) (string, error) {
	client, exists := s.clientManager.GetClient(userID)
	if !exists {
		return "", fmt.Errorf("no spotify client found for user %s", userID)
	}

	playlist, err := client.GetPlaylist(spotify.ID(playlistID))
	if err != nil {
		return "", fmt.Errorf("failed to get playlist: %v", err)
	}

	return playlist.Images[0].URL, nil
}
