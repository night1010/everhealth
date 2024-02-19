package entity

import (
	"time"

	"gorm.io/gorm"
)

type DoctorSpecialist struct {
	Id        uint   `gorm:"primaryKey;autoIncrement"`
	Name      string `gorm:"not null"`
	Image     string `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
