package entity

import (
	"time"

	"gorm.io/gorm"
)

type Province struct {
	Id          uint   `gorm:"primaryKey;autoIncrement"`
	ProvinceGid int    `gorm:"not null"`
	Name        string `gorm:"not null"`
	Code        string `gorm:"not null"`
	Cities      []*City
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt
}
