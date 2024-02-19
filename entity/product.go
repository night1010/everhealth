package entity

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Product struct {
	Id                uint   `gorm:"primaryKey;autoIncrement"`
	Name              string `gorm:"not null"`
	Manufacture       string `gorm:"not null"`
	Detail            string `gorm:"not null"`
	ProductCategoryId uint   `gorm:"not null"`
	ProductCategory   ProductCategory
	UnitInPack        string          `gorm:"not null"`
	Price             decimal.Decimal `gorm:"not null;type:numeric"`
	SellingUnit       string          `gorm:"not null"`
	Weight            decimal.Decimal `gorm:"not null;type:numeric"`
	Height            decimal.Decimal `gorm:"not null;type:numeric"`
	Length            decimal.Decimal `gorm:"not null;type:numeric"`
	Width             decimal.Decimal `gorm:"not null;type:numeric"`
	Image             string          `gorm:"not null"`
	ImageKey          string          `gorm:"not null"`
	IsHidden          bool            `gorm:"not null"`
	Drug              *Drug
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt
}

const (
	ProductFolder    = "product"
	ProductKeyPrefix = "product-"
)
