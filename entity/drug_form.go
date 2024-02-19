package entity

import (
	"time"

	"gorm.io/gorm"
)

type DrugForm struct {
	Id        uint   `gorm:"primaryKey;autoIncrement"`
	Name      string `gorm:"not null;unique"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
