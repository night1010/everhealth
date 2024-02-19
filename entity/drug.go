package entity

import (
	"time"

	"gorm.io/gorm"
)

type Drug struct {
	ProductId            uint    `gorm:"primaryKey"`
	Product              Product `gorm:"foreignKey:ProductId;references:Id"`
	GenericName          string  `gorm:"not null"`
	DrugFormId           uint    `gorm:"not nul"`
	DrugForm             DrugForm
	DrugClassificationId uint `gorm:"not null"`
	DrugClassification   DrugClassification
	Content              string `gorm:"not null"`
	CreatedAt            time.Time
	UpdatedAt            time.Time
	DeletedAt            gorm.DeletedAt
}
