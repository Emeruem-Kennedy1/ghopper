package models

import "time"

type Song struct {
	ID         string `gorm:"primaryKey"`
	Name       string `gorm:"column:song_name"` // explicitly name the column
	Artist     string `gorm:"column:artist_name"`
	SpotifyURL string `gorm:"column:spotify_url"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
