package entity

import (
	"time"

	"gorm.io/gorm"
)

type MessageType string

const (
	MessageTypeText  MessageType = "text"
	MessageTypeImage MessageType = "image"
	MessageTypePdf   MessageType = "pdf"
)

type Chat struct {
	Id             uint      `gorm:"primaryKey;autoIncrement"`
	TelemedicineId uint      `gorm:"not null"`
	UserId         uint      `gorm:"not null"`
	ChatTime       time.Time `gorm:"not null"`
	Telemedicine   Telemedicine
	Message        string
	MessageType    MessageType
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Deleted        gorm.DeletedAt
}
