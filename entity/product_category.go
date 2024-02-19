package entity

import (
	"time"

	"gorm.io/gorm"
)

type ProductCategory struct {
	Id        uint   `gorm:"primaryKey;autoIncrement"`
	Name      string `gorm:"not null"`
	IsDrug    bool   `gorm:"not null"`
	Image     string `gorm:"not null"`
	ImageKey  string `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

const (
	ProductCategoryFolder    = "product-category"
	ProductCategoryKeyPrefix = "product-category-"
)
