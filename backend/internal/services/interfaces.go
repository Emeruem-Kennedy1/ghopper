package services

import "github.com/zmb3/spotify"

type SpotifyServiceInterface interface {
	GetSongURL(userID, name, artist string) (string, error)
	CreatePlaylistFromSongs(userID string, songSpotifyIDs []spotify.ID, playlistName string, playlistDescription string) (string, error)
	DeletePlaylist(userID, playlistID string) error
	GetPlaylistImageURL(userID, playlistID string) (string, error)
}

type SpotifyClientInterface interface {
	Search(query string, t spotify.SearchType) (*spotify.SearchResult, error)
	CurrentUser() (*spotify.PrivateUser, error)
	CreatePlaylistForUser(userID, name, description string, public bool) (*spotify.FullPlaylist, error)
	GetPlaylistTracks(playlistID spotify.ID) (*spotify.PlaylistTrackPage, error)
	AddTracksToPlaylist(playlistID spotify.ID, trackIDs ...spotify.ID) (string, error)
	GetPlaylist(playlistID spotify.ID) (*spotify.FullPlaylist, error)
	UnfollowPlaylist(userID, playlistID spotify.ID) error
	CurrentUsersTopTracksOpt(opt *spotify.Options) (*spotify.FullTrackPage, error)
	CurrentUsersTopArtistsOpt(opt *spotify.Options) (*spotify.FullArtistPage, error)
}

type ClientManagerInterface interface {
	GetClient(userID string) (SpotifyClientInterface, bool)
	StoreClient(userID string, client SpotifyClientInterface)
	DeleteClient(userID string)
	RemoveClient(userID string)
}


var _ SpotifyServiceInterface = (*SpotifyService)(nil)
var _ SpotifyClientInterface = (*spotify.Client)(nil)
var _ ClientManagerInterface = (*ClientManager)(nil)
