package entity

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Telemedicine struct {
	Id                 uint               `gorm:"primaryKey;autoIncrement"`
	OrderedAt          time.Time          `gorm:"not null"`
	ExpiredAt          time.Time          `gorm:"not null"`
	Status             TelemedicineStatus `gorm:"not null"`
	TotalPayment       decimal.Decimal    `gorm:"not null;type:numeric"`
	Proof              string             `gorm:"not null"`
	ProofKey           string             `gorm:"not null"`
	SickLeavePdf       string             `gorm:"not null"`
	SickLeavePdfKey    string             `gorm:"not null"`
	PrescriptionPdf    string             `gorm:"not null"`
	PrescriptionPdfKey string             `gorm:"not null"`
	ProfileId          uint               `gorm:"not null"`
	Profile            *Profile           `gorm:"foreignKey:ProfileId;references:UserId"`
	DoctorId           uint               `gorm:"not null"`
	Doctor             *DoctorProfile     `gorm:"foreignKey:DoctorId;references:ProfileId"`
	Chats              []*Chat
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          gorm.DeletedAt
}

const (
	TelemedicineProofFolder = "telemedicine-proof"
	TelemedicineProofPrefix = "telemedicine-proof-"
	SickLeaveFolder         = "sick-leave"
	SickLeavePrefix         = "sick-leave-"
	PrescriptionFolder      = "prescription"
	PrescriptionPrefix      = "prescription-"
)

type TelemedicineStatus string

const (
	Waiting TelemedicineStatus = "waiting for payment"
	Ongoing TelemedicineStatus = "ongoing"
	End     TelemedicineStatus = "ended"
	Cancel  TelemedicineStatus = "canceled"
)
