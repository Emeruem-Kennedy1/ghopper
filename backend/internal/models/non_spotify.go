package models

import (
	"time"
)

type NonSpotifyUser struct {
	ID         string    `gorm:"primaryKey" json:"id"`
	Passphrase string    `gorm:"not null" json:"-"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type NonSpotifyPlaylist struct {
	ID          string    `gorm:"primaryKey" json:"id"`
	UserID      string    `gorm:"not null" json:"user_id"`
	Name        string    `gorm:"not null" json:"name"`
	Genre       string    `json:"genre"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type NonSpotifyPlaylistTrack struct {
	ID              string    `gorm:"primaryKey" json:"id"`
	PlaylistID      string    `gorm:"not null" json:"playlist_id"`
	Title           string    `gorm:"not null" json:"title"`
	Artist          string    `gorm:"not null" json:"artist"`
	AddedToPlaylist bool      `gorm:"default:false" json:"added_to_playlist"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type NonSpotifyPlaylistSeedTrack struct {
	ID         string    `gorm:"primaryKey" json:"id"`
	PlaylistID string    `gorm:"not null" json:"playlist_id"`
	Title      string    `gorm:"not null" json:"title"`
	Artist     string    `gorm:"not null" json:"artist"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type NonSpotifyPlaylistWithTracks struct {
	NonSpotifyPlaylist
	Tracks     []NonSpotifyPlaylistTrack     `json:"tracks"`
	SeedTracks []NonSpotifyPlaylistSeedTrack `json:"seed_tracks"`
}
