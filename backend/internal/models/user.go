package models

import (
	"time"
)

type User struct {
	ID           string    `gorm:"primaryKey;type:varchar(255)" json:"id"`
	DisplayName  string    `gorm:"not null" json:"display_name"`
	Email        string    `gorm:"uniqueIndex" json:"email"`
	SpotifyURI   string    `gorm:"uniqueIndex;not null" json:"uri"`
	Country      string    `json:"country"`
	ProfileImage string    `json:"profile_image"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
