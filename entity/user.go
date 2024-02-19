package entity

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	Id           uint   `gorm:"primaryKey;autoIncrement"`
	Email        string `gorm:"unique;not null"`
	Password     string
	IsVerified   bool   `gorm:"not null"`
	Token        string `gorm:"not null"`
	RoleId       RoleId `gorm:"not null"`
	Role         *Role  `gorm:"foreignKey:RoleId;references:Id"`
	AdminContact *AdminContact
	Profile      *Profile
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt
}
