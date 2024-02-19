package entity

import (
	"time"

	"gorm.io/gorm"
)

type StockMutation struct {
	Id                    uint                `gorm:"primaryKey;autoIncrement"`
	ToPharmacyProductId   uint                `gorm:"not null"`
	ToPharmacyProduct     *PharmacyProduct    `gorm:"notnull;foreignKey:ToPharmacyProductId;references:Id"`
	FromPharmacyProductId uint                `gorm:"not null"`
	FromPharmacyProduct   *PharmacyProduct    `gorm:"notnull;foreignKey:FromPharmacyProductId;references:Id"`
	Quantity              int                 `gorm:"not null"`
	Status                StockMutationStatus `gorm:"not null"`
	OrderId               uint
	MutatedAt             time.Time `gorm:"not null"`
	CreatedAt             time.Time
	UpdatedAt             time.Time
	DeletedAt             gorm.DeletedAt
}

type StockMutationStatus string

const (
	Pending StockMutationStatus = "pending"
	Accept  StockMutationStatus = "accept"
	Decline StockMutationStatus = "decline"
)
