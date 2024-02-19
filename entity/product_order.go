package entity

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type ProductOrder struct {
	Id            uint            `gorm:"primaryKey;autoIncrement"`
	OrderedAt     time.Time       `gorm:"not null"`
	OrderStatusId uint            `gorm:"not null"`
	OrderStatus   OrderStatus     `gorm:"foreignKey:OrderStatusId;references:Id"`
	ProfileId     uint            `gorm:"not null"`
	Profile       Profile         `gorm:"foreignKey:ProfileId;references:UserId"`
	ExpiredAt     time.Time       `gorm:"not null"`
	ShippingName  string          `gorm:"not null"`
	ShippingEta   string          `gorm:"not null"`
	ShippingPrice decimal.Decimal `gorm:"not null"`
	AddressId     uint            `gorm:"not null"`
	Address       Address         `gorm:"not null"`
	ItemOrderQty  int             `gorm:"not null"`
	TotalPayment  decimal.Decimal `gorm:"not null;type:numeric"`
	OrderItems    []OrderItem     `gorm:"foreignKey:OrderId"`
	PaymentMethod string
	PaymentProof  string
	ProofKey      string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt
}

const (
	PaymentProofFolder = "payment-proof"
	PaymentProofPrefix = "payment-proof-"
)
