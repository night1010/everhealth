package entity

import (
	"time"

	"gorm.io/gorm"
)

type AdminContact struct {
	UserId    uint   `gorm:"primaryKey"`
	User      User   `gorm:"foreignKey:UserId;references:Id"`
	Name      string `gorm:"not null"`
	Phone     string `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
