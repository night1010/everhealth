package entity

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Cart struct {
	UserId      uint            `gorm:"primaryKey"`
	User        User            `gorm:"foreignKey:UserId;references:Id"`
	TotalAmount decimal.Decimal `gorm:"not null;type:numeric"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt
}
