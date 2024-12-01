package models

import "time"

type Song struct {
	ID         string `gorm:"primaryKey"`
	Name       string `gorm:"column:song_name"` // explicitly name the column
	Artist     string `gorm:"column:artist_name"`
	SpotifyURL string `gorm:"column:spotify_url"`
	ImageURL   string `gorm:"column:image_url"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Playlist struct {
	ID          string `gorm:"primaryKey"`
	UserID      string `gorm:"column:user_id"`
	Name        string `gorm:"column:playlist_name"`
	Description string `gorm:"column:description"`
	URL         string `gorm:"column:url"`
}
