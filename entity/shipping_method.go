package entity

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type ShippingMethod struct {
	Id         uint            `gorm:"primaryKey;autoIncrement"`
	Name       string          `gorm:"not null"`
	Duration   string          `gorm:"not null"`
	PricePerKM decimal.Decimal `gorm:"not null;type:numeric"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt
}

type CalculatedShippingMethod struct {
	Name              string
	EstimatedDuration string
	Cost              string
}
