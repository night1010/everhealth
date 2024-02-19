package entity

import (
	"time"

	"gorm.io/gorm"
)

type ForgotPasswordToken struct {
	Id        uint      `gorm:"primaryKey;autoIncrement"`
	Token     string    `gorm:"not null"`
	ExpiredAt time.Time `gorm:"not null"`
	IsActive  bool      `gorm:"not null"`
	UserId    uint      `gorm:"not null"`
	User      User      `gorm:"foreignKey:UserId;references:Id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
