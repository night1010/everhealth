package entity

import (
	"time"

	"gorm.io/gorm"
)

type StatusOrder uint

const (
	WaitingForPayment StatusOrder = iota + 1
	WaitingForPaymentConfirmation
	Processed
	Sent
	OrderConfirmed
	Canceled
)

type OrderStatus struct {
	Id        uint   `gorm:"primaryKey;autoIncrement"`
	Name      string `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
