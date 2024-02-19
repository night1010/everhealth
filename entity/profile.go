package entity

import (
	"time"

	"gorm.io/gorm"
)

type Profile struct {
	UserId    uint   `gorm:"primaryKey"`
	User      User   `gorm:"foreignKey:UserId;references:Id"`
	Name      string `gorm:"not null"`
	Image     string
	ImageKey  string
	Birthdate time.Time `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

const (
	ProfilePhotoFolder    = "profile-photo"
	ProfilePhotoKeyPrefix = "profile-photo-"
)
