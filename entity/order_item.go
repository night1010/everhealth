package entity

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type OrderItem struct {
	Id                uint            `gorm:"primaryKey;autoIncrement"`
	OrderId           uint            `gorm:"not null"`
	Order             ProductOrder    `gorm:"foreignKey:OrderId;references:Id"`
	PharmacyProductId uint            `gorm:"not null"`
	PharmacyProduct   PharmacyProduct `gorm:"foreignKey:PharmacyProductId;references:Id"`
	Quantity          int             `gorm:"not null"`
	SubTotal          decimal.Decimal `gorm:"not null"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt
}
