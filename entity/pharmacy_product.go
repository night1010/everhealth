package entity

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type PharmacyProduct struct {
	Id         uint `gorm:"primaryKey;autoIncrement"`
	ProductId  uint `gorm:"not null"`
	Product    *Product
	PharmacyId uint `gorm:"not null"`
	Pharmacy   *Pharmacy
	Stock      int             `gorm:"not null"`
	Price      decimal.Decimal `gorm:"not null;type:numeric"`
	IsActive   bool            `gorm:"not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt
}
